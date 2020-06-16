package dao

import "github.com/jinzhu/gorm"

//合约实例
type ContractInstance struct {
	gorm.Model
	ChainId    uint   `json:"chain_id"` //链id ,合约部署在那条链上
	TxHash     string `json:"tx_hash"`
	Address    string `json:"address"`
	ContractId uint   `json:"contract_id"` //合约id
}

func (this *ContractInstance) TableName() string {
	return "contract_instances"
}

//合约部署
func (this *DataBaseAccessObject) CreateContractInstance(instance *ContractInstance) (uint, error) {
	err := this.db.Create(instance).Error
	if err != nil {
		return 0, err
	}
	return instance.ID, nil
}

//合约成功部署以后，更新地址
func (this *DataBaseAccessObject) UpdateContractAddress(id uint, address string) error {
	return this.db.Table((&ContractInstance{}).TableName()).
		Where("id=?", id).
		Update("address", address).Error
}

type CurrentContract struct {
	ChainId     uint   `json:"chain_id"` //链id ,合约部署在那条链上
	TxHash      string `json:"tx_hash"`
	Address     string `json:"address"`
	ContractId  uint   `json:"contract_id"` //合约id
	Description string `json:"description"`
	Sol         string `gorm:"type:text" json:"sol"`
	Abi         string `gorm:"type:text" json:"abi"`
	Bin         string `gorm:"type:text" json:"bin"`
}

//链使用的当前合约
func (this *DataBaseAccessObject) GetContractByChainId(chainId uint) (*CurrentContract, error) {
	var contractInstance CurrentContract
	err := this.db.Table((&ContractInstance{}).TableName()).Where("contract_instances.chain_id=?", chainId).
		Joins("inner join contracts  on contracts.id = contract_instances.contract_id").
		Select("contract_instances.chain_id,contract_instances.tx_hash,contract_instances.address,contract_instances.contract_id,contracts.description,contracts.sol,contracts.abi,contracts.bin").
		Scan(&contractInstance).Error
	return &contractInstance, err

}

//链上已经部署的合约
func (this *DataBaseAccessObject) GetContractsByChainId(chainId uint) ([]*ContractInstance, error) {
	contracts := make([]*ContractInstance, 0)
	err := this.db.Table((&ContractInstance{}).TableName()).Where("chain_id=?", chainId).
		Select("id,chain_id,tx_hash,address").Find(&contracts).Error
	return contracts, err
}
func (this *DataBaseAccessObject) GetContractInstances() ([]*ContractInstance, error) {
	contracts := make([]*ContractInstance, 0)
	err := this.db.Table((&ContractInstance{}).TableName()).Select("id,chain_id,tx_hash,address").Find(&contracts).Error
	return contracts, err
}
