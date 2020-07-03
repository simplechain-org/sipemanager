package dao

import "github.com/jinzhu/gorm"

func (this *DataBaseAccessObject) BeginTransaction() *gorm.DB {
	return this.db.Begin()
}
