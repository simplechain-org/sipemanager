package dao

import (
	"fmt"
	"math/big"
	"time"
)

//报销手续费日志
type ServiceChargeLog struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	AnchorNodeId    uint       `gorm:"anchor_node_id"`   //锚定节点编号
	TransactionHash string     `gorm:"transaction_hash"` //交易哈希
	Fee             string     `gorm:"fee"`              //报销手续费
	Coin            string     `gorm:"coin"`             //报销的币种
	Sender          string     `gorm:"sender"`           //出账账户地址
	Status          uint       `gorm:"status"`           //状态
}

type ServiceChargeLogView struct {
	ID              uint   `gorm:"id" json:"id"`
	CreatedAt       string `gorm:"created_at" json:"created_at"`
	AnchorNodeId    uint   `gorm:"anchor_node_id" json:"anchor_node_id"`     //锚定节点编号
	TransactionHash string `gorm:"transaction_hash" json:"transaction_hash"` //交易哈希
	Fee             string `gorm:"fee" json:"fee"`                           //报销手续费
	Coin            string `gorm:"coin" json:"coin"`                         //报销的币种
	Sender          string `gorm:"sender" json:"sender"`                     //出账账户地址
	Status          uint   `gorm:"status" json:"status"`                     //状态
	AnchorNodeName  string `gorm:"anchor_node_name" json:"anchor_node_name"` //锚定节点名称
}

func (this *ServiceChargeLog) TableName() string {
	return "service_charge_logs"
}

func (this *DataBaseAccessObject) CreateServiceChargeLog(obj *ServiceChargeLog) (uint, error) {
	err := this.db.Create(obj).Error
	if err != nil {
		return 0, err
	}
	return obj.ID, nil
}
func (this *DataBaseAccessObject) UpdateServiceChargeLogSourceStatus(id uint, status uint) error {
	return this.db.Table((&ServiceChargeLog{}).TableName()).
		Where("id=?", id).
		Update("status", status).Error
}

func (this *DataBaseAccessObject) GetServiceChargeLogCount(anchorNodeId uint) (int, error) {
	var count int

	db := this.db.Table((&ServiceChargeLog{}).TableName()).Where("status=?", 1)

	if anchorNodeId != 0 {
		db = db.Where("anchor_node_id=?", anchorNodeId)
	}
	err := db.Count(&count).Error //表示已经成功上链的数据

	return count, err
}

func (this *DataBaseAccessObject) GetServiceChargeLogPage(start, pageSize int, anchorNodeId uint) ([]*ServiceChargeLogView, error) {
	result := make([]*ServiceChargeLogView, 0)
	sql := `select id,
    (select name from anchor_nodes where id=service_charge_logs.anchor_node_id) as anchor_node_name,
     date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    anchor_node_id,
    transaction_hash,
    fee,
    coin,
    sender,
    status from service_charge_logs where status=1 and deleted_at is null`
	if anchorNodeId != 0 {
		sql += fmt.Sprintf(" and anchor_node_id=%d", anchorNodeId)
	}
	err := this.db.Raw(sql).Offset(start).Limit(pageSize).Scan(&result).Error
	return result, err
}

func (this *DataBaseAccessObject) GetServiceChargeLog(id uint) (*ServiceChargeLog, error) {
	var obj ServiceChargeLog
	err := this.db.Table((&ServiceChargeLog{}).TableName()).
		Where("id=?", id).
		First(&obj).Error
	return &obj, err
}

func (this *DataBaseAccessObject) GetServiceChargeSumFee(anchorNodeId uint, coin string) (*big.Int, error) {
	sum := big.NewInt(0)
	result := make([]ServiceChargeLog, 0)
	err := this.db.Table((&ServiceChargeLog{}).TableName()).
		Where("anchor_node_id=?", anchorNodeId).
		Where("coin=?", coin).
		Where("status=?", 1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	for _, o := range result {
		fee, success := big.NewInt(0).SetString(o.Fee, 10)
		if !success {
			continue
		}
		sum = sum.Add(sum, fee)
	}
	return sum, nil
}
