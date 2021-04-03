package accounts

import (
	"camp/core/web"
	"log"
	"net/http"
)

var (
	LayoutDir   string = "apps/core/views/layouts/"
	TemplateDir string = "apps/accounts/views/"
)

type UserController struct {
	NewView *web.View
	us      UserService
}

func NewController(db *web.DB, cfg *web.AppConfig) *UserController {
	return &UserController{
		NewView: web.NewView(TemplateDir, LayoutDir, "bootstrap", "new"),
		us:      NewUserService(db.Conn, cfg.Pepper, cfg.HMACKey),
	}
}

func (uc *UserController) New(w http.ResponseWriter, r *http.Request) {
	hub.ErrorHandler(uc.NewView.Render(w, r, nil))
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var vd web.Data
	var form SignupForm
	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(uc.NewView.Render(w, r, vd))
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
		hub.ErrorHandler(uc.NewView.Render(w, r, vd))
		return
	}
	//err := uc.signIn(w, &user)
	//if err != nil {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}
	http.Redirect(w, r, "/", http.StatusFound)
}
