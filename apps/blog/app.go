package blog

import (
	"camp/core/utils"
	"camp/core/web"
)

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub("blog", cfg.IsProd())

var _ web.SubApp = &SubApp{}

type SubApp struct{}

func NewSubApp(db *web.DB, cfg *web.AppConfig) *SubApp {
	app := &SubApp{
		//uc: NewController(db, cfg),
	}
	if err := app.CollectModels(db); err != nil {
		panic(err)
	}
	return app
}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{
		&Article{}, &Comment{},
	}
	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
