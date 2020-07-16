package dao

import "github.com/jinzhu/gorm"

type Contract struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Sol  string `gorm:"type:text" json:"sol" binding:"required"`
	Abi  string `gorm:"type:text" json:"abi" binding:"required"`
	Bin  string `gorm:"type:text" json:"bin" binding:"required"`
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
	err := this.db.Table((&Contract{}).TableName()).Select("id,name").Find(&contracts).Error
	return contracts, err
}

//加载整个合约（sol,bin,abi）
func (this *DataBaseAccessObject) GetContractById(id uint) (*Contract, error) {
	var contract Contract
	err := this.db.Table((&Contract{}).TableName()).Where("id=?", id).
		First(&contract).Error
	return &contract, err

}
func (this *DataBaseAccessObject) UpdateContract(id uint, name string, sol string, abi string, bin string) error {
	return this.db.Table((&Contract{}).TableName()).
		Where("id=?", id).
		Updates(Contract{Name: name, Sol: sol, Abi: abi, Bin: bin}).Error
}

//判断合约是否可以删除
//合约还在使用，不可以删除
func (this *DataBaseAccessObject) ContractCanDelete(contractId uint) (bool, error) {
	var count int
	err := this.db.Table((&Chain{}).TableName()).
		Joins("inner join contract_instances on contract_instances.id=chains.contract_instance_id").
		Where("contract_instances.contract_id=?", contractId).
		Count(&count).Error
	return !(count > 0), err
}
func (this *DataBaseAccessObject) RemoveContract(contractId uint) error {
	return this.db.Where("id = ?", contractId).Delete(&Contract{}).Error
}

func (this *DataBaseAccessObject) GetContractPage(start, pageSize int) ([]*Contract, error) {
	result := make([]*Contract, 0)
	db := this.db.Table((&Contract{}).TableName()).
	Select("id,name")
	err := db.Offset(start).
		Limit(pageSize).
		Find(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetContractCount() (int, error) {
	var count int
	err := this.db.Table((&Contract{}).TableName()).Count(&count).Error
	return count, err
}
