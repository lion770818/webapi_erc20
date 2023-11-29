package entity

import (
	"errors"
	"webapi_erc20/common/logs"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/common/utils"
	"webapi_erc20/define"

	"github.com/gin-gonic/gin"

	"github.com/shopspring/decimal"
)

type AddressCreateReq struct {
	MerchantID string `json:"merchant_id" binding:"required"`

	ChainType string `json:"chain_type" binding:"required"`
}

func (req *AddressCreateReq) Validate() (response.Status, error) {
	_, ok := define.MerchantID2Type[req.MerchantID]
	if !ok {
		respStatus := response.CodeParamInvalid.WithMsg(", merchant_id is not exist")
		return respStatus, errors.New(respStatus.Messages)
	}

	if req.ChainType != define.ChainType {
		respStatus := response.CodeParamInvalid.WithMsg(", chain_type need use ETH")
		return respStatus, errors.New(respStatus.Messages)
	}

	return response.Status{}, nil
}

type AddressCreateResp struct {
	Address string `json:"address"`

	SecretKey string `json:"secret_key"`
	PublicKey string `json:"public_key"`
}

type AddressGetBalanceReq struct {
	Address string

	CryptoType string
	ChainType  string
}

func (req *AddressGetBalanceReq) Parse(c *gin.Context) (response.Status, error) {
	req.Address = c.Param("address")
	if utils.IsEmpty(req.Address) {
		respStatus := response.CodeParamInvalid.WithMsg(", address is empty")
		return respStatus, errors.New(respStatus.Messages)
	}

	req.CryptoType = c.Param("cryptoType")
	if utils.IsEmpty(req.CryptoType) {
		respStatus := response.CodeParamInvalid.WithMsg(", crypto_type is empty")
		return respStatus, errors.New(respStatus.Messages)
	}

	logs.Debugf("req:%+v", req)
	req.ChainType = c.Query("chain_type")
	if req.ChainType != define.ChainType {
		respStatus := response.CodeParamInvalid.WithMsg(", chain_type need use ETH")
		return respStatus, errors.New(respStatus.Messages)
	}

	return response.Status{}, nil
}

type AddressGetBalanceResp struct {
	Balance decimal.Decimal `json:"balance"`
}
