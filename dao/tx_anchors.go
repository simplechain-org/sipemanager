package dao

import "fmt"

type TxAnchors struct {
	AnchorAddress   string `gorm:"primary_key; column:anchorAddress"` //锚定节点地址
	ContractAddress string `gorm:"column:contractAddress"`            //跨链合约地址
	SourceChainId   uint   `gorm:"source_chain_id"`
	TargetChainId   uint   `gorm:"target_chain_id"`
	SourceNetworkId uint64 `gorm:"source_network_id"`
	TargetNetworkId uint64 `gorm:"target_network_id"`
	AnchorId        uint   `gorm:"column:anchor_id"`
	Fee             uint64 `gorm:"column:fee"`
	Date            string `gorm:"primary_key; column:date"`
	Count           string `gorm:"column:count"`
	ChainId         uint   `gorm:"primary_key; column:chain_id" sql:"type:INT UNSIGNED NOT NULL"`
	TimeType        string `gorm:"column:timeType" `
}
type TxAnchorsNode struct {
	AnchorAddress   string `json:"anchorAddress"`   //锚定节点地址
	ContractAddress string `json:"contractAddress"` //跨链合约地址
	SourceChainId   uint   `json:"source_chain_id"`
	TargetChainId   uint   `json:"target_chain_id"`
	SourceNetworkId uint64 `json:"source_network_id"`
	TargetNetworkId uint64 `json:"target_network_id"`
	AnchorId        uint   `json:"anchor_id"`
	Fee             uint64 `json:"fee"`
	Date            string `json:"date"`
	Count           string `json:"count"`
	ChainId         uint   `json:"chain_id"`
	TimeType        string `json:"timeType" `
	Name            string `json:"name" `
}

func (this *TxAnchors) TableName() string {
	return "tx_anchors"
}

//func (this *DataBaseAccessObject) QueryAnchors(startTime string, endTime string, chainId uint, timeType string) ([]TxAnchors, error) {
//	result := make([]TxAnchors, 0)
//	err := this.db.Table((&TxAnchors{}).TableName()).Where("date BETWEEN ? and ? and chain_id = ? and timeType = ? ", startTime, endTime, chainId, timeType).Order("chain_id desc").Find(&result).Error
//
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}

func (this *DataBaseAccessObject) QueryAnchors(startTime string, endTime string, chainId int, timeType string) ([]TxAnchorsNode, error) {
	fmt.Printf("----45----%+v, %+v, %+v, %+v\n", startTime, endTime, chainId, timeType)
	txAnchors := make([]TxAnchorsNode, 0)
	var sql = `
SELECT anchorAddress, contractAddress, source_chain_id,target_chain_id, source_network_id, target_network_id, anchor_id, fee, date, count,chain_id,timeType, name  FROM
(SELECT * from tx_anchors WHERE date BETWEEN ? and ? and chain_id =? and timeType = ? ORDER BY chain_id asc ) t1
LEFT JOIN (SELECT id, name from anchor_nodes) t2 on t1.anchor_id = t2.id
`
	rows, err := this.db.Raw(sql, startTime, endTime, chainId, timeType).Rows()
	defer rows.Close()
	var result TxAnchorsNode
	for rows.Next() {
		rows.Scan(
			&result.AnchorAddress,
			&result.ContractAddress,
			&result.SourceChainId,
			&result.TargetChainId,
			&result.SourceNetworkId,
			&result.TargetNetworkId,
			&result.AnchorId,
			&result.Fee,
			&result.Date,
			&result.Count,
			&result.ChainId,
			&result.TimeType,
			&result.Name)
		txAnchors = append(txAnchors, result)
	}
	fmt.Printf("----4545----%+v\n", txAnchors)
	return txAnchors, err
}

func (this *DataBaseAccessObject) TxAnchorsReplace(data TxAnchors) error {
	var sql = "REPLACE INTO tx_anchors(anchorAddress, contractAddress, source_chain_id, target_chain_id, source_network_id, target_network_id,  anchor_id, fee, date, count, chain_id, timeType) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	return this.db.Exec(sql,
		data.AnchorAddress, data.ContractAddress, data.SourceChainId,
		data.TargetChainId, data.SourceNetworkId, data.TargetChainId, data.AnchorId, data.Fee,
		data.Date, data.Count, data.ChainId, data.TimeType).Error
}
