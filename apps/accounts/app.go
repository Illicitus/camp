package accounts

import (
	"camp/core/utils"
	"camp/core/web"
)

var cfg = web.LoadConfig()
var hub = utils.NewLocalHub("accounts", cfg.IsProd())
