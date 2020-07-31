package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type Node struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	Address   string     `gorm:"size:255" json:"address" binding:"required"` //地址
	Port      int        `json:"port" binding:"required"`                    //端口
	IsHttps   bool       `json:"is_https"`
	UserId    uint       `json:"user_id"`
	Name      string     `gorm:"size:255" json:"name" binding:"required"`
	ChainId   uint       `json:"chain_id" binding:"required"` //链id
}

type NodeView struct {
	ID        uint   `gorm:"id" json:"id"`
	CreatedAt string `gorm:"created_at" json:"created_at"`
	Address   string `gorm:"size:255" json:"address"`             //地址
	Port      int    `gorm:"port" json:"port" binding:"required"` //端口
	Name      string `gorm:"size:255" json:"name" `
	ChainId   uint   `gorm:"chain_id" json:"chain_id"` //链id
	ChainName string `gorm:"chain_name" json:"chain_name"`
	NetworkId uint64 `gorm:"network_id" json:"network_id" binding:"required"` //链的网络编号
	CoinName  string `gorm:"coin_name" json:"coin_name" binding:"required"`   //币名
	Symbol    string `gorm:"symbol" json:"symbol" binding:"required"`
}

func (this *Node) TableName() string {
	return "nodes"
}

var nodeTableName = (&Node{}).TableName()

func (this *DataBaseAccessObject) CreateNode(node *Node) (uint, error) {
	var count int
	err := this.db.Table(nodeTableName).Where("address=?", node.Address).
		Where("port=?", node.Port).
		Where("created_at is null").
		Where("user_id=?", node.UserId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("记录已经存在")
	}
	err = this.db.Create(node).Error
	if err != nil {
		return 0, err
	}
	return node.ID, nil
}
func (this *DataBaseAccessObject) ListAllNode() ([]Node, error) {
	nodes := make([]Node, 0)
	err := this.db.Table(nodeTableName).Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
func (this *DataBaseAccessObject) GetNodeById(nodeId uint) (*Node, error) {
	var node Node
	err := this.db.Table(nodeTableName).Where("id=?", nodeId).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (this *DataBaseAccessObject) ListNodeByUserId(userId uint) ([]NodeView, error) {
	nodes := make([]NodeView, 0)
	sql := `select 
    date_format(nodes.created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    nodes.id,
    nodes.address,
    nodes.port,
    nodes.name,
    nodes.chain_id,
    chains.name as chain_name,
    chains.network_id, 
    chains.coin_name,
    chains.symbol from nodes,chains where nodes.chain_id=chains.id`
	sql += " and `nodes`.`deleted_at` IS NULL"
	sql += fmt.Sprintf(" and user_id=%d", userId)
	err := this.db.Raw(sql).Scan(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
//func (this *DataBaseAccessObject) UserHasNode(userId uint) bool {
//	var count int
//	err := this.db.Table(userNodeTableName).Where("user_id=?", userId).Count(&count).Error
//	if err != nil {
//		return false
//	}
//	return count > 0
//}
func (this *DataBaseAccessObject) GetNodeByUserIdAndNetworkId(userId uint, networkId uint64) (*Node, error) {
	var node Node
	err := this.db.Table(nodeTableName).Where("user_id=? and network_id=?", userId, networkId).Order("id desc").First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (this *DataBaseAccessObject) GetNodeByChainId(chainId uint) (*Node, error) {
	var node Node
	err := this.db.Table(nodeTableName).Where("chain_id=?", chainId).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}
func (this *DataBaseAccessObject) UpdateNode(id uint, address string, port int, name string, chainId uint) error {
	return this.db.Table(nodeTableName).
		Where("id=?", id).
		Updates(Node{Address: address, Port: port, Name: name, ChainId: chainId}).Error

}
func (this *DataBaseAccessObject) RemoveNode(nodeId uint) error {
	return this.db.Where("id = ?", nodeId).Delete(&Node{}).Error
}

//根据链id获取节点列表
func (this *DataBaseAccessObject) ListNodeByChainId(chainId uint) ([]Node, error) {
	nodes := make([]Node, 0)
	err := this.db.Table(nodeTableName).Where("chain_id=?", chainId).Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
func (this *DataBaseAccessObject) ChainIdExists(chainId uint) bool {
	var total Total
	sql := `select count(*) as total from nodes where chain_id=? and deleted_at is null`
	db := this.db.Raw(sql, chainId)
	err := db.Scan(&total).Error
	if err != nil {
		fmt.Println("ChainIdExists error", err)
	}
	return total.Total > 0
}

func (this *DataBaseAccessObject) ListNodeByUserIdPage(start, pageSize int, userId uint) ([]NodeView, error) {
	nodes := make([]NodeView, 0)
	sql := `select 
    date_format(nodes.created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    nodes.id,
    nodes.address,
    nodes.port,
    nodes.name,
    nodes.chain_id,
    chains.name as chain_name,
    chains.network_id, 
    chains.coin_name,
    chains.symbol from nodes,chains where nodes.chain_id=chains.id`
	sql += " and `nodes`.`deleted_at` IS NULL"
	sql += fmt.Sprintf(" and user_id=%d", userId)
	err := this.db.Raw(sql).Offset(start).Limit(pageSize).Scan(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (this *DataBaseAccessObject) GetNodeByUserIdCount(userId uint) (int, error) {
	var total Total
	sql := `select count(*) as total from nodes,chains where nodes.chain_id=chains.id`
	sql += " and `nodes`.`deleted_at` IS NULL"
	sql += fmt.Sprintf(" and user_id=%d", userId)
	err := this.db.Raw(sql).Scan(&total).Error
	return total.Total, err
}
