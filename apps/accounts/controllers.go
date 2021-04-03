package accounts

import (
	"camp/core/utils"
	"camp/core/web"
	"log"
	"net/http"
)

var (
	LayoutDir   = "apps/core/views/layouts/"
	TemplateDir = "apps/accounts/views/"
)

type UserController struct {
	SignUpView *web.View
	LoginView  *web.View
	us         UserService
}

func NewController(db *web.DB, cfg *web.AppConfig) *UserController {
	return &UserController{
		SignUpView: web.NewView(TemplateDir, LayoutDir, "bootstrap", "new"),
		LoginView:  web.NewView(TemplateDir, LayoutDir, "bootstrap", "login"),
		us:         NewUserService(db.Conn, cfg.Pepper, cfg.HMACKey),
	}
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var vd web.Data
	var form SignupForm
	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.SignUpView.Render(w, r, vd))
		return
	}

	user := User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := uc.us.Create(&user); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.SignUpView.Render(w, r, vd))
		return
	}
	err := uc.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/accounts/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	vd := web.Data{}
	var form LoginForm

	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.LoginView.Render(w, r, vd))
		return
	}

	user, err := uc.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case utils.GormErr.NotFound:
			vd.Alert = &web.Alert{
				Level:   web.AlertLvlDanger,
				Message: "Invalid email address.",
			}
		case utils.GormErr.InvalidPassword:
			vd.Alert = &web.Alert{
				Level:   web.AlertLvlDanger,
				Message: "Invalid password provided.",
			}
		default:
			vd.SetAlert(err)
		}
		hub.ErrorHandler(uc.LoginView.Render(w, r, vd))
		return
	}

	err = uc.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		hub.ErrorHandler(uc.LoginView.Render(w, r, vd))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (uc *UserController) signIn(w http.ResponseWriter, user *User) error {
	if user.Remember == "" {
		token, err := utils.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = uc.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
