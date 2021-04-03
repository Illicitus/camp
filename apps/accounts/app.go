package accounts

import (
	"camp/core/utils"
	"camp/core/web"
)

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub("accounts", cfg.IsProd())

var _ web.SubApp = &SubApp{}

type SubApp struct {
	uc *UserController
}

func NewSubApp(db *web.DB, cfg *web.AppConfig) *SubApp {
	app := &SubApp{
		uc: NewController(db, cfg),
	}
	hub.ErrorHandler(app.CollectModels(db))

	return app
}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{
		&User{},
	}
	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
