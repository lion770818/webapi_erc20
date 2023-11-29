package entity

import (
	"errors"
	response "webapi_erc20/common/rsp"
	"webapi_erc20/define"

	"github.com/shopspring/decimal"
)

type WithdrawCreateReq struct {
	MerchantID string `json:"merchant_id" binding:"required"`

	CryptoType string `json:"crypto_type" binding:"required"`
	ChainType  string `json:"chain_type" binding:"required"`

	SecretKey string `json:"secret_key" binding:"required"`

	FromAddress string `json:"from_address" binding:"required"`
	ToAddress   string `json:"to_address" binding:"required"`

	Amount decimal.Decimal `json:"amount" binding:"required"`
	Memo   string          `json:"memo"`
}

func (req *WithdrawCreateReq) Validate() (response.Status, error) {
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

type WithdrawCreateResp struct {
	TxHash string `json:"tx_hash"`
	Memo   string `json:"memo"`
}
