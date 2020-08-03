package dao

import (
	"time"
)

type MakerOrder struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId uint       `json:"source_chain_id"`
	TargetChainId uint       `json:"target_chain_id"`
	Maker         string     `json:"maker"`
	SourceValue   uint64     `json:"source_value"`
	TargetValue   uint64     `json:"target_value"`
	Status        int        `json:"status"`
	StatusText    string     `json:"status_text"`
	TxHash        string     `json:"tx_hash"`
}

func (this *MakerOrder) TableName() string {
	return "maker_orders"
}

func (this *DataBaseAccessObject) CreateMakerOrder(makerOrder *MakerOrder) (uint, error) {
	err := this.db.Create(makerOrder).Error
	if err != nil {
		return 0, err
	}
	return makerOrder.ID, nil
}

func (this *DataBaseAccessObject) UpdateMakerOrderStatus(id uint, status int) error {
	text := "失败"
	if status == 1 {
		text = "成功"
	}
	err := this.db.Model(&MakerOrder{}).Where("id=?", id).
		Updates(MakerOrder{Status: status, StatusText: text}).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *DataBaseAccessObject) ListMakerOrder() ([]MakerOrder, error) {
	result := make([]MakerOrder, 0)
	err := this.db.Table((&MakerOrder{}).TableName()).Find(&result).Order("id desc").
		Error
	return result, err
}

func (this *DataBaseAccessObject) ListMakerOrderByStatus(status int) ([]MakerOrder, error) {
	result := make([]MakerOrder, 0)
	err := this.db.Table((&MakerOrder{}).TableName()).
		Where("status=?", status).
		Order("id desc").
		Find(&result).Error
	return result, err
}
