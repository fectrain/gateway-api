package models

type UserResponse struct {
	Data      *UserData `json:"data" structs:"data"`
	Message   string    `json:"message" structs:"message"`
	ErrorCode int32     `json:"error-code" structs:"error-code"`
}

type UserData struct {
	Username string
}

type ItemInfoResponse struct {
	Data      *[]ItemInfo `json:"data" structs:"data"`
	Message   string      `json:"message" structs:"message"`
	ErrorCode int32       `json:"error-code" structs:"error-code"`
}

type PriceHistoryResponse struct {
	Data      *[]PriceHistory `json:"data" structs:"data"`
	Message   string          `json:"message" structs:"message"`
	ErrorCode int32           `json:"error-code" structs:"error-code"`
}

type ItemInfo struct {
	ItemId int64
	ItemName string
}

type PriceHistory struct {
	Price     int64
	Timestamp int64
}
