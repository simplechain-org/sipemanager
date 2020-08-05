package dao

import (
	"time"
)

//接单记录
type TakerOrder struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId uint       `json:"source_chain_id"` //这是链记录的id不是networkId
	TargetChainId uint       `json:"target_chain_id"`
	Taker         string     `json:"taker"`
	SourceValue   uint64     `json:"source_value"`
	TargetValue   uint64     `json:"target_value"`
	Status        int        `json:"status"`
	StatusText    string     `json:"status_text"`
	TxHash        string     `json:"tx_hash"`
	CtxId         string     `json:"ctx_id"`
}

func (this *TakerOrder) TableName() string {
	return "taker_orders"
}

func (this *DataBaseAccessObject) CreateTakerOrder(takerOrder *TakerOrder) (uint, error) {
	err := this.db.Create(takerOrder).Error
	if err != nil {
		return 0, err
	}
	return takerOrder.ID, nil
}
func (this *DataBaseAccessObject) UpdateTakerOrderStatus(id uint, status int) error {
	text := "失败"
	if status == 1 {
		text = "成功"
	}
	err := this.db.Model(&TakerOrder{}).Where("id=?", id).
		Updates(TakerOrder{Status: status, StatusText: text}).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *DataBaseAccessObject) ListTakerOrder() ([]TakerOrder, error) {
	result := make([]TakerOrder, 0)
	err := this.db.Table((&TakerOrder{}).TableName()).Order("id desc").Find(&result).Error
	return result, err
}

func (this *DataBaseAccessObject) ListTakerOrderByStatus(status int) ([]TakerOrder, error) {
	result := make([]TakerOrder, 0)
	err := this.db.Table((&TakerOrder{}).TableName()).
		Where("status=?", status).
		Order("id desc").
		Find(&result).Error
	return result, err
}
