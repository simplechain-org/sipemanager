package dao

import "github.com/jinzhu/gorm"

type Chain struct {
	gorm.Model
	Name               string `json:"name" binding:"required"`       //链的名称
	NetworkId          uint64 `json:"network_id" binding:"required"` //链的网络编号
	CoinName           string `json:"coin_name" binding:"required"`  //币名
	Symbol             string `json:"symbol" binding:"required"`     //符号
	ContractInstanceId uint   `gorm:"contract_instance_id"`          //合约实例
}

func (this *Chain) TableName() string {
	return "chains"
}

func (this *DataBaseAccessObject) CreateChain(chain *Chain) (uint, error) {
	err := this.db.Create(chain).Error
	if err != nil {
		return 0, err
	}
	return chain.ID, nil
}
func (this *DataBaseAccessObject) GetChain(chainId uint) (*Chain, error) {
	var chain Chain
	err := this.db.Table((&Chain{}).TableName()).Where("id=?", chainId).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}
func (this *DataBaseAccessObject) GetChainByNetWorkId(NetWorkId uint64) (*Chain, error) {
	var chain Chain
	err := this.db.Table((&Chain{}).TableName()).Where("network_id=?", NetWorkId).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}
func (this *DataBaseAccessObject) GetChains() ([]*Chain, error) {
	chains := make([]*Chain, 0)
	err := this.db.Table((&Chain{}).TableName()).Find(&chains).Error
	if err != nil {
		return nil, err
	}
	return chains, nil
}

//链上是否还存在节点
func (this *DataBaseAccessObject) ChainHasNode(chainId uint) bool {
	var count int
	err := this.db.Table((&Node{}).TableName()).Where("chain_id=?", chainId).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

//链上是否还存在合约实例
func (this *DataBaseAccessObject) ChainHasContractInstance(chainId uint) bool {
	var count int
	err := this.db.Table((&ContractInstance{}).TableName()).Where("chain_id=?", chainId).Count(&count).Error
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func (this *DataBaseAccessObject) ChainRemove(chainId uint) error {
	return this.db.Where("id = ?", chainId).Delete(&Chain{}).Error
}

//根据合约地址获取链id
func (this *DataBaseAccessObject) GetChainIdByContractAddress(address string) (uint, error) {
	//读取单个字段时，需要创建一个结构体，因为scan的参数是slice or struct
	type Result struct {
		ChainId uint
	}
	var result Result
	err := this.db.Table((&ContractInstance{}).TableName()).Where("address=?", address).
		Order("id desc").Limit(1).Scan(&result).Error
	return result.ChainId, err
}

func (this *DataBaseAccessObject) GetChainIdByNetworkId(networkId uint64) (uint, error) {
	//读取单个字段时，需要创建一个结构体，因为scan的参数是slice or struct
	type Result struct {
		Id uint
	}
	var result Result
	err := this.db.Table((&Chain{}).TableName()).Where("network_id=?", networkId).
		Order("id desc").Limit(1).Scan(&result).Error
	return result.Id, err
}

func (this *DataBaseAccessObject) UpdateChain(id uint, name string, networkId uint64, coinName string, symbol string, contractInstanceId uint) error {
	var chain *Chain
	chain = &Chain{
		Name:      name,
		NetworkId: networkId,
		CoinName:  coinName,
		Symbol:    symbol,
	}
	if contractInstanceId != 0 {
		chain = &Chain{
			ContractInstanceId: contractInstanceId,
			Name:               name,
			NetworkId:          networkId,
			CoinName:           coinName,
			Symbol:             symbol,
		}
	}
	return this.db.Table((&Chain{}).TableName()).
		Where("id=?", id).
		Updates(chain).Error
}

type ChainInfo struct {
	gorm.Model
	Name               string `json:"name" gorm:"name"`                                 //链的名称
	NetworkId          uint64 `json:"network_id" gorm:"network_id"`                     //链的网络编号
	CoinName           string `json:"coin_name" gorm:"coin_name"`                       //币名
	Symbol             string `json:"symbol" gorm:"coin_name"`                          //符号
	ContractInstanceId uint   `json:"contract_instance_id" gorm:"contract_instance_id"` //合约实例
	Address            string `json:"address" gorm:"address"`                           //合约地址
}

func (this *DataBaseAccessObject) GetChainInfoPage(start, pageSize int) ([]*ChainInfo, error) {
	result := make([]*ChainInfo, 0)
	db := this.db.Table((&Chain{}).TableName()).Joins("left join contract_instances on contract_instances.id=chains.contract_instance_id").
		Select("chains.id,chains.name,chains.network_id,chains.coin_name,chains.symbol,chains.contract_instance_id,contract_instances.address")
	err := db.Offset(start).
		Limit(pageSize).
		Find(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetChainInfoCount() (int, error) {
	var count int
	err := this.db.Table((&Chain{}).TableName()).Count(&count).Error
	return count, err
}
func (this *DataBaseAccessObject) GetChainInfo(chainId uint) (*ChainInfo, error) {
	var chain ChainInfo
	err := this.db.Table((&Chain{}).TableName()).Joins("left join contract_instances on contract_instances.id=chains.contract_instance_id").
		Select("chains.id,chains.name,chains.network_id,chains.coin_name,chains.symbol,chains.contract_instance_id,contract_instances.address").Where("id=?", chainId).First(&chain).Error
	if err != nil {
		return nil, err
	}
	return &chain, nil
}
