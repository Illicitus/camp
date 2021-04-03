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

	// Generate web app
	apps := []web.SubApp{
		accounts.NewSubApp(db, cfg),
		blog.NewSubApp(db, cfg),
		core.NewSubApp(db, cfg),
	}

	//if err := db.DestructiveReset(); err != nil {
	//	panic(err)
	//}
	hub.ErrorHandler(db.AutoMigrate())

	webApp := web.NewApp(apps, cfg)
	webApp.Run()

}
