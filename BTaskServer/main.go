package main

import (
	"BTaskServer/common"
	"BTaskServer/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	if !common.InitConfig() {
		return
	}
	if !common.InitDB() {
		return
	}
	if !common.InitRedis() {
		return
	}
	if !common.InintTrans() {
		return
	}

	InitServer()
}

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r = routes.CollectRoute(r)
	addr := viper.GetString("server.addr")
	port := viper.GetString("server.port")
	fmt.Println(addr)
	panic(r.Run(addr + ":" + port))
}
