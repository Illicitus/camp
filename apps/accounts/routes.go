package accounts

import (
	"github.com/gorilla/mux"
)

func (s *SubApp) CollectRoutes(r *mux.Router) {
	r.Handle("/accounts/signup", s.uc.SignUpView).Methods("GET")
	r.HandleFunc("/accounts/signup", s.uc.Create).Methods("POST")
	r.Handle("/accounts/login", s.uc.LoginView).Methods("GET")
	r.HandleFunc("/accounts/login", s.uc.Login).Methods("POST")
}
