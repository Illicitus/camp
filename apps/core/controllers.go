package core

import (
	"camp/core/web"
)

var (
	LayoutDir   string = "apps/core/views/layouts/"
	TemplateDir string = "apps/core/views/"
)

type Controller struct {
	About,
	Contact,
	Home,
	MoreInfo *web.View
}

func NewController(db *web.DB, cfg *web.AppConfig) *Controller {
	return &Controller{
		About:    web.NewView(TemplateDir, LayoutDir, "bootstrap", "about"),
		Contact:  web.NewView(TemplateDir, LayoutDir, "bootstrap", "contact"),
		Home:     web.NewView(TemplateDir, LayoutDir, "bootstrap", "home"),
		MoreInfo: web.NewView(TemplateDir, LayoutDir, "bootstrap", "more_info"),
	}
}
