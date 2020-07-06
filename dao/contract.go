package dao

import "github.com/jinzhu/gorm"

type Contract struct {
	gorm.Model
	Description string `json:"description" binding:"required"`
	Sol         string `gorm:"type:text" json:"sol" binding:"required"`
	Abi         string `gorm:"type:text" json:"abi" binding:"required"`
	Bin         string `gorm:"type:text" json:"bin" binding:"required"`
}

func (this *Contract) TableName() string {
	return "contracts"
}

//合约创建
func (this *DataBaseAccessObject) CreateContract(contract *Contract) (uint, error) {
	err := this.db.Create(contract).Error
	if err != nil {
		return 0, err
	}
	return contract.ID, nil
}

func (this *DataBaseAccessObject) UpdateContractSol(id uint, sol string) error {
	return this.db.Table((&Contract{}).TableName()).
		Where("id=?", id).
		Update("sol", sol).Error
}
func (this *DataBaseAccessObject) UpdateContractAbi(id uint, abi string) error {
	return this.db.Table((&Contract{}).TableName()).
		Where("id=?", id).
		Update("abi", abi).Error
}
func (this *DataBaseAccessObject) UpdateContractBin(id uint, bin string) error {
	return this.db.Table((&Contract{}).TableName()).
		Where("id=?", id).
		Update("bin", bin).Error
}

//列出所有的合约，因为Sol，bin,abi较大，不加载
func (this *DataBaseAccessObject) GetContracts() ([]*Contract, error) {
	contracts := make([]*Contract, 0)
	err := this.db.Table((&Contract{}).TableName()).Select("id,description").Find(&contracts).Error
	return contracts, err
}

//加载整个合约（sol,bin,abi）
func (this *DataBaseAccessObject) GetContractById(id uint) (*Contract, error) {
	var contract Contract
	err := this.db.Table((&Contract{}).TableName()).Where("id=?", id).
		First(&contract).Error
	return &contract, err

}

//链使用哪个合约进行跨链
type ChainContract struct {
	ID                 uint `gorm:"primary_key"`
	ChainId            uint //链id
	ContractInstanceId uint `gorm:"contract_instance_id"` //合约实例id

}

func (this *ChainContract) TableName() string {
	return "chain_contracts"
}
