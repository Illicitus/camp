package blog

import (
	"camp/core/web"
	"github.com/gorilla/mux"
)

var _ web.SubApp = &SubApp{}

type SubApp struct{}

func NewSubApp(db *web.DB, cfg web.AppConfig) *SubApp {
	app := &SubApp{
		//uc: NewController(db, cfg),
	}
	if err := app.CollectModels(db); err != nil {
		panic(err)
	}
	return app
}
func (s SubApp) CollectRoutes(r *mux.Router) {}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{
		&Article{}, &Comment{},
	}
	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
