package command

import (
	"context"
	"time"

	"github.com/mudream4869/gokoakumachan/koakumachan/moonphase"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var PHASE_EMOJI = map[string]string{
	"New Moon":        "🌑",
	"Waxing Crescent": "🌒",
	"First Quarter":   "🌓",
	"Waxing Gibbous":  "🌔",
	"Full Moon":       "🌕",
	"Waning Gibbous":  "🌖",
	"Third Quarter":   "🌗",
	"Waning Crescent": "🌘",
}

type MoonCommand struct {
}

func (c *MoonCommand) Register(b *bot.Bot) {
	b.RegisterHandler(
		bot.HandlerTypeMessageText, "/moon", bot.MatchTypeExact, c.HandleMoon)
}

func (c *MoonCommand) HandleMoon(ctx context.Context, b *bot.Bot, update *models.Update) {
	phase := moonphase.New(time.Now()).PhaseName()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   PHASE_EMOJI[phase] + " " + phase,
	})
}
