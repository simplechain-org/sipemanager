package dao

import "github.com/jinzhu/gorm"

type TxAnchors struct {
	gorm.Model
	From          string `gorm:"column:from"` //锚定节点地址
	To            string `gorm:"column:to"`   //合约地址
	SourceChainId uint   `gorm:"source_chain_id"`
	TargetChainId uint   `gorm:"target_chain_id"`
	AnchorId      uint   `gorm:"column:anchor_id"`
	Fee           uint64 `gorm:"column:fee"`
}

func (this *TxAnchors) TableName() string {
	return "tx_anchors"
}

func (this *DataBaseAccessObject) CreateTxAnchors(txAnchors *TxAnchors) (uint, error) {
	err := this.db.Create(txAnchors).Error
	if err != nil {
		return 0, err
	}
	return txAnchors.ID, nil
}
