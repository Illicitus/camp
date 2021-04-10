package blog

import (
	"camp/apps/accounts"
	"camp/core/utils"
	"camp/core/web"
)

const SubAppName = "blog"

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub(SubAppName, cfg.IsProd())

var _ web.SubApp = &SubApp{}

type SubApp struct {
	ac                    *ArticleController
	RequireUserMiddleware accounts.RequireUser
}

func NewSubApp(db *web.DB, cfg *web.AppConfig, rm accounts.RequireUser) *SubApp {
	app := &SubApp{
		ac:                    NewArticleController(db, cfg),
		RequireUserMiddleware: rm,
	}
	hub.ErrorHandler(app.CollectModels(db))

	return app
}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{
		&ArticleModel{}, &CommentModel{},
	}
	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
