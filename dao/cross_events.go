package dao

import "github.com/jinzhu/gorm"

type CrossEvents struct {
	gorm.Model
	BlockNumber     int    `json:"block_name"`
	TxId            string `json:"tx_id"`
	Event           string `json:"event"` //事件名
	From            string `json:"from"`
	To              string `json:"to"`
	ChainId         uint64 `json:"chain_id"`
	RemoteChainId   string `json:"remote_chain_id"`
	Value           string `json:"value"`
	DestValue       string `json:"dest_value"` //目标金额
	TransactionHash string `json:"transaction_hash"`
	crossAddress    string `json:"cross_address"`
}

func (this *CrossEvents) TableName() string {
	return "cross_events"
}

func (this *DataBaseAccessObject) CreateCrossEvents(cross *CrossEvents) (uint, error) {
	err := this.db.Create(cross).Error
	if err != nil {
		return 0, err
	}
	return cross.ID, nil
}
