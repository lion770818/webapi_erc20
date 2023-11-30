package dao

import (
	"context"
	"webapi_erc20/common/logs"
	"webapi_erc20/common/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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

type TransactionRepository struct {
	Transaction
	db *gorm.DB
}

func (trans *Transaction) TableName() string {
	return "transaction"
}

var transactionInstance *TransactionRepository

func GetTransactionInstance() *TransactionRepository {
	if transactionInstance == nil {
		transactionInstance = &TransactionRepository{
			db: SqlSession,
		}
	}
	return transactionInstance
}

// 創建table
func (r *TransactionRepository) CreateTable() {

	exist := r.db.Migrator().HasTable(&Transaction{})
	if !exist {
		logs.Debugf("創建table")
		r.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易清單'").Create(&Transaction{})
	}

	// 自動遷移 schema
	r.db.AutoMigrate(&Transaction{})
}

func (repo *TransactionRepository) Create(ctx context.Context, trans Transaction) (int, error) {
	trans.CreateTime = utils.TimeNowUnix()
	trans.UpdateTime = utils.TimeNowUnix()

	result := repo.db.WithContext(ctx).Create(&trans)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(trans.ID), nil
}

func (repo *TransactionRepository) Update(ctx context.Context, trans Transaction) error {
	trans.UpdateTime = utils.TimeNowUnix()

	return repo.db.WithContext(ctx).Model(&trans).
		Where("id = ?", trans.ID).
		Updates(&trans).Error
}

func (repo *TransactionRepository) GetByTxHash(ctx context.Context, txHash string) (Transaction, error) {
	result := Transaction{}
	err := repo.db.WithContext(ctx).
		Where("`tx_hash` = ?", txHash).
		Take(&result).Error
	if err != nil {
		return Transaction{}, err
	}

	return result, nil
}

func (repo *TransactionRepository) GetListByStatusAndBlockHeight(ctx context.Context, status, limit int, blockHeight int64) ([]Transaction, error) {
	result := make([]Transaction, 0, 0)
	err := repo.db.WithContext(ctx).
		Where("`status` = ?", status).
		Where("`block_height` <= ?", blockHeight).
		Limit(limit).
		Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *TransactionRepository) GetListByNotifyStatus(ctx context.Context, notifyStatus int) ([]Transaction, error) {
	result := make([]Transaction, 0, 0)
	err := repo.db.WithContext(ctx).
		Where("`notify_status` = ?", notifyStatus).
		Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
