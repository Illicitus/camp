package accounts

import (
	"camp/core/utils"
	"github.com/jinzhu/gorm"
	"io"
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

	DB := ug.db.Where("id = ?", id)
	return &user, utils.First(DB, &user)
}

func (ug *userGorm) ByEmail(email string) (*UserModel, error) {
	var user UserModel

	DB := ug.db.Where("email = ?", email)
	return &user, utils.First(DB, &user)
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

func (ug *userGorm) ProfileByUserID(userID uint) (*ProfileForm, error) {
	var profile ProfileForm

	query := ug.db.Where("id = ?", userID).Table("accounts_users")
	query = query.Joins("LEFT JOIN accounts_user_avatars a on a.user_id = id")
	query = query.Select("accounts_users.*, a.fileName")

	err := utils.First(query, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

var _ UserAvatarDB = &userAvatarGorm{}

type userAvatarGorm struct {
	db *gorm.DB
}

func (uag *userAvatarGorm) Create(userID uint, r io.ReadCloser, filename string) error {
	userAvatar := UserAvatarModel{
		UserID:   userID,
		Filename: filename,
	}
	return uag.db.Create(userAvatar).Error
}
