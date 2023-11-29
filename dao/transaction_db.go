package dao

import "github.com/shopspring/decimal"

type Transaction struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	TxType           int    `gorm:"column:tx_type"`
	BlockHeight      int64  `gorm:"column:block_height"`
	TransactionIndex int    `gorm:"column:transaction_index"`
	TxHash           string `gorm:"column:tx_hash"`

	CryptoType string `gorm:"column:crypto_type"`
	ChainType  string `gorm:"column:chain_type"`

	ContractAddr string `gorm:"column:contract_addr"`
	FromAddress  string `gorm:"column:from_address"`
	ToAddress    string `gorm:"column:to_address"`

	Amount decimal.Decimal `gorm:"column:amount"`

	Gas      int64           `gorm:"column:gas"`
	GasUsed  int64           `gorm:"column:gas_used"`
	GasPrice decimal.Decimal `gorm:"column:gas_price"`

	Fee       decimal.Decimal `gorm:"column:fee"`
	FeeCrypto string          `gorm:"column:fee_crypto"`

	Confirm int64  `gorm:"column:confirm"`
	Status  int    `gorm:"column:status"`
	Memo    string `gorm:"column:memo"`

	NotifyStatus int `gorm:"column:notify_status"`
}

func (trans *Transaction) TableName() string {
	return "transaction"
}
