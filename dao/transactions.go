package dao

//todo anchors From crossAddress makerFinish status
type Transaction struct {
	BlockHash        string `gorm:"column:blockHash"`
	BlockNumber      int64  `gorm:"column:blockNumber"`
	Hash             string `gorm:"primary_key" gorm:"column:hash"`
	From             string `gorm:"column:from"`
	GasUsed          uint64 `gorm:"column:gasUsed"`
	GasPrice         string `gorm:"column:gasPrice"`
	Input            []byte `gorm:"column:input;type:varbinary(50000);"`
	Nonce            uint64 `gorm:"column:nonce"`
	To               string `gorm:"column:to"`
	TransactionIndex int64  `gorm:"column:transactionIndex"`
	Value            string `gorm:"column:value"`
	Timestamp        uint64 `gorm:"column:timestamp"`
	Status           uint64 `gorm:"column:status"`
	ChainId          uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
}

func (this *Transaction) TableName() string {
	return "transactions"
}

func (this *DataBaseAccessObject) TxReplace(data Transaction) error {
	var sql = "REPLACE INTO transactions(blockHash, blockNumber, hash, `from`, gasUsed, gasPrice, input, nonce, `to`, transactionIndex, value, timestamp, status, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockHash, data.BlockNumber, data.Hash,
		data.From, data.GasUsed, data.GasPrice,
		data.Input, data.Nonce, data.To,
		data.TransactionIndex, data.Value, data.Timestamp,
		data.Status, data.ChainId).Error
}

func (this *DataBaseAccessObject) GetTxByHash(hash string) (*Transaction, error) {
	var tx Transaction
	err := this.db.Table((&Transaction{}).TableName()).Where("hash=?", hash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
