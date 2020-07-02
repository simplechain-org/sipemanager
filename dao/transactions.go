package dao

//todo anchors From crossAddress makerFinish status
type Transaction struct {
	BlockHash        string `gorm:"column:blockHash"`
	BlockNumber      int64  `gorm:"column:blockNumber"`
	Hash             string `gorm:"primary_key" gorm:"column:hash"`
	From             string `gorm:"column:from"`
	GasUsed          uint64 `gorm:"column:gasUsed"`
	GasPrice         uint64 `gorm:"column:gasPrice"`
	Input            string `gorm:"column:input"`
	Nonce            uint64 `gorm:"column:nonce"`
	To               string `gorm:"column:to"`
	TransactionIndex int64  `gorm:"column:transactionIndex"`
	Value            string `gorm:"column:value"`
	Timestamp        uint64 `gorm:"column:timestamp"`
	Status           string `gorm:"column:status"`
	ChainId          uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
}

func (this *Transaction) TableName() string {
	return "transactions"
}
