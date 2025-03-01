package kastate

import (
	"github.com/go-telegram/bot"
	"github.com/mudream4869/gokoakumachan/koakumachan/kautil"
)

var State = kautil.NewLockMap[int64, bot.HandlerFunc]()
