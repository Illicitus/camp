package accounts

import (
	"camp/core/utils"
	"github.com/jinzhu/gorm"
)

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User

	db := ug.db.Where("email = ?", email)
	return &user, utils.First(db, &user)
}
