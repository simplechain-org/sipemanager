package dao

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Punishment struct {
	gorm.Model
	Value string //惩罚数量
	Coin  string //惩罚币种
	//suspend recovery token
	ManageType   string //管理类型
	AnchorNodeId uint   `gorm:"anchor_node_id"` //锚定节点编号
}

type PunishmentView struct {
	ID             uint   `gorm:"id" json:"ID"`
	CreatedAt      string `gorm:"created_at" json:"CreatedAt"`
	Value          string `gorm:"value" json:"value"`                       //惩罚数量
	Coin           string `gorm:"coin" json:"coin"`                         //惩罚币种
	ManageType     string `gorm:"manage_type" json:"manage_type"`           //管理类型
	AnchorNodeId   uint   `gorm:"anchor_node_id" json:"anchor_node_id"`     //锚定节点编号
	AnchorNodeName string `gorm:"anchor_node_name" json:"anchor_node_name"` //锚定节点名称
}

func (this *Punishment) TableName() string {
	return "punishments"
}

//添加惩罚记录
func (this *DataBaseAccessObject) CreatePunishment(instance *Punishment) (uint, error) {
	err := this.db.Create(instance).Error
	if err != nil {
		return 0, err
	}
	return instance.ID, nil
}

func (this *DataBaseAccessObject) GetPunishmentPage(start, pageSize int, anchorNodeId uint) ([]*PunishmentView, error) {
	result := make([]*PunishmentView, 0)
	sql := `select id,
    (select name from anchor_nodes where id=punishments.anchor_node_id) as anchor_node_name,
    value,
    coin,
    manage_type,
    date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    anchor_node_id from punishments`
	if anchorNodeId != 0 {
		sql += fmt.Sprintf(" where anchor_node_id=%d", anchorNodeId)
	}
	db := this.db.Raw(sql)
	err := db.Offset(start).
		Limit(pageSize).
		Find(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetPunishmentCount(anchorNodeId uint) (int, error) {
	var count int
	db := this.db.Table((&Punishment{}).TableName())
	if anchorNodeId != 0 {
		db = db.Where("anchor_node_id=?", anchorNodeId)
	}
	err := db.Count(&count).Error
	return count, err
}
