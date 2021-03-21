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
	New *web.View
	us  UserService
}

func NewController(db *web.DB, cfg web.AppConfig) *UserController {
	return &UserController{
		New: web.NewView(TemplateDir, LayoutDir, "bootstrap", "new"),
		us:  NewUserService(db.Conn, cfg.Pepper, cfg.HMACKey),
	}
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var vd web.Data
	var form SignupForm
	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		if err := uc.New.Render(w, r, vd); err != nil {
			panic(err)
		}
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
		if err := uc.New.Render(w, r, vd); err != nil {
			panic(err)
		}
		return
	}
	//err := uc.signIn(w, &user)
	//if err != nil {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}
	http.Redirect(w, r, "/", http.StatusFound)
}
