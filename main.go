package main

import (
	"camp/apps/accounts"
	"camp/apps/blog"
	"camp/apps/core"
	"camp/core/utils"
	"camp/core/web"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func main() {
	cfg := web.LoadConfig()

	// Init sentry
	err := sentry.Init(sentry.ClientOptions{Dsn: cfg.Sentry.Dns})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	hub := utils.NewLocalHub("main", cfg.IsProd())

	defer sentry.Flush(time.Second * 5)
	defer sentry.Recover()

	// Init db
	db, err := web.NewDbConnection(cfg.Database.Dialect, cfg.Database.ConnectionInfo(), cfg.IsProd())
	hub.ErrorHandler(err)

	defer func(db *web.DB) {
		hub.ErrorHandler(db.Close())
	}(db)

	// Generate web sub app
	accountsSubApp := accounts.NewSubApp(db, cfg)
	apps := []web.SubApp{
		accountsSubApp,
		blog.NewSubApp(db, cfg, accountsSubApp.RequireUserMiddleware),
		core.NewSubApp(db, cfg, accountsSubApp.RequireUserMiddleware),
	}

	//if err := db.DestructiveReset(); err != nil {
	//	panic(err)
	//}
	hub.ErrorHandler(db.AutoMigrate())

	webApp := web.NewApp(apps, cfg)
	webApp.Run()

}
