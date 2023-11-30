package dao

import (
	"context"
	"webapi_erc20/common/logs"
	"webapi_erc20/common/utils"

	"gorm.io/gorm"
)

type BlockHeight struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	BlockHeight int64 `gorm:"column:block_height"`
}

type blockHeightRepository struct {
	db *gorm.DB
}

func (bh *BlockHeight) TableName() string {
	return "block_height"
}

var blockHeightInstance *blockHeightRepository

func GetBlockHeightInstance() *blockHeightRepository {
	if blockHeightInstance == nil {
		blockHeightInstance = &blockHeightRepository{
			db: SqlSession,
		}
	}
	return blockHeightInstance
}

// func (repo *blockHeightRepository) New(db *gorm.DB) BlockHeightRepository {
// 	result := &blockHeightRepository{
// 		db: db,
// 	}

// 	return result
// }

// 創建table
func (r *blockHeightRepository) CreateTable() {

	exist := r.db.Migrator().HasTable(&BlockHeight{})
	if !exist {
		logs.Debugf("創建table")
		r.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='區塊鍊高度'").Create(&BlockHeight{})
	}

	// 自動遷移 schema
	r.db.AutoMigrate(&BlockHeight{})

}

func (repo *blockHeightRepository) Create(ctx context.Context, bh BlockHeight) (int, error) {
	bh.CreateTime = utils.TimeNowUnix()
	bh.UpdateTime = utils.TimeNowUnix()

	result := repo.db.WithContext(ctx).Create(&bh)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(bh.ID), nil
}

func (repo *blockHeightRepository) Update(ctx context.Context, bh BlockHeight) error {
	bh.UpdateTime = utils.TimeNowUnix()

	return repo.db.WithContext(ctx).Model(&bh).
		Where("id = ?", bh.ID).
		Updates(&bh).Error
}

func (repo *blockHeightRepository) Get(ctx context.Context) (BlockHeight, error) {
	result := BlockHeight{}
	err := repo.db.WithContext(ctx).
		Take(&result).Error
	if err != nil {
		return BlockHeight{}, err
	}

	return result, nil
}
