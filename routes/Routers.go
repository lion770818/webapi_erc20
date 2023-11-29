package routes

import (
	"webapi_erc20/controller"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	/**
	用户User路由组
	*/
	userGroup := r.Group("user")
	{
		//增加用户User
		userGroup.POST("/users", controller.CreateUser)
		//查看所有的User
		userGroup.GET("/users", controller.GetUserList)
		//修改某个User
		userGroup.PUT("/users/:id", controller.UpdateUser)
		//删除某个User
		userGroup.DELETE("/users/:id", controller.DeleteUserById)
	}

	v1 := r.Group("v1")
	{
		// 創建錢包
		v1.POST("/address", controller.AddressCreate)
		v1.GET("/:address/balance/:cryptoType", controller.GetBalance)
	}

	return r
}
