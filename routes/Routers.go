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
		// 取得錢包餘額
		v1.GET("/:address/balance/:cryptoType", controller.GetBalance) // http://127.0.0.1:8888/v1/0x3DE84eDBa9d3829F1eE18F297BAAD59E9Fe3855F/balance/ETH/?chain_type=ETH
		// 錢包交易
		v1.POST("/withdraw", controller.Withdraw)
		// 取得交易清單明細
		//v1.GET("/tx/:txHash", controller.GetTxHash)

	}

	return r
}
