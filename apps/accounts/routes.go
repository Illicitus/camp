package accounts

import (
	"github.com/gorilla/mux"
)

func (s *SubApp) CollectRoutes(r *mux.Router) {
	r.HandleFunc("/"+SubAppName+"/signup", s.uc.SignUpPage).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/signup", s.uc.Create).Methods("POST")
	r.HandleFunc("/"+SubAppName+"/login", s.uc.LoginPage).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/login", s.uc.Login).Methods("POST")
	r.HandleFunc("/"+SubAppName+"/logout", s.RequireUserMiddleware.ApplyFn(s.uc.Logout)).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/profile", s.RequireUserMiddleware.ApplyFn(s.uc.ProfilePage)).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/avatar", s.RequireUserMiddleware.ApplyFn(s.uc.UpdateAvatar)).Methods("POST")
	r.HandleFunc("/"+SubAppName+"/{id:[0-9]+}/update", s.RequireUserMiddleware.ApplyFn(s.uc.UpdatePage)).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/{id:[0-9]+}/update", s.RequireUserMiddleware.ApplyFn(s.uc.Update)).Methods("POST")
}
