package main

import (
	"camp/apps/accounts"
	"camp/apps/blog"
	"camp/apps/core"
	"camp/core/web"
)

func main() {
	cfg := web.LoadConfig()
	db, err := web.NewDbConnection(cfg.Database.Dialect(), cfg.Database.ConnectionInfo(), cfg.IsProd())
	if err != nil {
		panic(err)
	}

	apps := []web.SubApp{
		accounts.NewSubApp(db, cfg),
		blog.NewSubApp(db, cfg),
		core.NewSubApp(db, cfg),
	}

	//if err := db.DestructiveReset(); err != nil {
	//	panic(err)
	//}

	if err := db.AutoMigrate(); err != nil {
		panic(err)
	}

	defer func(db *web.DB) {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}(db)

	webApp := web.NewApp(apps, cfg)
	webApp.Run()
}
