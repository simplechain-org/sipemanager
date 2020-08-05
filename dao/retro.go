package dao

import (
	"errors"
	"fmt"
	"time"
)

type RetroActive struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	NetworkId uint64     `json:"network_id"`
	Status    int        `json:"status"` //1 待签  2 签名完成
	TxHash    string     `json:"tx_hash"`
	CtxId     string     `json:"ctx_id"`
	Event     int        `json:"event"` //1 maker  2 taker
}

func (this *RetroActive) TableName() string {
	return "retro_active"
}

func (this *DataBaseAccessObject) CreateRetroActive(retroActive *RetroActive) (uint, error) {
	err := this.db.Create(retroActive).Error
	if err != nil {
		return 0, err
	}
	return retroActive.ID, nil
}
func (this *DataBaseAccessObject) UpdateRetroActiveStatus(id uint, status int) error {
	err := this.db.Model(&RetroActive{}).Where("id=?", id).
		Updates(RetroActive{Status: status}).Error
	if err != nil {
		return err
	}
	return nil
}

func (this *DataBaseAccessObject) ListRetroActive(offset, limit uint32) ([]RetroActive, error) {
	var count uint32

	if err := this.db.Model(&RetroActive{}).Count(&count).Error; err != nil {
		return nil, err
	}
	fmt.Println(count, offset, limit, this.db)
	if offset <= count {
		result := make([]RetroActive, 0)
		err := this.db.Table((&RetroActive{}).TableName()).Order("id desc").Offset(offset).Limit(limit).Find(&result).Error
		return result, err
	}
	return nil, errors.New("offset > count")
}

func (this *DataBaseAccessObject) ListRetroActiveByStatus(status int, offset, limit uint32) ([]RetroActive, error) {
	var count uint32
	if err := this.db.Model(&RetroActive{}).Where("status=?", status).Count(&count).Error; err != nil {
		return nil, err
	}
	if offset <= count {
		result := make([]RetroActive, 0)
		err := this.db.Table((&RetroActive{}).TableName()).
			Where("status=?", status).
			Order("id desc").Offset(offset).Limit(limit).
			Find(&result).Error
		return result, err
	}
	return nil, errors.New("offset > count")
}
func (this *DataBaseAccessObject) QueryRetroActive(txHash string) (*RetroActive, error) {
	var retro RetroActive
	err := this.db.Table((&RetroActive{}).TableName()).Where("tx_hash=?", txHash).First(&retro).Error
	if err != nil {
		return nil, err
	}
	return &retro, nil
}
