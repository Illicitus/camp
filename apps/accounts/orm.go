package accounts

import "github.com/jinzhu/gorm"

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}
