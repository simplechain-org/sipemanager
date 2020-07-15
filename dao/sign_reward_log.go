package dao

import (
	"github.com/jinzhu/gorm"
	"math/big"
)

//签名奖励发放
type SignRewardLog struct {
	gorm.Model
	AnchorNodeId    uint   `gorm:"anchor_node_id"`   //锚定节点id
	AnchorNodeName  string `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TotalReward     string `gorm:"total_reward"`     //奖励池总额
	Rate            string `gorm:"rate"`             //签名量占比,存一个格式化后的结果
	Reward          string `gorm:"reward"`           //奖励值
	TransactionHash string `gorm:"transaction_hash"` //交易哈希
	Coin            string `gorm:"coin"`             //报销的币种
	Sender          string `gorm:"sender"`           //出账账户地址
	Status          uint   `gorm:"status"`           //状态
	//计算时，链上获取到最新的签名数-数据库中最新的记录，将得到本地签名数
	FinishCount uint64 `gorm:"finish_count"` //记录发放时的签名数finish
	SignCount   uint64 `gorm:"sign_count"`   //记录发放时的签名数sign
}

func (this *SignRewardLog) TableName() string {
	return "sign_reward_logs"
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
func (this *DataBaseAccessObject) GetSignRewardLogPage(start, pageSize int, anchorNodeId uint) ([]*SignRewardLog, error) {
	result := make([]*SignRewardLog, 0)
	db := this.db.Table((&SignRewardLog{}).TableName()).
		Where("status=?", 1)
	if anchorNodeId != 0 {
		db = db.Where("anchor_node_id=?", anchorNodeId)
	}
	err := db.Offset(start).
		Limit(pageSize).
		Find(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetSignRewardLogCount(anchorNodeId uint) (int, error) {
	var count int

	db := this.db.Table((&SignRewardLog{}).TableName()).Where("status=?", 1)

	if anchorNodeId != 0 {
		db = db.Where("anchor_node_id=?", anchorNodeId)
	}
	err := db.Count(&count).Error //表示已经成功上链的数据

	return count, err
}

func (this *DataBaseAccessObject) GetSignRewardLogSumFee(anchorNodeId uint, coin string) (*big.Int, error) {
	sum:=big.NewInt(0)
	result:=make([]SignRewardLog,0)
	err := this.db.Table((&SignRewardLog{}).TableName()).
		Where("anchor_node_id=?", anchorNodeId).
		Where("coin=?", coin).
		Where("status=?", 1).Find(&result).Error
	if err!=nil{
		return nil,err
	}
	for _,o:=range result{
		fee,success:=big.NewInt(0).SetString(o.Reward,10)
		if !success{
			continue
		}
		sum=sum.Add(sum,fee)
	}
	return sum, nil
}
