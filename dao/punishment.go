package dao

import (
	"fmt"
	"time"
)

type Punishment struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	Value     string     `gorm:"value" json:"value"` //惩罚数量
	Coin      string     `gorm:"coin" json:"coin"`   //惩罚币种
	//suspend recovery token
	ManageType   string `gorm:"manage_type" json:"manage_type"`   //管理类型
	AnchorNodeId uint   `gorm:"anchor_node_id" json:"created_at"` //锚定节点编号
}

type PunishmentView struct {
	ID             uint   `gorm:"id" json:"id"`
	CreatedAt      string `gorm:"created_at" json:"created_at"`
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
	err := db.Offset(start).Limit(pageSize).Scan(&result).Error
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
func (this *DataBaseAccessObject) PunishmentRecordNotFound(anchorNodeId uint, manageType string) bool {
	var punishment Punishment
	return this.db.Table((&Punishment{}).TableName()).Where("anchor_node_id=?", anchorNodeId).
		Where("manage_type=?", manageType).First(&punishment).RecordNotFound()

}
func (this *DataBaseAccessObject) RemovePunishmentByManageType(anchorNodeId uint, manageType string) error {
	return this.db.Where("anchor_node_id = ?", anchorNodeId).Where("manage_type=?", manageType).Delete(&Punishment{}).Error
}
