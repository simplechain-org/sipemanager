package dao

import "github.com/jinzhu/gorm"

type Punishment struct {
	gorm.Model
	Value string //惩罚数量
	Coin  string //惩罚币种
	//suspend recovery token
	ManageType     string //管理类型
	AnchorNodeId   uint   `gorm:"anchor_node_id"`   //锚定节点编号
	AnchorNodeName string `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
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

func (this *DataBaseAccessObject) GetPunishmentPage(start, pageSize int, anchorNodeId uint) ([]*Punishment, error) {
	result := make([]*Punishment, 0)
	db := this.db.Table((&Punishment{}).TableName())
	if anchorNodeId != 0 {
		db = db.Where("anchor_node_id=?", anchorNodeId)
	}
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
