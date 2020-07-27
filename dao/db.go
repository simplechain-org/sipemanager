package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"net/url"
)

type DBConfig struct {
	Username string `default:"root" yaml:"username"`
	Password string `default:"" yaml:"password"`
	Address  string `yaml:"address"`
	Port     uint   `default:"3306" yaml:"port"`
	DbName   string `required:"true" yaml:"db_name"`
	Charset  string `default:"utf8" yaml:"charset"`
	MaxIdle  int    `default:"1000" yaml:"max_idle"`
	MaxOpen  int    `default:"2000" yaml:"max_open"`
	LogMode  bool   `yaml:"log_mode"`
	Loc      string `required:"true" yaml:"loc"`
}

func GetDBConnection(conf *DBConfig) (*gorm.DB, error) {
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	dsn := fmt.Sprintf(format, conf.Username, conf.Password, conf.Address, conf.Port, conf.DbName, conf.Charset, url.QueryEscape(conf.Loc))
	logrus.Infof("dsn=%s", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(conf.LogMode)
	db.DB().SetMaxIdleConns(conf.MaxIdle)
	db.DB().SetMaxOpenConns(conf.MaxOpen)
	return db, nil
}

type DataBaseAccessObject struct {
	db *gorm.DB
}

func NewDataBaseAccessObject(db *gorm.DB) *DataBaseAccessObject {
	return &DataBaseAccessObject{db: db}
}

func AutoMigrate(db *gorm.DB) {
	if db == nil {
		return
	}
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&UserNode{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Wallet{})
	db.AutoMigrate(&Chain{})
	db.AutoMigrate(&Contract{})
	db.AutoMigrate(&ContractInstance{})
	db.AutoMigrate(&MakerOrder{})
	db.AutoMigrate(&TakerOrder{})
	db.AutoMigrate(&ChainRegister{})
	db.AutoMigrate(&RetroActive{})
	db.AutoMigrate(&CrossEvents{})
	db.AutoMigrate(&Block{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&Uncle{})
	db.AutoMigrate(&AnchorNode{})
	db.AutoMigrate(&TxAnchors{})
	db.AutoMigrate(&CrossAnchors{})
	db.AutoMigrate(&ServiceChargeLog{})
	db.AutoMigrate(&SignRewardLog{})
	db.AutoMigrate(&Punishment{})
	db.AutoMigrate(&WorkCount{})
	db.AutoMigrate(&RewardConfig{})
	db.AutoMigrate(&PrepareReward{})
}
