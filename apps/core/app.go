package core

import (
	"camp/core/utils"
	"camp/core/web"
)

const SubAppName = "core"

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub(SubAppName, cfg.IsProd())

var _ web.SubApp = &SubApp{}

type SubApp struct {
	c *Controller
}

func NewSubApp(db *web.DB, cfg *web.AppConfig) *SubApp {
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
