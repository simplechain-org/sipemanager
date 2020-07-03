package dao

type Block struct {
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
	BlockHash    string `gorm:"column:blockHash"`
	ChainId      uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
}

func (this *Block) TableName() string {
	return "blocks"
}

func (this *DataBaseAccessObject) Create(block Block) (int64, error) {
	err := this.db.Create(block).Error
	if err != nil {
		return 0, err
	}
	return block.Number, nil
}

func (this *DataBaseAccessObject) GetNewBlockNumber(chainId uint) (int64, error) {
	var block Block
	err := this.db.Table((&Block{}).TableName()).Where("chain_id = ?", chainId).Order("number desc").Limit(1).Find(&block).Error
	if err != nil {
		return 0, err
	}
	return block.Number, nil
}

func (this *DataBaseAccessObject) BlockReplace(data Block) error {
	var sql = "REPLACE INTO blocks(parentHash, sha3Uncles, miner, difficulty, number, gasLimit, gasUsed, timestamp, nonce, transactions, blockHash, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.ParentHash, data.UncleHash, data.CoinBase,
		data.Difficulty, data.Number, data.GasLimit,
		data.GasUsed, data.Time, data.Nonce,
		data.Transactions, data.BlockHash, data.ChainId).Error
}
