package dao

type BlockHeight struct {
	ID int64 `gorm:"column:id"`

	CreateTime int64 `gorm:"column:create_time"`
	UpdateTime int64 `gorm:"column:update_time"`

	BlockHeight int64 `gorm:"column:block_height"`
}

func (bh *BlockHeight) TableName() string {
	return "block_height"
}
