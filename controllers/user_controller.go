package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"gateway-api/models"
	"gateway-api/proto/bp"
	"git.garena.com/shopee/platform/tracing"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

var userClient bp.UserServiceClient

const (
	keyStringSessionId = "sessionId"
	userServerAddress  = "192.168.18.63:8081"
	itemServerAddress  = "192.168.18.63:50051"
)

const (
	success             int32 = 0
	notLogin            int32 = 1
	sessionExpired      int32 = 2
	serverInternalError int32 = 3
	illegalParamError   int32 = 4
)

func init() {
	//gRpc client
	userConn, e := grpc.Dial(userServerAddress, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}
	userClient = bp.NewUserServiceClient(userConn)
	itemConn, e := grpc.Dial(itemServerAddress, grpc.WithInsecure())
	if e != nil {
		panic(e)
	}
	itemClient = bp.NewItemPriceServiceClient(itemConn)
	//defer conn.Close()
}

func LoadTest(c *gin.Context)  {
	time.Sleep(100 * time.Microsecond)
	log.Println("request done..")
	sc := tracing.GetSpanContext(c.Request.Context())
	log.Println(sc.TraceIDString())
	c.JSON(http.StatusOK, "{}")

}

func Login(c *gin.Context) {

	var userInfo models.UserInfo
	err := c.ShouldBind(&userInfo)
	log.Println("user info client: " + userInfo.Username + ":" + userInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Incorrect param", ErrorCode: illegalParamError})
		return
	}
	resp, err := userClient.Login(context.Background(), &bp.LoginRequest{Username: userInfo.Username, Password: userInfo.Password})
	if err != nil {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Login fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}

	log.Printf("resp code: %v, user id: %d\n", resp.GetIsSuccess(), resp.GetUserId())
	if !resp.IsSuccess {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Login fail: " + resp.GetError().Message, ErrorCode: serverInternalError})
		return
	}

	session := sessions.Default(c)

	sessionId := fmt.Sprintf("%x_%d", md5.Sum(IntToBytes(resp.GetUserId())), time.Now().Unix())
	session.Set(sessionId, resp.GetUserId())
	session.Save()

	c.SetCookie(keyStringSessionId, sessionId, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, models.UserResponse{Message: "Login successfully", ErrorCode: success, Data: &models.UserData{Username: resp.Username}})
}

func Register(c *gin.Context) {
	var userInfo models.UserInfo
	err := c.ShouldBind(&userInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Incorrect param", ErrorCode: illegalParamError})
		return
	}
	resp, err := userClient.Register(context.Background(), &bp.RegisterRequest{Username: userInfo.Username, Password: userInfo.Password})
	if err != nil {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Register fail: " + err.Error(), ErrorCode: serverInternalError})
		return
	}
	if !resp.IsSuccess {
		c.JSON(http.StatusOK, models.UserResponse{Message: "Register fail: " + resp.GetError().Message, ErrorCode: serverInternalError})
		return
	}
	c.JSON(http.StatusOK, models.UserResponse{Message: "Register successfully", ErrorCode: success})

}

func IntToBytes(n int64) []byte {
	//data := int64(n)
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, n)
	return byteBuf.Bytes()
}
