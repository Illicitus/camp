package accounts

import (
	"camp/core/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	Authenticate(email, password string) (*UserModel, error)
	UserDB
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := utils.NewHMAC(hmacKey)
	uv := NewUserValidator(ug, hmac, pepper)

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

func (us *userService) Authenticate(email, password string) (*UserModel, error) {
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

// UserAvatarService is a set of methods used to manipulate and work with the user avatar model
type UserAvatarService interface {
	UserAvatarDB
}

func NewUserAvatarService(db *gorm.DB) UserAvatarService {
	ug := &userAvatarGorm{db}
	uv := MewUserAvatarValidator(ug)

	return &userAvatarService{
		UserAvatarDB: uv,
	}
}

var _ UserAvatarService = &userAvatarService{}

type userAvatarService struct {
	UserAvatarDB
}

func (uas *userAvatarService) Create(userID uint, r io.ReadCloser, filename string) error {
	defer func() {
		hub.ErrorHandler(r.Close())
	}()

	path, err := uas.mkImagePath(userID)
	if err != nil {
		return err
	}

	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer func() {
		hub.ErrorHandler(dst.Close())
	}()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	err = uas.UserAvatarDB.Create(userID, r, filename)
	if err != nil {
		return err
	}
	return nil
}

func (uas *userAvatarService) mkImagePath(userID uint) (string, error) {
	userPath := fmt.Sprintf("apps/%s/assets/images/%v/", SubAppName, userID)
	err := os.MkdirAll(userPath, 0755)
	if err != nil {
		return "", err
	}
	return userPath, nil
}
