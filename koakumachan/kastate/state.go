package kastate

import (
	"github.com/go-telegram/bot"
	"github.com/mudream4869/gokoakumachan/koakumachan/kautil"
)

// State is a lock map that mapping chat id to handler function
var State = kautil.NewLockMap[int64, bot.HandlerFunc]()
