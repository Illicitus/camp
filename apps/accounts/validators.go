package accounts

import (
	"camp/core/utils"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

var (
	emailRegex        = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	passwordMinLength = 8
	rememberMinBytes  = 8
)

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac       utils.HMAC
	emailRegex *regexp.Regexp
	pepper     string
}

func newUserValidator(udb UserDB, hmac utils.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: emailRegex,
		pepper:     pepper,
	}
}

func (uv *userValidator) defaultRemember(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := utils.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token

	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := utils.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < rememberMinBytes {
		return utils.ValErr.RememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return utils.ValErr.RememberRequired
	}
	return nil
}

func (uv *userValidator) idGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return utils.ValErr.InvalidID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return utils.ValErr.EmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return utils.ValErr.InvalidEmail
	}
	return nil
}

func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)

	// Email doesn't exist, no error to user
	if err == utils.GormErr.NotFound {
		return nil
	}
	if err != nil {
		return err
	}
	// Email exists, return validation error then
	if user.ID != existing.ID {
		return utils.ValErr.EmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < passwordMinLength {
		return utils.ValErr.PasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return utils.ValErr.PasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return utils.ValErr.PasswordRequired
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)

	return nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return nil
}

func (uv *userValidator) Create(user *User) error {
	valFuncs := []userValFunc{
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.defaultRemember,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable,
	}
	if err := runUserValFuncs(user, valFuncs...); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{Email: email}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}
