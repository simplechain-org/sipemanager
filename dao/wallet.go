package dao

import (
	"fmt"
	"time"
)

type Wallet struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"deleted_at" json:"deleted_at"`
	Name      string     `gorm:"size:255" json:"name"`
	Content   string     `gorm:"type:text" json:"content"`
	UserId    uint
	Address   string `gorm:"size:255" json:"address"`
}

func (this *Wallet) TableName() string {
	return "wallets"
}

func (this *DataBaseAccessObject) CreateWallet(wallet *Wallet) (uint, error) {
	err := this.db.Create(wallet).Error
	if err != nil {
		return 0, err
	}
	return wallet.ID, nil
}
func (this *DataBaseAccessObject) ListWalletByUserId(userId uint) ([]Wallet, error) {
	wallets := make([]Wallet, 0)
	err := this.db.Table((&Wallet{}).TableName()).Select("id,name,user_id,address").Where("user_id=?", userId).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}
func (this *DataBaseAccessObject) GetWallet(id uint) (*Wallet, error) {
	var wallet Wallet
	err := this.db.Table((&Wallet{}).TableName()).Where("id=?", id).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
func (this *DataBaseAccessObject) UpdateWallet(id uint, content []byte) error {
	return this.db.Table((&Wallet{}).TableName()).
		Where("id=?", id).
		Update("content", content).Error
}

func (this *DataBaseAccessObject) WalletExists(address string) bool {
	var count int

	db := this.db.Table((&Wallet{}).TableName()).Where("address=? and deleted_at is null", address)

	err := db.Count(&count).Error

	if err != nil {
		fmt.Println(err)
	}
	return count != 0
}

func (this *DataBaseAccessObject) RemoveWallet(id uint) error {
	return this.db.Where("id = ?", id).Delete(&Wallet{}).Error
}

type WalletView struct {
	ID        uint   `gorm:"id" json:"id"`
	CreatedAt string `gorm:"created_at" json:"created_at"`
	Name      string `gorm:"size:255" json:"name"`
	Content   string `gorm:"type:text" json:"content"`
	Address   string `gorm:"size:255" json:"address"`
}

func (this *DataBaseAccessObject) ListWalletViewByUserId(userId uint) ([]WalletView, error) {
	sql := `select
			id,
			name,
			content,
			address,
			date_format(wallets.created_at,'%Y-%m-%d %H:%i:%S') as created_at
            from wallets where wallets.user_id=? and wallets.deleted_at is null`
	wallets := make([]WalletView, 0)
	err := this.db.Raw(sql, userId).Scan(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}
func (this *DataBaseAccessObject) GetWalletViewCount(userId uint) (int, error) {
	sql := `select count(*) as total from wallets where wallets.user_id=? and wallets.deleted_at is null`
	var total Total
	err := this.db.Raw(sql, userId).Scan(&total).Error
	return total.Total, err
}
func (this *DataBaseAccessObject) GetWalletViewPage(userId uint, start, pageSize int) ([]*WalletView, error) {
	sql := `select
			id,
			name,
			content,
			address,
			date_format(wallets.created_at,'%Y-%m-%d %H:%i:%S') as created_at
            from wallets where wallets.user_id=? and wallets.deleted_at is null`
	wallets := make([]*WalletView, 0)
	err := this.db.Raw(sql, userId).Offset(start).
		Limit(pageSize).Scan(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, err
}
