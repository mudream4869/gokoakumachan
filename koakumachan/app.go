package koakumachan

import (
	"context"
	"gokoakumachan/koakumachan/tarotdata"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AppConfig struct {
	BotToken          string `yaml:"bot_token"`
	TarotDataFilename string `yaml:"tarot_data_filename"`
}

type App struct {
	bot       *bot.Bot
	tarotDeck *tarotdata.Deck
}

func NewApp(conf *AppConfig) (*App, error) {
	tarotDeck, err := tarotdata.LoadDeck(conf.TarotDataFilename)
	if err != nil {
		return nil, err
	}

	app := &App{
		tarotDeck: tarotDeck,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(app.echoHandler),
	}

	tgbot, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		return nil, err
	}

	app.bot = tgbot

	return app, nil
}

func (app *App) Start(ctx context.Context) {
	app.bot.Start(ctx)
}

func (app *App) echoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	app.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
