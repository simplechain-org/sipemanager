package dao

import "github.com/jinzhu/gorm"

type WorkCount struct {
	gorm.Model
	SourceChainId       uint `gorm:"source_chain_id" json:"source_chain_id"`
	TargetChainId       uint `gorm:"target_chain_id" json:"target_chain_id"`
	AnchorNodeId        uint `gorm:"anchor_node_id" json:"anchor_node_id"` //锚定节点编号
	SignCount           uint `gorm:"sign_count" json:"sign_count"`
	FinishCount         uint `gorm:"finish_count" json:"finish_count"`
	PreviousSignCount   uint `gorm:"previous_sign_count" json:"previous_sign_count"`
	PreviousFinishCount uint `gorm:"previous_finish_count" json:"previous_finish_count"`
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