package koakumachan

import (
	"context"
	"gokoakumachan/koakumachan/moonphase"
	"gokoakumachan/koakumachan/tarotdata"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const HELP_MESSAGE = `
/help - 顯示幫助信息
/moon - 顯示目前月相
/tarot - 開始占卜
`

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

type AppConfig struct {
	// When set, the bot runs in debug mode.
	Debug bool `yaml:"debug"`

	// When set, only the users with these usernames are allowed to use the bot.
	UsernameWhitelist []string `yaml:"username_whitelist"`

	// Telegram bot token.
	BotToken string `yaml:"bot_token"`

	// Filename of the tarot data.
	TarotDataFilename string `yaml:"tarot_data_filename"`
}

type App struct {
	bot               *bot.Bot
	tarotDeck         *tarotdata.Deck
	UsernameWhitelist map[string]bool
	conf              *AppConfig
}

func NewApp(conf *AppConfig) (*App, error) {
	tarotDeck, err := tarotdata.LoadDeck(conf.TarotDataFilename)
	if err != nil {
		return nil, err
	}

	usernameWhitelist := make(map[string]bool)
	for _, username := range conf.UsernameWhitelist {
		usernameWhitelist[username] = true
	}

	app := &App{
		tarotDeck:         tarotDeck,
		UsernameWhitelist: usernameWhitelist,
		conf:              conf,
	}

	opts := []bot.Option{}
	if conf.Debug {
		opts = append(opts, bot.WithDebug())
	}

	tgbot, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		return nil, err
	}

	tgbot.RegisterHandler(
		bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, app.handleHelp)

	tgbot.RegisterHandler(
		bot.HandlerTypeMessageText, "/moon", bot.MatchTypeExact, app.handleMoon)

	app.bot = tgbot

	return app, nil
}

func (app *App) Start(ctx context.Context) {
	app.bot.Start(ctx)
}

func (app *App) checkWhitelist(username string) bool {
	if len(app.UsernameWhitelist) == 0 {
		return true
	}

	return app.UsernameWhitelist[username]
}

func (app *App) handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !app.checkWhitelist(update.Message.From.Username) {
		log.Printf("Unauthorized access from %s", update.Message.From.Username)
		return
	}

	app.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   HELP_MESSAGE,
	})
}

func (app *App) handleMoon(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if !app.checkWhitelist(update.Message.From.Username) {
		log.Printf("Unauthorized access from %s", update.Message.From.Username)
		return
	}

	phase := moonphase.New(time.Now()).PhaseName()
	app.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   PHASE_EMOJI[phase] + " " + phase,
	})
}
