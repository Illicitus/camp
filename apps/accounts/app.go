package accounts

import (
	"camp/core/utils"
	"camp/core/web"
)

const SubAppName = "accounts"

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub(SubAppName, cfg.IsProd())

var _ web.SubApp = &SubApp{}

type SubApp struct {
	uc                    *UserController
	RequireUserMiddleware RequireUser
}

func NewSubApp(db *web.DB, cfg *web.AppConfig) *SubApp {
	app := &SubApp{
		uc: NewController(db, cfg),
	}
	app.RequireUserMiddleware = RequireUser{
		User: User{
			UserService: app.uc.us,
		},
	}

	hub.ErrorHandler(app.CollectModels(db))

	return app
}

func (s *SubApp) CollectModels(db *web.DB) error {
	models := []web.Model{
		&UserModel{}, &UserAvatarModel{},
	}
	for _, m := range models {
		db.Models = append(db.Models, m)
	}
	return nil
}
