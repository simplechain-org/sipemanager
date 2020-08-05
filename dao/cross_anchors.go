package dao

import (
	"github.com/sirupsen/logrus"
	"strconv"
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
	Hash            string `grom:"column:hash"`
}

func (this *CrossAnchors) TableName() string {
	return "cross_anchors"
}

func (this *DataBaseAccessObject) CrossAnchorsReplace(data CrossAnchors) error {
	var sql = "REPLACE INTO cross_anchors(blockNumber, gasUsed, gasPrice, contractAddress, timestamp, status, chain_id, remote_chain_id, eventType, networkId, remoteNetworkId, anchorAddress, txId, hash) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.BlockNumber, data.GasUsed, data.GasPrice,
		data.ContractAddress, data.Timestamp, data.Status,
		data.ChainId, data.RemoteChainId, data.EventType,
		data.NetworkId, data.RemoteNetworkId, data.AnchorAddress,
		data.TxId, data.Hash).Error
}

func (this *DataBaseAccessObject) QueryTxByHours(txAnchors TxAnchors, EventType string) error {
	var sql = `
SELECT date_list date, IFNULL(count,0) count, IFNULL(fee,0) fee FROM
(
	(
	select FROM_UNIXTIME(timestamp,'%Y-%m-%d %H:00:00') as cross_date,COUNT(*) count, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee FROM cross_anchors
	WHERE anchorAddress = ? and contractAddress = ? and eventType = ? and networkId= ? and remoteNetworkId = ?
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
	rows, err := this.db.Raw(sql, txAnchors.AnchorAddress, txAnchors.ContractAddress, EventType, txAnchors.SourceNetworkId, txAnchors.TargetNetworkId).Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&txAnchors.Date,
			&txAnchors.Count,
			&txAnchors.Fee)
		txAnchors.TimeType = "hour"
		err := this.TxAnchorsReplace(txAnchors)
		if err != nil {
			logrus.Error("QueryTxByHours:", err.Error())
			return err
		}
	}
	return err
}

func (this *DataBaseAccessObject) QueryTxByDays(txAnchors TxAnchors, EventType string) error {
	var sql = `
SELECT date_list date, IFNULL(count,0) count, IFNULL(fee,0) fee FROM
(
	(
	select FROM_UNIXTIME(timestamp,'%Y-%m-%d') as cross_date,COUNT(*) count, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee FROM cross_anchors
	WHERE anchorAddress = ? and contractAddress = ? and eventType = ? and networkId= ? and remoteNetworkId = ?
	GROUP BY cross_date
	) t1 
	RIGHT JOIN
	(
		SELECT @cdate:= DATE_ADD(@cdate,INTERVAL - 1 day) AS date_list
		FROM (SELECT @cdate:=DATE_ADD(date_format(now(),'%Y-%m-%d'),INTERVAL + 1 day) FROM transactions) tmp1,(SELECT @mindt:=min(timestamp) from cross_anchors) s
		WHERE @cdate > FROM_UNIXTIME(@mindt,'%Y-%m-%d')
	) t2 
	ON t1.cross_date= t2.date_list
) ORDER BY t2.date_list desc
`
	rows, err := this.db.Raw(sql, txAnchors.AnchorAddress, txAnchors.ContractAddress, EventType, txAnchors.SourceNetworkId, txAnchors.TargetNetworkId).Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&txAnchors.Date,
			&txAnchors.Count,
			&txAnchors.Fee)
		txAnchors.TimeType = "day"
		err := this.TxAnchorsReplace(txAnchors)
		if err != nil {
			logrus.Error("QueryTxByDays:", err.Error())
			return err
		}
	}
	return err
}

func (this *DataBaseAccessObject) QueryTxByWeeks(txAnchors TxAnchors, EventType string) error {
	var sql = `
SELECT date_list date, IFNULL(count,0) count, IFNULL(fee,0) fee FROM
	(
	SELECT
			count(*) count,
			sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee,
			yearweek(FROM_UNIXTIME(timestamp,'%Y-%m-%d')) cross_date
	FROM
			cross_anchors
	WHERE anchorAddress = ? and contractAddress = ? and eventType = ? and networkId= ? and remoteNetworkId = ?
	GROUP BY
		 cross_date
	) t1 
RIGHT JOIN
	(
		SELECT @cdate:= (@cdate - 1) AS date_list
		FROM (SELECT @cdate:=( yearweek(now()) +1) FROM transactions) tmp1,(SELECT @mindt:=min(timestamp) from cross_anchors) s
		WHERE @cdate > YEARWEEK(FROM_UNIXTIME(@mindt,'%Y-%m-%d'))
	) t2 
ON t1.cross_date= t2.date_list
`
	rows, err := this.db.Raw(sql, txAnchors.AnchorAddress, txAnchors.ContractAddress, EventType, txAnchors.SourceNetworkId, txAnchors.TargetNetworkId).Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(
			&txAnchors.Date,
			&txAnchors.Count,
			&txAnchors.Fee)
		txAnchors.TimeType = "week"
		err := this.TxAnchorsReplace(txAnchors)
		if err != nil {
			logrus.Error("QueryTxByWeeks:", err.Error())
			return err
		}
	}
	return err
}

func (this *DataBaseAccessObject) QueryFinishList(offset, limit uint32, startTimeParam, endTimeParam, anchorIdParam string) ([]CrossAnchors, uint32, error) {
	var count uint32
	result := make([]CrossAnchors, 0)
	switch true {
	case startTimeParam != "" && endTimeParam != "" && anchorIdParam != "":
		startTime, stringErr := strconv.Atoi(startTimeParam)
		endTime, stringErr := strconv.Atoi(endTimeParam)
		anchorId, stringErr := strconv.Atoi(anchorIdParam)
		if stringErr != nil {
			return nil, 0, stringErr
		}
		anchor, anchorErr := this.GetAnchorNode(uint(anchorId))
		if anchorErr != nil {
			return nil, 0, anchorErr
		}
		if err := this.db.Model(&CrossAnchors{}).Where("timestamp between ? and ? ", startTime, endTime).Where("anchorAddress = ?", anchor.Address).Count(&count).Error; err != nil {
			return nil, 0, err
		}

		err := this.db.Table((&CrossAnchors{}).TableName()).Where("timestamp between ? and ? ", startTime, endTime).Where("anchorAddress = ?", anchor.Address).Order("timestamp desc").Offset(offset).Limit(limit).Find(&result).Error
		return result, count, err

	case startTimeParam != "" && endTimeParam != "":
		startTime, stringErr := strconv.Atoi(startTimeParam)
		endTime, stringErr := strconv.Atoi(endTimeParam)
		if stringErr != nil {
			return nil, 0, stringErr
		}
		if err := this.db.Model(&CrossAnchors{}).Where("timestamp between ? and ? ", startTime, endTime).Count(&count).Error; err != nil {
			return nil, 0, err
		}

		err := this.db.Table((&CrossAnchors{}).TableName()).Where("timestamp between ? and ? ", startTime, endTime).Order("timestamp desc").Offset(offset).Limit(limit).Find(&result).Error
		return result, count, err

	case startTimeParam == "" && endTimeParam == "" && anchorIdParam != "":
		anchorId, stringErr := strconv.Atoi(anchorIdParam)
		if stringErr != nil {
			return nil, 0, stringErr
		}
		anchor, anchorErr := this.GetAnchorNode(uint(anchorId))
		if anchorErr != nil {
			return nil, 0, anchorErr
		}
		if err := this.db.Model(&CrossAnchors{}).Where("anchorAddress = ?", anchor.Address).Count(&count).Error; err != nil {
			return nil, 0, err
		}
		err := this.db.Table((&CrossAnchors{}).TableName()).Where("anchorAddress = ?", anchor.Address).Order("timestamp desc").Offset(offset).Limit(limit).Find(&result).Error
		return result, count, err
	default:
		if err := this.db.Model(&CrossAnchors{}).Count(&count).Error; err != nil {
			return nil, 0, err
		}

		err := this.db.Table((&CrossAnchors{}).TableName()).Order("timestamp desc").Offset(offset).Limit(limit).Find(&result).Error
		return result, count, err

	}
}
