package accounts

import (
	"camp/core/utils"
	"camp/core/web"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	LayoutDir   = "apps/core/views/layouts/"
	TemplateDir = "apps/accounts/views/"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte (2 << 20 = 2 mb), 2 << 10 = 1 kb
)

type UserController struct {
	SignUpView  *web.View
	LoginView   *web.View
	UpdateView  *web.View
	ProfileView *web.View
	us          UserService
	uas         UserAvatarService
}

func NewController(db *web.DB, cfg *web.AppConfig) *UserController {
	return &UserController{
		SignUpView:  web.NewView(TemplateDir, LayoutDir, "bootstrap", "new"),
		LoginView:   web.NewView(TemplateDir, LayoutDir, "bootstrap", "login"),
		UpdateView:  web.NewView(TemplateDir, LayoutDir, "bootstrap", "update"),
		ProfileView: web.NewView(TemplateDir, LayoutDir, "bootstrap", "profile"),
		us:          NewUserService(db.Conn, cfg.Pepper, cfg.HMACKey),
		uas:         NewUserAvatarService(db.Conn),
	}
}

func (uc *UserController) SignUpPage(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	// Ignore parse url errors
	hub.IgnoreErrorHandler(web.ParseURLParams(r, &form))
	hub.ErrorHandler(uc.SignUpView.Render(w, r, form))
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var vd web.Data
	var form SignupForm
	vd.Yield = &form
	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.SignUpView.Render(w, r, vd))
		return
	}

	user := UserModel{
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

func (uc *UserController) LoginPage(w http.ResponseWriter, r *http.Request) {
	var form LoginForm

	// Ignore parse url errors
	hub.IgnoreErrorHandler(web.ParseURLParams(r, &form))
	hub.ErrorHandler(uc.LoginView.Render(w, r, form))
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	vd := web.Data{}
	var form LoginForm
	vd.Yield = &form
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

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	user := UserContext(r.Context())
	token, _ := utils.RememberToken()
	user.Remember = token
	hub.ErrorHandler(uc.us.Update(user))
	http.Redirect(w, r, "/", http.StatusFound)
}

func (uc *UserController) signIn(w http.ResponseWriter, user *UserModel) error {
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
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(user.ID)),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}

func (uc *UserController) userByID(w http.ResponseWriter, r *http.Request) (*UserModel, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return nil, err
	}
	user, err := uc.us.ByID(uint(id))
	if err != nil {
		switch err {
		case utils.GormErr.NotFound:
			http.Error(w, "Invalid user ID", http.StatusNotFound)
		default:
			http.Error(w, "Something wend wrong", http.StatusInternalServerError)
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserController) userByRemember(w http.ResponseWriter, r *http.Request) (*UserModel, error) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		return nil, err
	}
	user, err := uc.us.ByRemember(cookie.Value)
	if err != nil {
		switch err {
		case utils.GormErr.NotFound:
			http.Error(w, "Invalid user ID", http.StatusNotFound)
		default:
			http.Error(w, "Something wend wrong", http.StatusInternalServerError)
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserController) ProfilePage(w http.ResponseWriter, r *http.Request) {
	user, err := uc.userByRemember(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uc.us.ProfileByUserID(user.ID)
	vd := web.Data{}
	vd.Yield = &user

	// Ignore parse url errors
	hub.IgnoreErrorHandler(web.ParseURLParams(r, vd))
	hub.ErrorHandler(uc.ProfileView.Render(w, r, vd))
}

func (uc *UserController) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user, err := uc.userByRemember(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vd := web.Data{}

	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		hub.ErrorHandler(uc.ProfileView.Render(w, r, vd))
		return
	}

	files := r.MultipartForm.File["avatar"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			hub.ErrorHandler(uc.ProfileView.Render(w, r, vd))
			return
		}
		defer func() {
			hub.ErrorHandler(file.Close())
		}()

		err = uc.uas.Create(user.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			hub.ErrorHandler(uc.ProfileView.Render(w, r, vd))
			return
		}
	}
	// Ignore parse url errors
	hub.IgnoreErrorHandler(web.ParseURLParams(r, vd))
	hub.ErrorHandler(uc.ProfileView.Render(w, r, vd))
}

func (uc *UserController) UpdatePage(w http.ResponseWriter, r *http.Request) {
	user, err := uc.userByID(w, r)
	if err != nil {
		return
	}
	vd := web.Data{}

	vd.Yield = &user

	// Ignore parse url errors
	hub.IgnoreErrorHandler(web.ParseURLParams(r, vd))
	hub.ErrorHandler(uc.UpdateView.Render(w, r, vd))
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	user := UserContext(r.Context())

	var form UpdateForm
	var vd web.Data
	vd.Yield = &user

	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.UpdateView.Render(w, r, vd))
		return
	}
	user.Name = form.Name
	err := uc.us.Update(user)
	if err != nil {
		vd.SetAlert(err)
		hub.ErrorHandler(uc.UpdateView.Render(w, r, vd))
		return
	}
	vd.Alert = &web.Alert{
		Level:   web.AlertLvlSuccess,
		Message: "User successfully updated!",
	}
	hub.ErrorHandler(uc.UpdateView.Render(w, r, vd))
}
