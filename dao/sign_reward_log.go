package dao

import (
	"math/big"

	"github.com/jinzhu/gorm"
)

//签名奖励发放
type SignRewardLog struct {
	gorm.Model
	AnchorNodeId    uint     `gorm:"anchor_node_id"`   //锚定节点id
	AnchorNodeName  string   `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TotalReward     *big.Int `gorm:"total_reward"`     //奖励池总额
	Rate            int      `gorm:"rate"`             //签名量占比
	Reward          *big.Int `gorm:"reward"`           //奖励值
	TransactionHash string   `gorm:"transaction_hash"` //交易哈希
	Coin            string   `gorm:"coin"`             //报销的币种
	Sender          string   `gorm:"sender"`           //出账账户地址
	BlockNumber     uint     `gorm:"block_number"`     //区块高度
	Status          uint     `gorm:"status"`           //状态
	Signatures      uint     `gorm:"signatures"`       //签名数
}

func (this *SignRewardLog) TableName() string {
	return "service_charge_logs"
}

func (this *DataBaseAccessObject) CreateSignRewardLog(obj *SignRewardLog) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}
func (this *DataBaseAccessObject) UpdateSignRewardLogStatus(id uint, status uint) error {
	return this.db.Table((&SignRewardLog{}).TableName()).
		Where("id=?", id).
		Update("status", status).Error
}
