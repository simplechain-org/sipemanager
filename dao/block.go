package dao

type Block struct {
	ParentHash   string `gorm:"column:parentHash"`
	UncleHash    string `gorm:"column:sha3Uncles"`
	CoinBase     string `gorm:"column:miner"`
	Difficulty   int64  `gorm:"column:difficulty"`
	Number       int64  `gorm:"primary_key; column:number" sql:"type:BIGINT UNSIGNED NOT NULL"`
	GasLimit     uint64 `gorm:"column:gasLimit"`
	GasUsed      uint64 `gorm:"column:gasUsed"`
	Time         uint64 `gorm:"column:timestamp"`
	Nonce        uint64 `gorm:"column:nonce"`
	Transactions int    `gorm:"column:transactions"`
	BlockHash    string `gorm:"column:blockHash"`
	ChainId      uint   `gorm:"primary_key; column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
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

func (this *DataBaseAccessObject) Delete(number int64, chainId uint) error {
	return this.db.Where("number < ? and chain_id =? ", number, chainId).Delete(&Block{}).Error
}

func (this *DataBaseAccessObject) GetNewBlockNumber(chainId uint) (int64, error) {
	var block Block
	err := this.db.Table((&Block{}).TableName()).Where("chain_id = ?", chainId).Order("number desc").Limit(1).Find(&block).Error
	if err != nil {
		return 0, err
	}
	return block.Number, nil
}

func (this *DataBaseAccessObject) GetMaxBlockNumber(chainId uint) int64 {
	var Number int64
	row := this.db.Raw("select IFNULL(max(number),0) number from blocks where chain_id = ?", chainId).Row()
	row.Scan(&Number)
	return Number
}

func (this *DataBaseAccessObject) BlockReplace(data Block) error {
	var sql = "REPLACE INTO blocks(parentHash, sha3Uncles, miner, difficulty, number, gasLimit, gasUsed, timestamp, nonce, transactions, blockHash, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.ParentHash, data.UncleHash, data.CoinBase,
		data.Difficulty, data.Number, data.GasLimit,
		data.GasUsed, data.Time, data.Nonce,
		data.Transactions, data.BlockHash, data.ChainId).Error
}

func (this *DataBaseAccessObject) UpdateBlock(data Block) error {

	number := this.db.Table((&Block{}).TableName()).Where("chain_id=?", data.ChainId).
		Update(Block{
			ParentHash:   data.ParentHash,
			UncleHash:    data.UncleHash,
			CoinBase:     data.CoinBase,
			Difficulty:   data.Difficulty,
			Number:       data.Number,
			GasLimit:     data.GasLimit,
			GasUsed:      data.GasUsed,
			Time:         data.Time,
			Nonce:        data.Nonce,
			Transactions: data.Transactions,
			BlockHash:    data.BlockHash,
			ChainId:      data.ChainId,
		}).RowsAffected
	if number == 0 {
		err := this.BlockReplace(data)
		if err != nil {
			return err
		}
	}
	return nil
}
