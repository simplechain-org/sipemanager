package dao

import (
	"github.com/sirupsen/logrus"
)

type Uncle struct {
	ParentHash  string `gorm:"column:parentHash"`
	UncleHash   string `gorm:"column:sha3Uncles"`
	CoinBase    string `gorm:"column:miner"`
	Difficulty  int64  `gorm:"column:difficulty"`
	Number      int64  `gorm:"primary_key" gorm:"column:number" sql:"type:BIGINT UNSIGNED NOT NULL"`
	GasLimit    uint64 `gorm:"column:gasLimit"`
	GasUsed     uint64 `gorm:"column:gasUsed"`
	Time        uint64 `gorm:"column:timestamp"`
	Nonce       uint64 `gorm:"column:nonce"`
	ChainId     uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
	BlockNumber int64  `gorm:"column:blockNumber"`
	BlockHash   string `gorm:"column:blockHash"`
}

func (this *Uncle) TableName() string {
	return "uncles"
}

func (this *DataBaseAccessObject) UncleReplace(data Uncle) error {
	var sql = "REPLACE INTO uncles(parentHash, sha3Uncles, miner, difficulty, number, gasLimit, gasUsed, timestamp, nonce, chain_id, blockNumber, blockHash) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.ParentHash, data.UncleHash, data.CoinBase,
		data.Difficulty, data.Number, data.GasLimit,
		data.GasUsed, data.Time, data.Nonce,
		data.ChainId, data.BlockNumber, data.BlockHash).Error
}

type MaxUncle struct {
	BlockNumber int64
	ChainId     uint
	ChainName   string
}

func (this *DataBaseAccessObject) QueryMaxUncle() ([]MaxUncle, error) {
	maxUncle := make([]MaxUncle, 0)
	rows, err := this.db.Raw("select IFNULL(max(blockNumber),0) blockNumber, chain_id from uncles GROUP BY chain_id").Rows()
	defer rows.Close()
	var result MaxUncle
	for rows.Next() {
		rows.Scan(&result.BlockNumber, &result.ChainId)
		chain, err := this.GetChain(result.ChainId)
		if err != nil {
			logrus.Error("QueryMaxUncle:", err.Error())
			continue
		}
		result.ChainName = chain.Name
		maxUncle = append(maxUncle, result)

	}
	return maxUncle, err
}
