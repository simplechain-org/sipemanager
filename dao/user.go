package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255"`
	Password string `gorm:"size:255"`
}

func (u *User) TableName() string {
	return "users"
}

var userTableName = (&User{}).TableName()

func cryptoPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	return encodePW, nil
}
func CheckPassword(encodePW string, password string) bool {
	// 正确密码验证
	return bcrypt.CompareHashAndPassword([]byte(encodePW), []byte(password)) == nil
}

func (this *DataBaseAccessObject) CreateUser(user *User) (uint, error) {
	var count int
	err := this.db.Table(userTableName).Where("username=?", user.Username).Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("Record already exists")
	}
	encodePwd, err := cryptoPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = encodePwd
	err = this.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (this *DataBaseAccessObject) GetUserByUsername(username string) (*User, error) {
	var user User
	err := this.db.Table(userTableName).Where("username=?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (this *DataBaseAccessObject) UserIsValid(username, password string) bool {
	user, err := this.GetUserByUsername(username)
	if err != nil {
		return false
	}
	return CheckPassword(user.Password, password)
}


