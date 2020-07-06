package dao

import (
	"math/big"

	"github.com/jinzhu/gorm"
)

//报销手续费日志
type ServiceChargeLog struct {
	gorm.Model
	AnchorNodeId    uint     `gorm:"anchor_node_id"`   //锚定节点编号
	AnchorNodeName  string   `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TransactionHash string   `gorm:"transaction_hash"` //交易哈希
	Fee             *big.Int `gorm:"fee"`              //报销手续费
	Coin            string   `gorm:"coin"`             //报销的币种
	Sender          string   `gorm:"sender"`           //出账账户地址
	BlockNumber     uint     `gorm:"block_number"`     //区块高度
	Status          uint     `gorm:"status"`           //状态
}

func (this *ServiceChargeLog) TableName() string {
	return "service_charge_logs"
}

func (this *DataBaseAccessObject) CreateServiceChargeLog(obj *ServiceChargeLog) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}
func (this *DataBaseAccessObject) UpdateServiceChargeLogSourceStatus(id uint, status uint) error {
	return this.db.Table((&ServiceChargeLog{}).TableName()).
		Where("id=?", id).
		Update("status", status).Error
}
