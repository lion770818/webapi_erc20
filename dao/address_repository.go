package dao

import (
	"context"
	"webapi_erc20/common/logs"
	"webapi_erc20/common/utils"

	"gorm.io/gorm"
)

// db repository

type Address struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	MerchantType int    `gorm:"column:merchant_type"`
	Address      string `gorm:"column:address"` // 錢包地址

	ChainType string `gorm:"column:chain_type"`
}

type AddressRepository struct {
	Address
	db *gorm.DB
}

func (addr *Address) TableName() string {
	return "address"
}

var addressInstance *AddressRepository

func GetAddressInstance() *AddressRepository {
	if addressInstance == nil {
		addressInstance = &AddressRepository{
			db: SqlSession,
		}
	}
	return addressInstance
}

// 創建table
func (r *AddressRepository) CreateTable() {

	// 自動遷移 schema
	r.db.AutoMigrate(&Address{})

	// 創建table
	//r.db.Create(&Address{})
	r.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用戶錢包地址'").Create(&Address{})
}

// 創建錢包地址
func (repo *AddressRepository) Create(ctx context.Context, addr Address) (int, error) {
	addr.CreateTime = utils.TimeNowUnix()
	addr.UpdateTime = utils.TimeNowUnix()

	logs.Debugf("創建錢包地址 addr:%+v", addr)
	result := repo.db.WithContext(ctx).Create(&addr)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(addr.ID), nil
}

// 取得錢包地址
func (repo *AddressRepository) GetByAddress(ctx context.Context, addr string) (Address, error) {
	result := Address{}
	err := repo.db.WithContext(ctx).
		Where("`address` = ?", addr).
		Take(&result).Error
	if err != nil {
		return Address{}, err
	}

	return result, nil
}
