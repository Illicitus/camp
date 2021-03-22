package utils

import (
	"github.com/jinzhu/gorm"
)

func First(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return GormErr.NotFound
	}
	return err

}
