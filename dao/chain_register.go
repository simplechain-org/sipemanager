package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type ChainRegister struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	SourceChainId   uint       `json:"source_chain_id"`
	TargetChainId   uint       `json:"target_chain_id"`
	Confirm         uint       `json:"confirm"`
	AnchorAddresses string     `json:"anchor_addresses"`
	Status          int        `json:"status"`
	StatusText      string     `json:"status_text"`
	TxHash          string     `json:"tx_hash"`
	Address         string     `json:"address"` // 合约地址
}

func (this *ChainRegister) TableName() string {
	return "chain_registers"
}

func (this *DataBaseAccessObject) CreateChainRegister(chainRegister *ChainRegister) (uint, error) {
	err := this.db.Create(chainRegister).Error
	if err != nil {
		return 0, err
	}
	return chainRegister.ID, nil
}
func (this *DataBaseAccessObject) CreateChainRegisterByTx(db *gorm.DB, chainRegister *ChainRegister) (uint, error) {
	err := db.Create(chainRegister).Error
	if err != nil {
		return 0, err
	}
	return chainRegister.ID, nil
}
func (this *DataBaseAccessObject) UpdateChainRegisterStatus(id uint, status int) error {
	text := "失败"
	if status == 1 {
		text = "成功"
	}
	err := this.db.Model(&ChainRegister{}).Where("id=?", id).
		Updates(ChainRegister{Status: status, StatusText: text}).Error
	if err != nil {
		return err
	}
	return nil
}

//根据当前所在的链和目标的网络id获取对应的链id
func (this *DataBaseAccessObject) GetTargetChainId(sourceChainId uint, targetNetworkId uint64) (uint, error) {
	type Result struct {
		TargetChainId uint
	}
	var result Result
	err := this.db.Model(&ChainRegister{}).
		Where("chain_registers.source_chain_id=?", sourceChainId).
		Where("chains.network_id=?", targetNetworkId).
		Joins("inner join chains on chain_registers.target_chain_id=chains.id").
		Select("chain_registers.target_chain_id").
		Order("chain_registers.id desc").
		Limit(1).Scan(&result).Error
	return result.TargetChainId, err
}

func (this *DataBaseAccessObject) ListChainRegister() ([]ChainRegister, error) {
	result := make([]ChainRegister, 0)
	err := this.db.Table((&ChainRegister{}).TableName()).Order("id desc").Find(&result).Error
	return result, err
}

func (this *DataBaseAccessObject) ListChainRegisterByStatus(status int) ([]ChainRegister, error) {
	result := make([]ChainRegister, 0)
	err := this.db.Table((&ChainRegister{}).TableName()).
		Where("status=?", status).
		Order("id desc").
		Find(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetTargetChainIdBySourceChainId(sourceChainId uint) (uint, error) {
	type Result struct {
		TargetChainId uint
	}
	var result Result
	err := this.db.Model(&ChainRegister{}).
		Where("source_chain_id=?", sourceChainId).
		Select("target_chain_id").
		Order("id desc").
		Limit(1).Scan(&result).Error
	return result.TargetChainId, err
}
func (this *DataBaseAccessObject) GetChainRegisterByChaiId(sourceChainId uint, targetChainId uint) (*ChainRegister, error) {
	var result ChainRegister
	err := this.db.Table((&ChainRegister{}).TableName()).
		Where("source_chain_id=? and target_chain_id=?", sourceChainId, targetChainId).
		Or("source_chain_id=? and target_chain_id=?", targetChainId, sourceChainId).First(&result).Error
	return &result, err
}

type TokenListInterface struct {
	ChainID            uint
	RemoteChainID      uint
	AnchorAddresses    string
	SourceCrossAddress string
	RemoteCrossAddress string
	NetworkId          uint64
	RemoteNetworkId    uint64
	Name               string
	Count              int
}

func (this *DataBaseAccessObject) GetTxTokenList() (map[string]TokenListInterface, error) {
	result := make([]ChainRegister, 0)
	opposite := make([]ChainRegister, 0)
	TokenList := make(map[string]TokenListInterface, 0)
	err := this.db.Table((&ChainRegister{}).TableName()).Where("status = 1").Order("id desc").Find(&result).Error
	for _, item := range result {
		err := this.db.Table((&ChainRegister{}).TableName()).Where("source_chain_id=? and target_chain_id =? and status = 1", item.TargetChainId, item.SourceChainId).Find(&opposite).Error
		if err != nil {
			return nil, err
		}
		if len(opposite) > 1 {
			sourceId := strconv.Itoa(int(item.SourceChainId))
			targetId := strconv.Itoa(int(item.TargetChainId))
			source, err := this.GetChain(item.SourceChainId)
			target, err := this.GetChain(item.TargetChainId)
			if err != nil {
				return nil, err
			}
			tokenList := TokenListInterface{
				ChainID:            item.SourceChainId,
				RemoteChainID:      item.TargetChainId,
				AnchorAddresses:    item.AnchorAddresses,
				SourceCrossAddress: item.Address,
				RemoteCrossAddress: opposite[0].Address,
				NetworkId:          source.NetworkId,
				RemoteNetworkId:    target.NetworkId,
				Name:               source.Name + " <=> " + target.Name,
			}
			if len(TokenList) != 0 {
				for key, _ := range TokenList {
					if !(strings.Contains(key, sourceId) && strings.Contains(key, targetId)) {
						TokenList[sourceId+","+targetId] = tokenList
					}
				}
			} else {
				TokenList[sourceId+","+targetId] = tokenList
			}
		}
	}

	return TokenList, err
}

type ChainRegisterView struct {
	ID              uint   `json:"id" gorm:"id"`
	CreatedAt       string `json:"created_at" gorm:"created_at"` //创建时间
	SourceChainId   uint   `json:"source_chain_id" gorm:"source_chain_id"`
	TargetChainId   uint   `json:"target_chain_id" gorm:"target_chain_id"`
	SourceChainName string `json:"source_chain_name" gorm:"source_chain_name"`
	TargetChainName string `json:"target_chain_name" gorm:"target_chain_name"`
	Confirm         uint   `json:"confirm" gorm:"confirm"`
	AnchorAddresses string `json:"anchor_addresses" gorm:"anchor_addresses"`
	TxHash          string `json:"tx_hash" gorm:"tx_hash"`
	Status          int    `json:"status" gorm:"status"`
	StatusText      string `json:"status_text" gorm:"status_text"`
}

func (this *DataBaseAccessObject) GetChainRegisterPage(start, pageSize int) ([]*ChainRegisterView, error) {
	sql := `select 
     id,
    (select name from chains where chains.id=chain_registers.source_chain_id) as source_chain_name,
    (select name from chains where chains.id=chain_registers.target_chain_id) as target_chain_name,
    date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    source_chain_id,
    target_chain_id,
    anchor_addresses,
    confirm,
    status,
    status_text,
    tx_hash from chain_registers where deleted_at is null`
	result := make([]*ChainRegisterView, 0)
	db := this.db.Raw(sql)
	err := db.Offset(start).
		Limit(pageSize).
		Scan(&result).Error
	return result, err
}
func (this *DataBaseAccessObject) GetChainRegisterCount() (int, error) {
	sql := `select count(*) as total from chain_registers where deleted_at is null`
	var total Total
	err := this.db.Raw(sql).Scan(&total).Error
	return total.Total, err
}
func (this *DataBaseAccessObject) GetChainRegister(id uint) (*ChainRegisterView, error) {
	sql := `select 
     id,
    (select name from chains where chains.id=chain_registers.source_chain_id) as source_chain_name,
    (select name from chains where chains.id=chain_registers.target_chain_id) as target_chain_name,
    date_format(created_at,'%Y-%m-%d %H:%i:%S') as created_at,
    source_chain_id,
    target_chain_id,
    anchor_addresses,
    confirm,
	status,
	status_text,
    tx_hash from chain_registers where deleted_at is null and chain_registers.id=?`
	var result ChainRegisterView
	db := this.db.Raw(sql, id).Limit(1)
	err := db.Scan(&result).Error
	return &result, err
}
func (this *DataBaseAccessObject) GetChainRegisterById(id uint) (*ChainRegister, error) {
	var result ChainRegister
	err := this.db.Model(&ChainRegister{}).Where("id=?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (this *DataBaseAccessObject) UpdateChainRegisterAnchorAddresses(sourceChainId, targetChainId uint, address string, symbol string, anchorNodeId uint) error {
	var result ChainRegister
	err := this.db.Model(&ChainRegister{}).
		Where("source_chain_id=?", sourceChainId).
		Where("target_chain_id=?", targetChainId).
		Where("address=?", address).First(&result).Error
	if err != nil {
		return err
	}
	anchorNodeIdStr := fmt.Sprintf("%d", anchorNodeId)
	if symbol == "+" {
		if result.AnchorAddresses != "" {
			ids := strings.Split(result.AnchorAddresses, ",")
			for _, id := range ids {
				if id == anchorNodeIdStr {
					//已经存在
					return errors.New("锚定节点已经存在")
				}
			}
			ids = append(ids, anchorNodeIdStr)
			anchorNodeAddress := strings.Join(ids, ",")
			return this.db.Table((&ChainRegister{}).TableName()).Where("id = ?", result.ID).Update("anchor_addresses", anchorNodeAddress).Error
		}
	} else if symbol == "-" {
		if result.AnchorAddresses != "" {
			updateAddress := make([]string, 0)
			ids := strings.Split(result.AnchorAddresses, ",")
			for _, id := range ids {
				if id != anchorNodeIdStr {
					//把anchorNodeIdStr从anchor_addresses中移除
					updateAddress = append(updateAddress, id)
				}
			}
			anchorNodeAddress := strings.Join(updateAddress, ",")
			return this.db.Table((&ChainRegister{}).TableName()).Where("id = ?", result.ID).Update("anchor_addresses", anchorNodeAddress).Error
		}
	}
	return nil
}
func (this *DataBaseAccessObject) ChainRegisterRecordNotFound(sourceChainId, targetChainId uint, address string, status int) bool {
	var result ChainRegister
	return this.db.Model(&ChainRegister{}).
		Where("source_chain_id=?", sourceChainId).
		Where("target_chain_id=?", targetChainId).
		Where("address=?", address).
		Where("status=?", status).First(&result).RecordNotFound()

}
func (this *DataBaseAccessObject) GetChainRegisterWithAddress(sourceChainId, targetChainId uint, address string, status int) (*ChainRegister, error) {
	var result ChainRegister
	err := this.db.Model(&ChainRegister{}).
		Where("source_chain_id=?", sourceChainId).
		Where("target_chain_id=?", targetChainId).
		Where("address=?", address).
		Where("status=?", status).First(&result).Error
	if err != nil {
		err = this.db.Model(&ChainRegister{}).
			Where("source_chain_id=?", targetChainId).
			Where("target_chain_id=?", sourceChainId).
			Where("address=?", address).
			Where("status=?", status).First(&result).Error
	}
	return &result, err

}

func (this *DataBaseAccessObject) GetRegisterLatestBySourChainId(sourceChainId int) (*ChainRegister, error) {
	var result ChainRegister
	err := this.db.Table((&ChainRegister{}).TableName()).Where("source_chain_id = ? and status = 1", sourceChainId).Order("id desc").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
