package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"math/big"
)

//锚定节点
type AnchorNode struct {
	gorm.Model
	Name          string `gorm:"name"`    //锚定节点名称
	Address       string `gorm:"address"` //锚定节点地址
	SourceChainId uint   `gorm:"source_chain_id"`
	TargetChainId uint   `gorm:"target_chain_id"`
	SourceHash    string `gorm:"source_hash"`   //链上的交易哈希
	TargetHash    string `gorm:"target_hash"`   //链上的交易哈希
	SourceStatus  uint   `gorm:"source_status"` //链上达成的状态  锚定节点添加成功
	TargetStatus  uint   `gorm:"target_status"` //链上达成的状态  锚定节点添加成功
	Pledge        string `gorm:"pledge"`        //质押sipc的金额
	Status        bool   `gorm:"status"`
}

func (this *AnchorNode) TableName() string {
	return "anchor_nodes"
}

//添加锚定节点
func (this *DataBaseAccessObject) CreateAnchorNode(instance *AnchorNode) (uint, error) {
	err := this.db.Create(instance).Error
	if err != nil {
		return 0, err
	}
	return instance.ID, nil
}
func (this *DataBaseAccessObject) CreateAnchorNodeByTx(db *gorm.DB, instance *AnchorNode) (uint, error) {
	err := db.Create(instance).Error
	if err != nil {
		return 0, err
	}
	return instance.ID, nil
}

func (this *DataBaseAccessObject) UpdateSourceStatus(id uint, status uint) error {
	return this.db.Table((&AnchorNode{}).TableName()).
		Where("id=?", id).
		Update("source_status", status).Error
}
func (this *DataBaseAccessObject) UpdateTargetStatus(id uint, status uint) error {
	return this.db.Table((&AnchorNode{}).TableName()).
		Where("id=?", id).
		Update("target_status", status).Error
}
func (this *DataBaseAccessObject) GetAnchorNode(id uint) (*AnchorNode, error) {
	var obj AnchorNode
	err := this.db.Table((&AnchorNode{}).TableName()).Where("id=?", id).First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
func (this *DataBaseAccessObject) RemoveAnchorNode(id uint) error {
	return this.db.Where("id = ?", id).Delete(&AnchorNode{}).Error

}

func (this *DataBaseAccessObject) ListAnchorNode() ([]AnchorNode, error) {
	anchorNodes := make([]AnchorNode, 0)
	err := this.db.Table((&AnchorNode{}).TableName()).Where("target_status=?", 1).
		Where("source_status=?", 1).
		Find(&anchorNodes).Error
	if err != nil {
		return nil, err
	}
	return anchorNodes, nil
}
func (this *DataBaseAccessObject) GetAnchorNodeCount() (int, error) {
	var count int
	err := this.db.Table((&AnchorNode{}).TableName()).Where("target_status=?", 1).
		Where("source_status=?", 1).
		Count(&count).Error
	return count, err
}

func (this *DataBaseAccessObject) GetAnchorNodePage(start, pageSize int) ([]*AnchorNode, error) {
	result := make([]*AnchorNode, 0)
	err := this.db.Table((&AnchorNode{}).TableName()).Where("target_status=?", 1).
		Where("source_status=?", 1).Offset(start).
		Limit(pageSize).
		Find(&result).Error
	return result, err
}

func (this *DataBaseAccessObject) SubPledge(anchorNodeId uint, value string) error {
	var obj AnchorNode
	err := this.db.Table((&AnchorNode{}).TableName()).Where("id=?", anchorNodeId).First(&obj).Error
	if err != nil {
		return err
	}
	fee, success := big.NewInt(0).SetString(value, 10)
	if !success {
		return errors.New("value数据非法")
	}
	sub, success := big.NewInt(0).SetString(obj.Pledge, 10)
	if !success {
		return errors.New("pledge数据非法")
	}
	sub = sub.Sub(sub, fee)
	return this.db.Table((&AnchorNode{}).TableName()).
		Where("id=?", anchorNodeId).
		Update("pledge", sub.String()).Error
}
