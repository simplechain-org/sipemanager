package dao

type Uncle struct {
	ParentHash   string `gorm:"column:parentHash"`
	UncleHash    string `gorm:"column:sha3Uncles"`
	CoinBase     string `gorm:"column:miner"`
	Difficulty   int64  `gorm:"column:difficulty"`
	Number       int64  `gorm:"primary_key" gorm:"column:number" sql:"type:BIGINT UNSIGNED NOT NULL"`
	GasLimit     uint64 `gorm:"column:gasLimit"`
	GasUsed      uint64 `gorm:"column:gasUsed"`
	Time         uint64 `gorm:"column:timestamp"`
	Nonce        uint64 `gorm:"column:nonce"`
	Transactions int    `gorm:"column:transactions"`
	ChainId      uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
	BlockNumber  int64  `gorm:"column:blockNumber"`
}

func (this *Uncle) TableName() string {
	return "uncles"
}
