package core

import (
	"camp/core/web"
	"github.com/gorilla/mux"
)

var _ web.SubApp = &SubApp{}

type SubApp struct {
	c *Controller
}

func NewSubApp(db *web.DB, cfg web.AppConfig) *SubApp {
	app := &SubApp{
		c: NewController(db, cfg),
	}
	if err := app.CollectModels(db); err != nil {
		panic(err)
	}
	return app
}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{}

	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
func (sa SubApp) CollectRoutes(r *mux.Router) {
	r.Handle("/", sa.c.Home).Methods("GET")
	r.Handle("/about", sa.c.About).Methods("GET")
	r.Handle("/contact", sa.c.Contact).Methods("GET")
	r.Handle("/more-info", sa.c.MoreInfo).Methods("GET")
}
