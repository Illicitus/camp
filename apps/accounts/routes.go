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

	r.HandleFunc("/accounts/signup", s.uc.SignUpPage).Methods("GET")
	r.HandleFunc("/accounts/signup", s.uc.Create).Methods("POST")
	r.HandleFunc("/accounts/login", s.uc.LoginPage).Methods("GET")
	r.HandleFunc("/accounts/login", s.uc.Login).Methods("POST")
	r.HandleFunc("/accounts/logout", requireUserMiddleware.ApplyFn(s.uc.Logout)).Methods("GET")
	r.HandleFunc("/accounts/update", requireUserMiddleware.ApplyFn(s.uc.UpdatePage)).Methods("GET")
	r.HandleFunc("/accounts/update", requireUserMiddleware.ApplyFn(s.uc.Update)).Methods("POST")
}
