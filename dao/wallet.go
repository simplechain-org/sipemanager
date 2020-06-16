package dao

import "github.com/jinzhu/gorm"

type Wallet struct {
	gorm.Model
	Name    string `gorm:"size:255" json:"name"`
	Content []byte `gorm:"type:text" json:"content"`
	UserId  uint
	Address string `gorm:"size:255" json:"address"`
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
