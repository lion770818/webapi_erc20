package controller

import (
	"context"
	"net/http"
	"webapi_erc20/common/logs"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/common/utils"
	"webapi_erc20/entity"
	"webapi_erc20/service"

	"github.com/gin-gonic/gin"
)

// 創建錢包
func AddressCreate(c *gin.Context) {

	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	var req entity.AddressCreateReq
	if err := response.ShouldBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(response.CodeInternalError))
		return
	}

	if respStatus, err := req.Validate(); err != nil {

		logs.Errorf("req.Validate err:%v", err)

		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	result, respStatus, err := service.Create(ctx, req)
	if err != nil {
		logs.Errorf("create address err:%v", err)

		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}

// 取得錢包餘額
func GetBalance(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	defer cancelCtx()

	var req entity.AddressGetBalanceReq
	if respStatus, err := req.Parse(c); err != nil {
		logs.Errorf("req.Parse err:%v", err)

		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	result, respStatus, err := service.GetBalance(ctx, req)
	if err != nil {
		logs.Errorf("getBalance err:%v", err)
		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}
