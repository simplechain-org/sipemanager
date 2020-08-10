package dao

import (
	"fmt"
	"math/big"
	"time"
)

//签名奖励发放
type SignRewardLog struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	AnchorNodeId    uint       `gorm:"anchor_node_id"`   //锚定节点id
	AnchorNodeName  string     `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TotalReward     string     `gorm:"total_reward"`     //奖励池总额
	Rate            string     `gorm:"rate"`             //签名量占比,存一个格式化后的结果
	Reward          string     `gorm:"reward"`           //奖励值
	TransactionHash string     `gorm:"transaction_hash"` //交易哈希
	Coin            string     `gorm:"coin"`             //报销的币种
	Sender          string     `gorm:"sender"`           //出账账户地址
	Status          uint       `gorm:"status"`           //状态
	//计算时，链上获取到最新的签名数-数据库中最新的记录，将得到本地签名数
	FinishCount uint64 `gorm:"finish_count"` //记录发放时的签名数finish
	SignCount   uint64 `gorm:"sign_count"`   //记录发放时的签名数sign
}
type SignRewardLogView struct {
	ID              uint   `gorm:"id" json:"id"`
	CreatedAt       string `gorm:"created_at" json:"created_at"`
	AnchorNodeId    uint   `gorm:"anchor_node_id" json:"anchor_node_id"`     //锚定节点id
	AnchorNodeName  string `gorm:"anchor_node_name" json:"anchor_node_name"` //锚定节点名称
	TotalReward     string `gorm:"total_reward" json:"total_reward"`         //奖励池总额
	Rate            string `gorm:"rate" json:"rate"`                         //签名量占比,存一个格式化后的结果
	Reward          string `gorm:"reward" json:"reward"`                     //奖励值
	TransactionHash string `gorm:"transaction_hash" json:"transaction_hash"` //交易哈希
	Coin            string `gorm:"coin" json:"coin"`                         //报销的币种
	Sender          string `gorm:"sender" json:"sender"`                     //出账账户地址
	Status          uint   `gorm:"status" json:"status"`                     //状态
	//计算时，链上获取到最新的签名数-数据库中最新的记录，将得到本地签名数
	FinishCount uint64 `gorm:"finish_count" json:"finish_count"` //记录发放时的签名数finish
	SignCount   uint64 `gorm:"sign_count" json:"sign_count"`     //记录发放时的签名数sign
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
func (this *DataBaseAccessObject) GetSignRewardLogPage(start, pageSize int, anchorNodeId uint) ([]*SignRewardLogView, error) {
	result := make([]*SignRewardLogView, 0)
	sql := `select id,
    (select name from anchor_nodes where id=sign_reward_logs.anchor_node_id) as anchor_node_name,
     date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    anchor_node_id,
    transaction_hash,
    total_reward,
    rate,
    reward,
    coin,
    sender,
    finish_count,
    sign_count,
    status from sign_reward_logs where status=1 and deleted_at is null`
	if anchorNodeId != 0 {
		sql += fmt.Sprintf(" and anchor_node_id=%d", anchorNodeId)
	}
	db := this.db.Raw(sql)
	err := db.Offset(start).Limit(pageSize).Scan(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetSignRewardLogCount(anchorNodeId uint) (int, error) {
	sql := `select count(*) as total from sign_reward_logs where status=1 and deleted_at is null`
	if anchorNodeId != 0 {
		sql += fmt.Sprintf(" and anchor_node_id=%d", anchorNodeId)
	}
	var total Total
	err := this.db.Raw(sql).Scan(&total).Error
	return total.Total, err
}

func (this *DataBaseAccessObject) GetSignRewardLogSumFee(anchorNodeId uint, coin string) (*big.Int, error) {
	sum := big.NewInt(0)
	result := make([]SignRewardLog, 0)
	err := this.db.Table((&SignRewardLog{}).TableName()).
		Where("anchor_node_id=?", anchorNodeId).
		Where("coin=?", coin).
		Where("status=?", 1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	for _, o := range result {
		fee, success := big.NewInt(0).SetString(o.Reward, 10)
		if !success {
			continue
		}
		sum = sum.Add(sum, fee)
	}
	return sum, nil
}
