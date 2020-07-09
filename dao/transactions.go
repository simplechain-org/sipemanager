package dao

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Transaction struct {
	BlockHash        string `gorm:"column:blockHash"`
	BlockNumber      int64  `gorm:"column:blockNumber"`
	Hash             string `gorm:"primary_key" gorm:"column:hash"`
	From             string `gorm:"column:from"`
	GasUsed          string `gorm:"column:gasUsed"`
	GasPrice         string `gorm:"column:gasPrice"`
	Input            []byte `gorm:"column:input;type:varbinary(50000);"`
	Nonce            uint64 `gorm:"column:nonce"`
	To               string `gorm:"column:to"`
	TransactionIndex int64  `gorm:"column:transactionIndex"`
	Value            string `gorm:"column:value"`
	Timestamp        uint64 `gorm:"column:timestamp"`
	Status           uint64 `gorm:"column:status"`
	ChainId          uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
	EventType        string `gorm:"column:eventType"`
}

type Recept struct {
	TxId   common.Hash
	TxHash common.Hash
	From   common.Address
	To     common.Address
}

type MakerFinish struct {
	Rtx           Recept
	RemoteChainId *big.Int
}

func (this *Transaction) TableName() string {
	return "transactions"
}

func (this *DataBaseAccessObject) TxReplace(data Transaction) error {
	var sql = "REPLACE INTO transactions(blockHash, blockNumber, hash, `from`, gasUsed, gasPrice, input, nonce, `to`, transactionIndex, value, timestamp, status, chain_id, eventType) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockHash, data.BlockNumber, data.Hash,
		data.From, data.GasUsed, data.GasPrice,
		data.Input, data.Nonce, data.To,
		data.TransactionIndex, data.Value, data.Timestamp,
		data.Status, data.ChainId, data.EventType).Error
}

func (this *DataBaseAccessObject) GetTxByHash(hash string) (*Transaction, error) {
	var tx Transaction
	err := this.db.Table((&Transaction{}).TableName()).Where("hash=?", hash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (this *DataBaseAccessObject) GetTxByAnchors(chainId uint, from string, to string) ([]Transaction, error) {
	result := make([]Transaction, 0)
	err := this.db.Table((&Transaction{}).TableName()).Where("chain_id=? and from =? and to =?", chainId, from, to).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
