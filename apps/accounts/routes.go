package accounts

import (
	"github.com/gorilla/mux"
)

func (s *SubApp) CollectRoutes(r *mux.Router) {

	requireUserMiddleware := RequireUser{
		User: User{
			UserService: s.uc.us,
		},
	}

	r.HandleFunc("/"+SubAppName+"/signup", s.uc.SignUpPage).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/signup", s.uc.Create).Methods("POST")
	r.HandleFunc("/"+SubAppName+"/login", s.uc.LoginPage).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/login", s.uc.Login).Methods("POST")
	r.HandleFunc("/"+SubAppName+"/logout", requireUserMiddleware.ApplyFn(s.uc.Logout)).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/update", requireUserMiddleware.ApplyFn(s.uc.UpdatePage)).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/update", requireUserMiddleware.ApplyFn(s.uc.Update)).Methods("POST")
}
