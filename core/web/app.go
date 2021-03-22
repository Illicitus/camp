package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type SubApp interface {
	CollectRoutes(r *mux.Router)
	CollectModels(db *DB) error
}

type App struct {
	cfg               AppConfig
	appsRootDirectory string
	apps              []SubApp
	router            *mux.Router
}

func NewApp(apps []SubApp, cfg AppConfig) *App {
	r := mux.NewRouter()

	for _, app := range apps {
		app.CollectRoutes(r)
	}

	return &App{
		cfg:    cfg,
		apps:   apps,
		router: r,
	}
}

func (a *App) Run() {
	fmt.Printf("Started on:%d...\n", a.cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Port), a.router); err != nil {
		panic(err)
	}
}
