package dao

import (
	"github.com/jinzhu/gorm"
)

type RewardConfig struct {
	gorm.Model
	SourceChainId   uint   `json:"source_chain_id"`
	TargetChainId   uint   `json:"target_chain_id"`
	RegulationCycle uint   `json:"regulation_cycle"` //调控周期
	SignReward      string `json:"sign_reward"`      //单笔签名奖励
}

type RewardConfigView struct {
	CreatedAt        string `gorm:"created_at" json:"created_at"`
	ID               uint   `gorm:"id" json:"ID"`
	SourceChainId    uint   `gorm:"source_chain_id" json:"source_chain_id"`
	TargetChainId    uint   `gorm:"target_chain_id" json:"target_chain_id"`
	SourceChainName  string `json:"source_chain_name" gorm:"source_chain_name"`
	TargetChainName  string `json:"target_chain_name" gorm:"target_chain_name"`
	RegulationCycle  uint   `gorm:"regulation_cycle" json:"regulation_cycle"`   //调控周期
	SignReward       string `gorm:"sign_reward" json:"sign_reward"`             //单笔签名奖励
	InProgress       uint   `gorm:"in_progress" json:"in_progress"`             //已进行多少天，计算字段
	TransactionCount uint64 `gorm:"transaction_count" json:"transaction_count"` //本周期交易笔数，计算字段
}

func (this *RewardConfig) TableName() string {
	return "reward_configs"
}
func (this *DataBaseAccessObject) CreateRewardConfig(obj *RewardConfig) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}

func (this *DataBaseAccessObject) GetRewardConfig(id uint) (*RewardConfigView, error) {
	sql := `select id,
			date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
			source_chain_id,
			target_chain_id,
			regulation_cycle,
			sign_reward,
			(select TIMESTAMPDIFF(DAY,date_format(created_at,'%Y-%m-%d %H:%i:%S'),date_format(NOW(), '%Y-%m-%d %H:%i:%S'))) as in_progress,
			(select name from chains where chains.id=reward_configs.source_chain_id) as source_chain_name,
			(select name from chains where chains.id=reward_configs.target_chain_id) as target_chain_name,
			(select sum(sign_count+finish_count) from work_counts where source_chain_id=reward_configs.source_chain_id and target_chain_id=reward_configs.target_chain_id and created_at > reward_configs.created_at) as transaction_count
			from reward_configs where id=?`
	var rewardConfigView RewardConfigView
	err := this.db.Raw(sql, id).Scan(&rewardConfigView).Error
	return &rewardConfigView, err
}

//删除
func (this *DataBaseAccessObject) RemoveRewardConfig(id uint) error {
	return this.db.Where("id = ?", id).Delete(&RewardConfig{}).Error
}

type IndexRewardConfigView struct {
	SourceChainId uint `gorm:"source_chain_id" json:"source_chain_id"`
	TargetChainId uint `gorm:"target_chain_id" json:"target_chain_id"`
}

//分页
func (this *DataBaseAccessObject) GetRewardConfigPage(start, pageSize int) ([]*RewardConfigView, error) {
	search := make([]*IndexRewardConfigView, 0)
	sql := `SELECT DISTINCT source_chain_id,target_chain_id FROM reward_configs where deleted_at is null`
	db := this.db.Raw(sql)
	err := db.Offset(start).
		Limit(pageSize).
		Find(&search).Error
	result := make([]*RewardConfigView, 0)
	for _, o := range search {
		v, err := this.GetRewardConfigBySourceAndTarget(o.SourceChainId, o.TargetChainId)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, err
}

type Total struct {
	Total int `gorm:"total"`
}

func (this *DataBaseAccessObject) GetRewardConfigCount() (int, error) {
	var total Total
	sql := `select count(*) as total from (SELECT DISTINCT source_chain_id,target_chain_id FROM reward_configs where deleted_at is null) as reward_configs_temp`
	db := this.db.Raw(sql)
	err := db.Scan(&total).Error
	return total.Total, err
}

func (this *DataBaseAccessObject) GetLatestRewardConfig(sourceChainId, targetChainId uint) (*RewardConfig, error) {
	var result RewardConfig
	sql := `select *  from reward_configs where source_chain_id=? and target_chain_id=? and deleted_at is null`
	db := this.db.Raw(sql, sourceChainId, targetChainId).Order("id DESC", true)
	err := db.Limit(1).Scan(&result).Error
	return &result, err
}

//根据ID 获取 source_chain_id和 target_chain_id
//然后标记删除（软删除）所有的相同的source_chain_id和 target_chain_id的记录
//为了保存以往的记录而采取的做法
func (this *DataBaseAccessObject) RemoveRelativeRewardConfig(id uint) error {
	var result RewardConfig
	err := this.db.Table((&RewardConfig{}).TableName()).Where("id=?", id).First(&result).Error
	if err != nil {
		return err
	}
	return this.db.Where("source_chain_id=?", result.SourceChainId).
		Where("target_chain_id=?", result.TargetChainId).
		Delete(RewardConfig{}).Error
}

func (this *DataBaseAccessObject) GetRewardConfigBySourceAndTarget(sourceChainId, targetChainId uint) (*RewardConfigView, error) {
	sql := `select id,
			date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
			source_chain_id,
			target_chain_id,
			regulation_cycle,
			sign_reward,
			(select TIMESTAMPDIFF(DAY,date_format(created_at,'%Y-%m-%d %H:%i:%S'),date_format(NOW(), '%Y-%m-%d %H:%i:%S'))) as in_progress,
			(select name from chains where chains.id=reward_configs.source_chain_id) as source_chain_name,
			(select name from chains where chains.id=reward_configs.target_chain_id) as target_chain_name,
			(select sum(sign_count+finish_count) from work_counts where source_chain_id=reward_configs.source_chain_id and target_chain_id=reward_configs.target_chain_id and created_at > reward_configs.created_at) as transaction_count
			from reward_configs where source_chain_id=? and target_chain_id=? and deleted_at is null order by id desc limit 1`
	var rewardConfigView RewardConfigView
	err := this.db.Raw(sql, sourceChainId, targetChainId).Scan(&rewardConfigView).Error
	return &rewardConfigView, err
}
func (this *DataBaseAccessObject) GetRewardConfigById(id uint) (*RewardConfig, error) {
	var result RewardConfig
	err := this.db.Table((&RewardConfig{}).TableName()).Where("id=?", id).First(&result).Error
	return &result, err
}
