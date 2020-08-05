package dao

import (
	"fmt"
	"time"
)

type Contract struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	Name      string     `json:"name" binding:"required"`
	Sol       string     `gorm:"type:text" json:"sol" binding:"required"`
	Abi       string     `gorm:"type:text" json:"abi" binding:"required"`
	Bin       string     `gorm:"type:text" json:"bin" binding:"required"`
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

//列出所有的合约
func (this *DataBaseAccessObject) GetContracts() ([]*Contract, error) {
	contracts := make([]*Contract, 0)
	err := this.db.Table((&Contract{}).TableName()).Find(&contracts).Error
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
	contract := Contract{Name: name}
	if sol != "" {
		contract.Sol = sol
	}
	if abi != "" {
		contract.Abi = abi
	}
	if bin != "" {
		contract.Bin = bin
	}
	return this.db.Table((&Contract{}).TableName()).Where("id=?", id).Updates(contract).Error
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
func (this *DataBaseAccessObject) GetContractPage(start, pageSize int, status string) ([]*Contract, error) {
	if status == "not_deployed" || status == "deployed" {
		sql := `select id,
        name,
        sol,
        bin,
        abi,
        created_at,
        updated_at,
        deleted_at from contracts`
		if status == "deployed" {
			sql += " where id in (select distinct contract_id from contract_instances)"
		} else {
			sql += " where id not in (select distinct contract_id from contract_instances)"
		}
		sql += " and `contracts`.`deleted_at` IS NULL "
		sql += fmt.Sprintf(" order by id limit %d offset %d", pageSize, start)
		result := make([]*Contract, 0)
		err := this.db.Raw(sql).Scan(&result).Error
		return result, err

	} else {
		result := make([]*Contract, 0)
		err := this.db.Table((&Contract{}).TableName()).Offset(start).Limit(pageSize).Find(&result).Error
		return result, err
	}
}
func (this *DataBaseAccessObject) GetContractCount(status string) (int, error) {
	if status == "not_deployed" || status == "deployed" {
		sql := `select count(*) as total from contracts`
		if status == "deployed" {
			sql += " where id in (select distinct contract_id from contract_instances)"
		} else {
			sql += " where id not in (select distinct contract_id from contract_instances)"
		}
		sql += " and `contracts`.`deleted_at` IS NULL "
		var total Total
		err := this.db.Raw(sql).Scan(&total).Error
		return total.Total, err
	} else {
		var total Total
		sql := `select count(*) as total from contracts where contracts.deleted_at IS NULL`
		err := this.db.Raw(sql).Scan(&total).Error
		return total.Total, err
	}
}
