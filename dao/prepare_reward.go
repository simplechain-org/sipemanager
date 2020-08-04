package dao

import (
	"time"
)

//预扣费用，针对合约中的reward
type PrepareReward struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId uint       `json:"source_chain_id"`
	TargetChainId uint       `json:"target_chain_id"`
	SourceReward  string     `json:"source_reward"`
	TargetReward  string     `json:"target_reward"`
	SourceHash    string     `json:"source_hash"`
	TargetHash    string     `json:"target_hash"`
}
type PrepareRewardView struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId   uint       `json:"source_chain_id"`
	TargetChainId   uint       `json:"target_chain_id"`
	SourceReward    string     `json:"source_reward"`
	TargetReward    string     `json:"target_reward"`
	SourceChainName string     `json:"source_chain_name"`
	TargetChainName string     `json:"target_chain_name"`
	SourceChainCoin string     `json:"source_chain_coin"`
	TargetChainCoin string     `json:"target_chain_coin"`
}

func (this *PrepareReward) TableName() string {
	return "prepare_rewards"
}

//添加预扣费用记录
func (this *DataBaseAccessObject) CreatePrepareReward(obj *PrepareReward) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}

//更新
func (this *DataBaseAccessObject) UpdatePrepareReward(obj *PrepareReward) error {
	return this.db.Table((&PrepareReward{}).TableName()).
		Where("source_chain_id=?", obj.SourceChainId).
		Where("target_chain_id=?", obj.TargetChainId).
		Updates(PrepareReward{
			SourceReward: obj.SourceReward,
			TargetReward: obj.TargetReward,
			SourceHash:   obj.SourceHash,
			TargetHash:   obj.TargetHash,
		}).Error
}
func (this *DataBaseAccessObject) GetPrepareRewardPage(start, pageSize int) ([]*PrepareRewardView, error) {
	search := make([]*PrepareRewardView, 0)
	sql := `SELECT id, 
			source_chain_id,
			target_chain_id,
			source_reward,
			target_reward,
			(select name from chains where id=source_chain_id) as source_chain_name,
			(select name  from chains where id=target_chain_id) as target_chain_name,
			(select coin_name from chains where id=source_chain_id) as source_chain_coin,
			(select coin_name from chains where id=target_chain_id) as target_chain_coin
			 FROM prepare_rewards where deleted_at is null`
	db := this.db.Raw(sql)
	//必须使用Scan，不能使用find
	err := db.Offset(start).
		Limit(pageSize).
		Scan(&search).Error
	return search, err
}

func (this *DataBaseAccessObject) GetPrepareRewardCount() (int, error) {
	var total Total
	sql := `select count(*) as total from prepare_rewards where deleted_at is null`
	db := this.db.Raw(sql)
	//必须使用Scan，不能使用find
	err := db.Scan(&total).Error
	return total.Total, err
}
