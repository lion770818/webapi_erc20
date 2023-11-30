package cryp_notify

import "github.com/shopspring/decimal"

type CreateTransactionNotifyReq struct {
	TxType int `json:"tx_type"`

	BlockHeight int64  `json:"block_height"`
	TxHash      string `json:"tx_hash"`

	CryptoType string `json:"crypto_type"`
	ChainType  string `json:"chain_type"`

	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`

	Amount    decimal.Decimal `json:"amount"`
	Fee       decimal.Decimal `json:"fee"`
	FeeCrypto string          `json:"fee_crypto"`

	Status int    `json:"status"`
	Memo   string `json:"memo"`
}
