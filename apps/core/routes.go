package core

import (
	"github.com/gorilla/mux"
)

func (s SubApp) CollectRoutes(r *mux.Router) {
	r.Handle("/", s.c.Home).Methods("GET")
	r.Handle("/about", s.c.About).Methods("GET")
	r.Handle("/contact", s.c.Contact).Methods("GET")
	r.Handle("/more-info", s.c.MoreInfo).Methods("GET")
}
