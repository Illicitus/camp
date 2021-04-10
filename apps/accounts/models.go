package accounts

import (
	"camp/core/utils"
	"camp/core/web"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io"
	"net/url"
	"reflect"
)

var _ web.Model = &UserAvatarModel{}

type UserAvatarModel struct {
	UserID   uint      `gorm:"not null"`
	User     UserModel `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Filename string    `gorm:"not null"`
}

func (ua UserAvatarModel) TableName() string {
	return utils.NormalizeModelName(SubAppName, reflect.TypeOf(ua).Name())
}

func (ua UserAvatarModel) IsGormModel() {}

func (ua *UserAvatarModel) Path() string {
	tmp := url.URL{
		Path: "/" + ua.RelativePath(),
	}
	return tmp.String()
}

func (ua *UserAvatarModel) RelativePath() string {
	return fmt.Sprintf("%s/asserts/%v/%v", SubAppName, ua.UserID, ua.Filename)
}

type UserAvatarDB interface {

	// Methods for altering users avatars
	Create(userID uint, r io.ReadCloser, filename string) error
}

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

	ProfileByUserID(id uint) (*ProfileForm, error)
}
