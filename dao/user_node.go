package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type UserNode struct {
	gorm.Model
	UserId uint
	NodeId uint
}

func (this *UserNode) TableName() string {
	return "user_nodes"
}

var userNodeTableName = (&UserNode{}).TableName()

func (this *DataBaseAccessObject) CreateUserNode(userNode *UserNode) (uint, error) {
	var count int
	err := this.db.Table(userNodeTableName).Where("user_id=?", userNode.UserId).
		Where("node_id=?", userNode.NodeId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("Record already exists")
	}
	//todo 判断节点是否存在
	_, err = this.GetNodeById(userNode.NodeId)
	if err != nil {
		return 0, err
	}
	err = this.db.Create(userNode).Error
	if err != nil {
		return 0, err
	}
	return userNode.ID, nil
}

func (this *DataBaseAccessObject) GetUserCurrentNode(userId uint) (*Node, error) {
	var userNode UserNode
	err := this.db.Table(userNodeTableName).Where("user_id=?", userId).First(&userNode).Error
	if err != nil {
		return nil, err
	}
	node, err := this.GetNodeById(userNode.NodeId)
	if err != nil {
		return nil, err
	}
	return node, nil
}
func (this *DataBaseAccessObject) UpdateUserCurrentNode(userNode *UserNode) error {
	var count int
	err := this.db.Table(userNodeTableName).Where("user_id=?", userNode.UserId).
		Where("node_id=?", userNode.NodeId).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = this.GetNodeById(userNode.NodeId)
	if err != nil {
		return err
	}
	return this.db.Model(&UserNode{}).
		Where("user_id = ?", userNode.UserId).
		Update("node_id", userNode.NodeId).Error
}
