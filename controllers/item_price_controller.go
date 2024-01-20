package controllers

import (
	"gateway-api/models"
	"gateway-api/proto/bp"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

var itemClient bp.ItemPriceServiceClient

//var UserVerifyEnable = true

const (
	defaultUserId int64 = 1
	notFound            = -1
)

func GetItemInfoList(c *gin.Context) {
	responseBody, userId := verifySessionId(c)
	if userId == notFound {
		c.JSON(http.StatusOK, responseBody)
		return
	}
	resp, err := itemClient.GetItemInfoList(context.Background(), &bp.ItemInfoRequest{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Query item fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}

	itemData := resp.GetItemInfo()
	returnData := make([]models.ItemInfo, len(itemData))
	for i, item := range itemData {
		returnData[i] = models.ItemInfo{ItemId: item.ItemId, ItemName: item.ItemName}
	}

	c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Query successfully", ErrorCode: success, Data: &returnData})
}

func GetItemPriceHistory(c *gin.Context) {
	responseBody, userId := verifySessionId(c)
	if userId == notFound {
		c.JSON(http.StatusOK, responseBody)
		return
	}

	itemIdStr := c.Query("item_id")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if itemIdStr == "" || err != nil {
		c.JSON(http.StatusOK, models.PriceHistoryResponse{Message: "Illegal item id", ErrorCode: illegalParamError})
		return
	}

	resp, err := itemClient.GetItemPriceHistoryByItem(context.Background(), &bp.ItemPriceRequest{ItemId: itemId})
	if err != nil {
		c.JSON(http.StatusOK, models.PriceHistoryResponse{Message: "Query items fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}

	itemData := resp.GetItemPrice()
	returnData := make([]models.PriceHistory, len(itemData))
	for i, item := range itemData {
		returnData[i] = models.PriceHistory{Timestamp: item.Timestamp, Price: item.Price}
	}
	c.JSON(http.StatusOK, models.PriceHistoryResponse{Message: resp.GetMessage(), ErrorCode: resp.ErrorCode, Data: &returnData})
}

func AddItem(c *gin.Context) {
	responseBody, userId := verifySessionId(c)
	if userId == notFound {
		c.JSON(http.StatusOK, responseBody)
		return
	}
	itemId := c.Query("item_id")
	itemName := c.Query("item_name")

	itemId64, err := strconv.ParseInt(itemId, 10, 64)

	if itemId == "" || err != nil {
		c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Illegal item id", ErrorCode: illegalParamError})
		return
	}

	resp, err := itemClient.AddItem(context.Background(), &bp.ItemOpRequest{UserId: userId, ItemId: itemId64, ItemName: itemName})
	if err != nil {
		c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Add item fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}
	c.JSON(http.StatusOK, models.ItemInfoResponse{Message: resp.GetMessage(), ErrorCode: resp.ErrorCode})
}

func DeleteItem(c *gin.Context) {
	responseBody, userId := verifySessionId(c)
	if userId == notFound {
		c.JSON(http.StatusOK, responseBody)
		return
	}
	itemId := c.Query("item_id")

	itemId64, err := strconv.ParseInt(itemId, 10, 64)

	if itemId == "" || err != nil {
		c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Illegal item id", ErrorCode: illegalParamError})
		return
	}

	resp, err := itemClient.DeleteItem(context.Background(), &bp.ItemOpRequest{UserId: userId, ItemId: itemId64})
	if err != nil {
		c.JSON(http.StatusOK, models.ItemInfoResponse{Message: "Delete item fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}

	c.JSON(http.StatusOK, models.ItemInfoResponse{Message: resp.GetMessage(), ErrorCode: resp.ErrorCode})

}

func verifySessionId(c *gin.Context) (models.ItemInfoResponse, int64) {
	//if UserVerifyEnable {
	requestSessionId, err := c.Cookie(keyStringSessionId)

	if err != nil || requestSessionId == "" {
		return models.ItemInfoResponse{Message: "You haven't login", ErrorCode: notLogin}, notFound
	}

	session := sessions.Default(c)
	userId := session.Get(requestSessionId)
	if userId == nil {
		return models.ItemInfoResponse{Message: "Session expired, please login", ErrorCode: sessionExpired}, notFound
	}
	return models.ItemInfoResponse{}, userId.(int64)
	//}
	//return models.ItemInfoResponse{}, defaultUserId
}
