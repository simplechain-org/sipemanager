package dao

import (
	"time"
)

type WorkCount struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId       uint `gorm:"source_chain_id" json:"source_chain_id"`
	TargetChainId       uint `gorm:"target_chain_id" json:"target_chain_id"`
	AnchorNodeId        uint `gorm:"anchor_node_id" json:"anchor_node_id"` //锚定节点编号
	SignCount           uint64 `gorm:"sign_count" json:"sign_count"`
	FinishCount         uint64 `gorm:"finish_count" json:"finish_count"`
	PreviousSignCount   uint64 `gorm:"previous_sign_count" json:"previous_sign_count"`
	PreviousFinishCount uint64 `gorm:"previous_finish_count" json:"previous_finish_count"`
}

func (this *WorkCount) TableName() string {
	return "work_counts"
}
func (this *DataBaseAccessObject) CreateWorkCount(obj *WorkCount) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}
func (this *DataBaseAccessObject) GetlatestWorkCount(sourceChainId, targetChainId, anchorNodeId uint) (*WorkCount, error) {
	var result WorkCount
	err := this.db.Table((&WorkCount{}).TableName()).
		Where("source_chain_id=?", sourceChainId).
		Where("target_chain_id=?", targetChainId).
		Where("anchor_node_id=?", anchorNodeId).Order("id desc", true).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
