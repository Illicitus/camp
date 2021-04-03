package blog

import (
	"camp/core/utils"
	"camp/core/web"
)

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub("blog", cfg.IsProd())
