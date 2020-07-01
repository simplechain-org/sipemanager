package dao

type CrossEvents struct {
	BlockNumber     uint64 `gorm:"column:blockNumber" json:"blockNumber"`
	TxId            string `gorm:"column:txId" json:"txId"`
	Event           string `gorm:"column:event" json:"event"` //事件名
	From            string `gorm:"column:from" json:"from"`
	To              string `gorm:"column:to" json:"to"`
	NetworkId       uint64 `gorm:"column:networkId" gorm:"primary_key" json:"networkId"`
	RemoteNetworkId int64  `gorm:"column:remoteNetworkId" json:"remoteNetworkId"`
	Value           string `gorm:"column:value" json:"value"`
	DestValue       string `gorm:"column:destValue" json:"destValue"` //目标金额
	TransactionHash string `gorm:"column:transactionHash" gorm:"primary_key" json:"transactionHash"`
	CrossAddress    string `gorm:"column:crossAddress" json:"crossAddress"`
	ChainId         uint   `json:"chain_id"`
}

func (this *CrossEvents) TableName() string {
	return "cross_events"
}

func (this *DataBaseAccessObject) MakerEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(blockNumber, txId, event, `from`, networkId, remoteNetworkId, value, destValue, transactionHash, crossAddress, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.From, data.NetworkId, data.RemoteNetworkId,
		data.Value, data.DestValue, data.TransactionHash,
		data.CrossAddress, data.ChainId).Error
}

func (this *DataBaseAccessObject) TakerEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(blockNumber, txId, event, `from`, `to`, networkId, remoteNetworkId, value, destValue, transactionHash, crossAddress, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.From, data.To, data.NetworkId,
		data.RemoteNetworkId, data.Value, data.DestValue,
		data.TransactionHash, data.CrossAddress, data.ChainId).Error
}

func (this *DataBaseAccessObject) MakerFinishEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(blockNumber, txId, event, `to`, networkId, transactionHash, crossAddress, chain_id) VALUES (?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.To, data.NetworkId, data.TransactionHash,
		data.CrossAddress, data.ChainId).Error
}
