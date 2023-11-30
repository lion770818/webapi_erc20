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

// 轉帳
func Withdraw(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*utils.Time30S)
	defer cancelCtx()

	var req entity.WithdrawCreateReq
	if err := response.ShouldBindJSON(c, &req); err != nil {
		logs.Errorf("shouldBindJSON err:%v", err)
		c.JSON(http.StatusBadRequest, response.NewError(response.CodeInternalError))
		return
	}

	if respStatus, err := req.Validate(); err != nil {
		logs.Errorf("validate err:%v", err)
		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	result, respStatus, err := service.Withdraw(ctx, req)
	if err != nil {
		logs.Errorf("withdraw err:%v", err)
		c.JSON(http.StatusBadRequest, response.NewError(respStatus))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(result))
}

func GetTxHash(c *gin.Context) {
	// ctx, cancelCtx := context.WithTimeout(context.Background(), utils.Time30S)
	// defer cancelCtx()

	// var req entity.TransGetTxHashReq
	// if respStatus, err := req.Parse(c); err != nil {
	// 	logs.Errorf("req.Parse err:%v", err)

	// 	c.JSON(http.StatusBadRequest, response.NewError(respStatus))
	// 	return
	// }

	// result, respStatus, err := service.GetByTxHash(ctx, req)
	// if err != nil {
	// 	logs.Errorf("getByTxHash err:%v", err)

	// 	c.JSON(http.StatusBadRequest, response.NewError(respStatus))
	// 	return
	// }

	// c.JSON(http.StatusOK, response.NewSuccess(result))
}
