package accounts

import (
	"camp/core/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := utils.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)

	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

var _ UserService = &userService{}

type userService struct {
	UserDB
	pepper string
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, utils.GormErr.InvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}
