package main

import (
	"gateway-api/controllers"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const baseGroup string = "api/v1"

var r *gin.Engine

const (
	redisAddr          = "192.168.18.63:6379"
	sessionStorageName = "userSession"

	serverPort = ":8090"
)

func init() {
	r = gin.Default()

	//redis session setup
	//store, _ := redis.NewStore(10, "tcp", redisAddr, "", []byte("secret"))
	//expire time for session storage
	//store.Options(sessions.Options{MaxAge: 30 * 60})
	//r.Use(sessions.Sessions(sessionStorageName, store)).Use(middleware.MwPrometheusHttp)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

}

func main() {
	api := r.Group(baseGroup)
	r.GET(baseGroup+"/login", controllers.LoadTest)
	api.POST("/register", controllers.Register)
	//api.GET("/item-list", controllers.GetItemInfoList)
	//api.GET("/item-price-history", controllers.GetItemPriceHistory)
	//api.GET("/add-item", controllers.AddItem)
	//api.GET("/delete-item", controllers.DeleteItem)
	pprof.Register(r)
	err := r.Run(serverPort)
	if err != nil {
		return
	}

}
