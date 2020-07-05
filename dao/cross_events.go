package dao

type CrossEvents struct {
	BlockNumber     uint64 `json:"block_number"`
	TxId            string `json:"tx_id"`
	Event           string `json:"event"`
	From            string `json:"from"`
	To              string `json:"to"`
	NetworkId       uint64 `gorm:"primary_key" json:"network_id"`
	RemoteNetworkId int64  `json:"remote_network_id"`
	Value           string `json:"value"`
	DestValue       string `json:"dest_value"` //目标金额
	TransactionHash string `gorm:"primary_key" json:"transaction_hash"`
	CrossAddress    string `json:"cross_address"`
	ChainId         uint   `json:"chain_id"`
}

func (this *CrossEvents) TableName() string {
	return "cross_events"
}

func (this *DataBaseAccessObject) MakerEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(block_number, tx_id, event, `from`, network_id, remote_network_id, value, dest_value, transaction_hash, cross_address, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.From, data.NetworkId, data.RemoteNetworkId,
		data.Value, data.DestValue, data.TransactionHash,
		data.CrossAddress, data.ChainId).Error
}

func (this *DataBaseAccessObject) TakerEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(block_number, tx_id, event, `from`, `to`, network_id, remote_network_id, value, dest_value, transaction_hash, cross_address, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.From, data.To, data.NetworkId,
		data.RemoteNetworkId, data.Value, data.DestValue,
		data.TransactionHash, data.CrossAddress, data.ChainId).Error
}

func (this *DataBaseAccessObject) MakerFinishEventUpsert(data CrossEvents) error {
	var sql = "REPLACE INTO cross_events(block_number, tx_id, event, `to`, network_id, transaction_hash, cross_address, chain_id) VALUES (?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.TxId, data.Event,
		data.To, data.NetworkId, data.TransactionHash,
		data.CrossAddress, data.ChainId).Error
}

func (this *DataBaseAccessObject) GetMaxCrossNumber(chainId uint) int64 {
	var blockNumber int64
	row := this.db.Raw("select IFNULL(max(block_number),0) blockNumber from cross_events where chain_id = ?", chainId).Row()
	row.Scan(&blockNumber)
	return blockNumber
}
