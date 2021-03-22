package accounts

import (
	"camp/core/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var _ web.Model = &User{}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func (u User) IsGormModel() {}

type UserDB interface {
	// Methods for altering users
	Create(user *User) error
	//Update(user *User) error
	//Delete(id uint) error

	// Methods for querying for single user
	//ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	//ByRemember(token string) (*User, error)
}
