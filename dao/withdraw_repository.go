package dao

import (
	"context"
	"webapi_erc20/common/logs"
	"webapi_erc20/common/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Withdraw struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	MerchantType int `gorm:"column:merchant_type"`

	TxHash     string `gorm:"column:tx_hash"`
	CryptoType string `gorm:"column:crypto_type"`
	ChainType  string `gorm:"column:chain_type"`

	FromAddress string `gorm:"column:from_address"`
	ToAddress   string `gorm:"column:to_address"`

	Amount decimal.Decimal `gorm:"column:amount"`
	Memo   string          `gorm:"column:memo"`
}

type WithdrawRepository struct {
	Withdraw
	db *gorm.DB
}

func (w *Withdraw) TableName() string {
	return "withdraw"
}

var withdrawInstance *WithdrawRepository

func GetWithdrawInstance() *WithdrawRepository {
	if withdrawInstance == nil {
		withdrawInstance = &WithdrawRepository{
			db: SqlSession,
		}
	}
	return withdrawInstance
}

// 創建table
func (r *WithdrawRepository) CreateTable() {

	exist := r.db.Migrator().HasTable(&Withdraw{})
	if !exist {
		logs.Debugf("創建table")
		r.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易清單'").Create(&Withdraw{})
	}

	// 自動遷移 schema
	r.db.AutoMigrate(&Withdraw{})

}

func (repo *WithdrawRepository) Create(ctx context.Context, w Withdraw) (int, error) {
	w.CreateTime = utils.TimeNowUnix()
	w.UpdateTime = utils.TimeNowUnix()

	result := repo.db.WithContext(ctx).Create(&w)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(w.ID), nil
}

func (repo *WithdrawRepository) GetByTxHash(ctx context.Context, txHash string) (Withdraw, error) {
	result := Withdraw{}
	err := repo.db.WithContext(ctx).
		Where("`tx_hash` = ?", txHash).
		Take(&result).Error
	if err != nil {
		return Withdraw{}, err
	}

	return result, nil
}
