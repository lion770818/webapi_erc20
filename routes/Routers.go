package routes

import (
	"webapi_erc20/controller"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("v1")
	{
		// 創建錢包
		v1.POST("/address", controller.AddressCreate)
		v1.GET("/:address/balance/:cryptoType", controller.GetBalance)
	}

	return r
}
