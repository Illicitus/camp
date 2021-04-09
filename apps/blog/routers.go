package blog

import (
	"github.com/gorilla/mux"
)

func (s SubApp) CollectRoutes(r *mux.Router) {
	r.HandleFunc("/"+SubAppName+"/articles", s.ac.List).Methods("GET")
	r.Handle("/"+SubAppName+"/articles/new", s.ac.ArticleCreateView).Methods("GET")
	r.HandleFunc("/"+SubAppName+"/articles/new", s.ac.Create).Methods("POST")
	//r.HandleFunc("/"+SubAppName+"/articles/{id:[0-9]+}/comments", s.cc.Create).Methods("POST")
}
