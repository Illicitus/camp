package accounts

import (
	"camp/core/utils"
	"github.com/jinzhu/gorm"
)

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) Create(user *UserModel) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) ByID(id uint) (*UserModel, error) {
	var user UserModel

	db := ug.db.Where("id = ?", id)
	return &user, utils.First(db, &user)
}

func (ug *userGorm) ByEmail(email string) (*UserModel, error) {
	var user UserModel

	db := ug.db.Where("email = ?", email)
	return &user, utils.First(db, &user)
}

func (ug *userGorm) ByRemember(rememberHash string) (*UserModel, error) {
	var user UserModel

	err := utils.First(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) Update(user *UserModel) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := UserModel{Model: gorm.Model{ID: id}}

	return ug.db.Delete(&user).Error
}
