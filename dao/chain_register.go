package dao

import "github.com/jinzhu/gorm"

type ChainRegister struct {
	gorm.Model
	SourceChainId   uint   `json:"source_chain_id"`
	TargetChainId   uint   `json:"target_chain_id"`
	Confirm         uint   `json:"confirm"`
	AnchorAddresses string `json:"anchor_addresses"`
	Status          int    `json:"status"`
	StatusText      string `json:"status_text"`
	TxHash          string `json:"tx_hash"`
}

func (this *ChainRegister) TableName() string {
	return "chain_registers"
}

func (this *DataBaseAccessObject) CreateChainRegister(chainRegister *ChainRegister) (uint, error) {
	err := this.db.Create(chainRegister).Error
	if err != nil {
		return 0, err
	}
	return chainRegister.ID, nil
}
func (this *DataBaseAccessObject) UpdateChainRegisterStatus(id uint, status int) error {
	text := "失败"
	if status == 1 {
		text = "成功"
	}
	err := this.db.Model(&ChainRegister{}).Where("id=?", id).
		Updates(ChainRegister{Status: status, StatusText: text}).Error
	if err != nil {
		return err
	}
	return nil
}

//根据当前所在的链和目标的网络id获取对应的链id
func (this *DataBaseAccessObject) GetTargetChainId(sourceChainId uint, targetNetworkId uint64) (uint, error) {
	type Result struct {
		TargetChainId uint
	}
	var result Result
	err := this.db.Model(&ChainRegister{}).
		Where("chain_registers.source_chain_id=?", sourceChainId).
		Where("chains.network_id=?", targetNetworkId).
		Joins("inner join chains on chain_registers.target_chain_id=chains.id").
		Select("chain_registers.target_chain_id").
		Order("chain_registers.id desc").
		Limit(1).Scan(&result).Error
	return result.TargetChainId, err
}

func (this *DataBaseAccessObject) ListChainRegister() ([]ChainRegister, error) {
	result := make([]ChainRegister, 0)
	err := this.db.Table((&ChainRegister{}).TableName()).Order("id desc").Find(&result).Error
	return result, err
}

func (this *DataBaseAccessObject) ListChainRegisterByStatus(status int) ([]ChainRegister, error) {
	result := make([]ChainRegister, 0)
	err := this.db.Table((&ChainRegister{}).TableName()).
		Where("status=?", status).
		Order("id desc").
		Find(&result).Error
	return result, err
}
