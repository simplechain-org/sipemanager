package dao

import (
	"github.com/sirupsen/logrus"
)

type CrossAnchors struct {
	BlockNumber     int64  `gorm:"column:blockNumber"`
	GasUsed         string `gorm:"column:gasUsed"`
	GasPrice        string `gorm:"column:gasPrice"`
	ContractAddress string `gorm:"column:contractAddress"`
	Timestamp       uint64 `gorm:"column:timestamp"`
	Status          uint64 `gorm:"column:status"`
	ChainId         uint   `gorm:"primary_key; column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
	RemoteChainId   uint   `gorm:"column:remote_chain_id"`
	EventType       string `gorm:"column:eventType"`
	NetworkId       uint64 `gorm:"column:networkId"`
	RemoteNetworkId uint64 `gorm:"column:remoteNetworkId"`
	AnchorAddress   string `gorm:"primary_key; column:anchorAddress"`
	TxId            string `gorm:"primary_key; column:txId" `
}

func (this *CrossAnchors) TableName() string {
	return "cross_anchors"
}

func (this *DataBaseAccessObject) CrossAnchorsReplace(data CrossAnchors) error {
	var sql = "REPLACE INTO cross_anchors(blockNumber, gasUsed, gasPrice, contractAddress, timestamp, status, chain_id, remote_chain_id, eventType, networkId, remoteNetworkId, anchorAddress, txId) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.GasUsed, data.GasPrice,
		data.ContractAddress, data.Timestamp, data.Status,
		data.ChainId, data.RemoteChainId, data.EventType,
		data.NetworkId, data.RemoteNetworkId, data.AnchorAddress,
		data.TxId).Error
}

//toDo 增加跨链合约地址进行查询
func (this *DataBaseAccessObject) QueryTxByHours(txAnchors TxAnchors, EventType string) error {
	var sql = `
SELECT date_list date, IFNULL(count,0) count, IFNULL(fee,0) fee FROM
(
	(
	select FROM_UNIXTIME(timestamp,'%Y-%m-%d %H:00:00') as cross_date,COUNT(*) count, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee FROM cross_anchors
	WHERE anchorAddress = ? and eventType = ? and networkId= ? and remoteNetworkId = ?
	GROUP BY cross_date
	) t1 
	RIGHT JOIN
	(
		SELECT @cdate:= DATE_ADD(@cdate,INTERVAL - 1 hour) AS date_list
		FROM (SELECT @cdate:=DATE_ADD(date_format(now(),'%Y-%m-%d %H:00:00'),INTERVAL + 1 hour) FROM transactions) tmp1,(SELECT @mindt:=min(timestamp) from cross_anchors) s
		WHERE @cdate > FROM_UNIXTIME(@mindt,'%Y-%m-%d %H:00:00')
	) t2 
	ON t1.cross_date= t2.date_list
) ORDER BY t2.date_list desc
`
	rows, err := this.db.Raw(sql, txAnchors.AnchorAddress, EventType, txAnchors.SourceNetworkId, txAnchors.TargetNetworkId).Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&txAnchors.Date,
			&txAnchors.Count,
			&txAnchors.Fee)
		err := this.TxAnchorsReplace(txAnchors)
		if err != nil {
			logrus.Error("QueryTxByHours", err.Error())
		}
	}
	return err
}
