package dao

import (
	"context"
	"webapi_erc20/common/logs"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Tokens struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	CryptoType   string `gorm:"column:crypto_type"`
	ChainType    string `gorm:"column:chain_type"`
	ContractAddr string `gorm:"column:contract_addr"`

	Decimals int             `gorm:"column:decimals"`
	GasLimit int64           `gorm:"column:gas_limit"`
	GasPrice decimal.Decimal `gorm:"column:gas_price"`

	ContractAbi string `gorm:"column:contract_abi"`
}

type TokensRepository struct {
	Tokens
	db *gorm.DB
}

func (t *Tokens) TableName() string {
	return "tokens"
}

var tokensInstance *TokensRepository

func GetTokenInstance() *TokensRepository {
	if tokensInstance == nil {
		tokensInstance = &TokensRepository{
			db: SqlSession,
		}
	}
	return tokensInstance
}

// 創建table
func (r *TokensRepository) CreateTable() {

	exist := r.db.Migrator().HasTable(&Tokens{})
	if !exist {
		logs.Debugf("創建table")
		r.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='合約token'").Create(&Tokens{})
	}

	// 自動遷移 schema
	r.db.AutoMigrate(&Tokens{})
}

func (repo *TokensRepository) GetByCryptoType(ctx context.Context, cryptoType string) (Tokens, error) {
	result := Tokens{}
	err := repo.db.WithContext(ctx).
		Where("`crypto_type` = ?", cryptoType).
		Take(&result).Error
	if err != nil {
		return Tokens{}, err
	}

	return result, nil
}

func (repo *TokensRepository) GetList(ctx context.Context) ([]Tokens, error) {
	result := make([]Tokens, 0, 200)
	err := repo.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
