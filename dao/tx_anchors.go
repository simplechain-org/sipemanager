package dao

type TxAnchors struct {
	From            string `gorm:"primary_key" gorm:"column:from"` //锚定节点地址
	SourceChainId   uint   `gorm:"source_chain_id"`
	TargetChainId   uint   `gorm:"target_chain_id"`
	SourceNetworkId uint64 `gorm:"source_network_id"`
	TargetNetworkId uint64 `gorm:"target_network_id"`
	AnchorId        uint   `gorm:"column:anchor_id"`
	Fee             uint64 `gorm:"column:fee"`
	Date            string `gorm:"primary_key" gorm:"column:date"`
	Count           string `gorm:"column:count"`
	ChainId         uint   `gorm:"primary_key" gorm:"column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
}

func (this *TxAnchors) TableName() string {
	return "tx_anchors"
}

func (this *DataBaseAccessObject) QueryAnchors(data Transaction) error {
	var sql = `
SELECT data.day,IFNULL(data.count, 0) count, IFNULL(data.fee,0) fee, day_list.day as date from
(select FROM_UNIXTIME(timestamp, '%Y-%m-%d %H:00:00') day, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED)) fee, count(1) count from transactions
where from = "0x17529b05513e5595ceff7f4fb1e06512c271a540" and to = "0xb11e0d62e216fc161fd7acfe7b4d36153ead89e0" and status = 1 and chain_id = 2
GROUP BY day) data
right join
(SELECT @date := DATE_ADD(@date, interval 1 hour) day from
(SELECT @date := DATE_ADD('2020-07-09', interval -1 hour) from transactions)
days LIMIT 24) day_list on day_list.day = data.day
`
	return this.db.Exec(sql).Error
}

func (this *DataBaseAccessObject) TxAnchorsReplace(data TxAnchors) error {
	var sql = "REPLACE INTO tx_anchors(`from`, source_chain_id, target_chain_id, source_network_id, target_network_id,  anchor_id, fee, date, count, chain_id) VALUES (?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.From, data.SourceChainId,
		data.TargetChainId, data.SourceNetworkId, data.TargetChainId, data.AnchorId, data.Fee,
		data.Date, data.Count, data.ChainId).Error
}

//`
//-- select FROM_UNIXTIME(`timestamp`,'%Y-%m-%d %H:00:00')as date,COUNT(*) count, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee
//-- FROM transactions
//-- GROUP BY FROM_UNIXTIME(`timestamp`,'%Y-%m-%d %H:00:00');
// SELECT IFNULL(data.count, 0) count, IFNULL(data.fee,0) fee, day_list.day as date from
//(select FROM_UNIXTIME(`timestamp`, '%Y-%m-%d %H:00:00') day, sum( CAST(gasUsed as SIGNED)* CAST(gasPrice as SIGNED) ) fee, count(1) count from transactions
//where `from` = "0x17529b05513e5595ceff7f4fb1e06512c271a540" and `to` = "0xb11e0d62e216fc161fd7acfe7b4d36153ead89e0" and `status` = 1
//GROUP BY day) data
//right join
//(SELECT @date := DATE_ADD(@date, interval 1 hour) day from
//(SELECT @date := DATE_ADD('2020-07-09', interval -1 hour) from transactions)
// days LIMIT 24) day_list on day_list.day = data.day`
