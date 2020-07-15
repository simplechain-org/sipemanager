package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Node struct {
	gorm.Model
	Address string `gorm:"size:255" json:"address" binding:"required"` //地址
	Port    int    `json:"port" binding:"required"`                    //端口
	IsHttps bool   `json:"is_https"`
	UserId  uint   `json:"user_id"`
	Name    string `gorm:"size:255" json:"name" binding:"required"`
	ChainId uint   `json:"chain_id" binding:"required"` //链id
}

type NodeView struct {
	gorm.Model
	Address   string `gorm:"size:255" json:"address"` //地址
	Port      int    `json:"port" binding:"required"` //端口
	Name      string `gorm:"size:255" json:"name" `
	ChainId   uint   `json:"chain_id"` //链id
	ChainName string `gorm:"chain_name" json:"chain_name"`
}

func (this *Node) TableName() string {
	return "nodes"
}

var nodeTableName = (&Node{}).TableName()

//{"address":"127.0.0.1","port":9545,"chain_id":1,"name":"主链节点1"}
func (this *DataBaseAccessObject) CreateNode(node *Node) (uint, error) {
	var count int
	err := this.db.Table(nodeTableName).Where("address=?", node.Address).
		Where("port=?", node.Port).
		Where("user_id=?", node.UserId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("Record already exists")
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
func (this *DataBaseAccessObject) ListNodeByUserId(userId uint) ([]Node, error) {
	nodes := make([]Node, 0)
	err := this.db.Table(nodeTableName).Where("user_id=?", userId).Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
func (this *DataBaseAccessObject) UserHasNode(userId uint) bool {
	var count int
	err := this.db.Table(userNodeTableName).Where("user_id=?", userId).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
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
