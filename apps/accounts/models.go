package accounts

import (
	"camp/core/utils"
	"camp/core/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
)

var _ web.Model = &UserModel{}

type UserModel struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func (m UserModel) TableName() string {
	return utils.NormalizeModelName(SubAppName, reflect.TypeOf(m).Name())
}

func (_ UserModel) IsGormModel() {}

type UserDB interface {
	// Methods for altering users
	Create(user *UserModel) error
	Update(user *UserModel) error
	Delete(id uint) error

	// Methods for querying for single user
	ByID(id uint) (*UserModel, error)
	ByEmail(email string) (*UserModel, error)
	ByRemember(token string) (*UserModel, error)
}
