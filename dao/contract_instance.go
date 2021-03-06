package dao

import (
	"time"
)

//合约实例
type ContractInstance struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	ChainId    uint       `gorm:"chain_id" json:"chain_id"` //链id ,合约部署在那条链上
	TxHash     string     `gorm:"tx_hash" json:"tx_hash"`
	Address    string     `gorm:"address" json:"address"`
	ContractId uint       `gorm:"contract_id" json:"contract_id"` //合约id
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
	ChainId    uint   `json:"chain_id"` //链id ,合约部署在那条链上
	TxHash     string `json:"tx_hash"`
	Address    string `json:"address"`
	ContractId uint   `json:"contract_id"` //合约id
	Name       string `json:"name"`
	Sol        string `gorm:"type:text" json:"sol"`
	Abi        string `gorm:"type:text" json:"abi"`
	Bin        string `gorm:"type:text" json:"bin"`
}

//链使用的当前合约
func (this *DataBaseAccessObject) GetContractByChainId(chainId uint) (*CurrentContract, error) {
	var contractInstance CurrentContract
	err := this.db.Table((&ContractInstance{}).TableName()).Where("contract_instances.chain_id=?", chainId).
		Joins("inner join contracts  on contracts.id = contract_instances.contract_id").
		Select("contract_instances.chain_id,contract_instances.tx_hash,contract_instances.address,contract_instances.contract_id,contracts.name,contracts.sol,contracts.abi,contracts.bin").
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

type InstanceNodes struct {
	CrossAddress       string `json:"cross_address"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	IsHttps            bool   `json:"is_https"`
	NetworkId          uint64 `json:"network_id"`
	Name               string `json:"name"`
	ChainId            uint   `json:"chain_id"`
	ContractId         uint   `json:"contract_id"`
	NodeId             uint   `json:"node_id"`
	ContractInstanceId uint   `json:"contract_instance_id"`
}

func (this *DataBaseAccessObject) GetInstancesJoinNode() ([]InstanceNodes, error) {
	insNodes := make([]InstanceNodes, 0)
	//var sql = `SELECT  t.cross_address,t.contract_id, n.address, n.port, n.is_https, n.network_id, n.name, n.chain_id
	//			from
	//			(select address cross_address, chain_id, contract_id from
	//				contract_instances
	//				WHERE id in
	//					(SELECT contract_instance_id id from chain_contracts )
	//					and
	//					deleted_at is null
	//			) t
	//			LEFT JOIN nodes n on n.chain_id = t.chain_id`
	var sql = `
SELECT  t.cross_address,t.contract_id, n.address, n.port, n.is_https, t.network_id, n.name, n.chain_id, n.id nodeId, t.contract_instance_id contractInstanceId  from (
	SELECT chain_id id, contract_instance_id, contract_instances.address cross_address, contract_instances.contract_id contract_id, network_id
	from 
	chains
	INNER JOIN contract_instances on chains.contract_instance_id = contract_instances.id and chains.deleted_at is null
) t 
LEFT JOIN nodes n on t.id = n.chain_id and n.deleted_at is null
`
	var result InstanceNodes
	rows, err := this.db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&result.CrossAddress,
			&result.ContractId,
			&result.Address,
			&result.Port,
			&result.IsHttps,
			&result.NetworkId,
			&result.Name,
			&result.ChainId,
			&result.NodeId,
			&result.ContractInstanceId)
		insNodes = append(insNodes, result)
	}
	return insNodes, err
}

type ContractInstanceView struct {
	ID         uint   `gorm:"id" json:"id"`
	CreatedAt  string `gorm:"created_at" json:"created_at"`
	ChainId    uint   `gorm:"chain_id" json:"chain_id"` //链id ,合约部署在那条链上
	TxHash     string `gorm:"tx_hash" json:"tx_hash"`
	Address    string `gorm:"address" json:"address"`
	ContractId uint   `gorm:"contract_id" json:"contract_id"` //合约id
	Name       string `gorm:"name" json:"name"`
	ChainName  string `gorm:"chain_name" json:"chain_name"`
}

func (this *DataBaseAccessObject) GetContractInstancePage(start, pageSize int) ([]*ContractInstanceView, error) {
	result := make([]*ContractInstanceView, 0)
	sql := `select 
			contract_instances.id,
			contract_instances.chain_id,
			contract_instances.tx_hash,
			contract_instances.address,
            date_format(contract_instances.created_at,'%Y-%m-%d %H:%i:%S') as created_at,
            (select name from chains where chains.id=contract_instances.chain_id) as chain_name,
			contracts.name from contract_instances,contracts 
            where contract_instances.contract_id=contracts.id and contract_instances.deleted_at is null`

	db := this.db.Raw(sql)

	err := db.Offset(start).Limit(pageSize).Scan(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetContractInstanceCount() (int, error) {
	var total Total
	sql := `select count(*) as total from (select 
			contract_instances.id,
			contract_instances.chain_id,
			contract_instances.tx_hash,
			contract_instances.address,
            (select name from chains where chains.id=contract_instances.chain_id) as chain_name,
			contracts.name from contract_instances ,contracts where contract_instances.contract_id=contracts.id and contract_instances.deleted_at is null) as temp`
	db := this.db.Raw(sql)
	err := db.Scan(&total).Error
	return total.Total, err
}

func (this *DataBaseAccessObject) GetContractInstanceById(id uint) (*ContractInstance, error) {
	var contractInstance ContractInstance
	err := this.db.Table((&ContractInstance{}).TableName()).Where("id=?", id).
		First(&contractInstance).Error
	return &contractInstance, err

}
