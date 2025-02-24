package koakumachan

import (
	"context"
	"errors"
	"gokoakumachan/koakumachan/command"
	"gokoakumachan/koakumachan/tarotdata"
	"log"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gopkg.in/yaml.v3"
)

var ErrBotCredentialFileNotSet = errors.New("bot credential file not set")

const HELP_MESSAGE = `
/help - 顯示幫助信息
/moon - 顯示目前月相
/tarot - 查詢托特塔羅牌牌義
`

type BotCredential struct {
	Token string `yaml:"token"`
}

type AppConfig struct {
	// When set, the bot runs in debug mode.
	Debug bool `yaml:"debug"`

	// When set, only the users with these usernames are allowed to use the bot.
	UsernameWhitelist []string `yaml:"username_whitelist"`

	// Telegram bot token.
	BotCredentialFile string `yaml:"bot_credential_file"`

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

	if len(conf.UsernameWhitelist) > 0 {
		opts = append(opts, bot.WithMiddlewares(app.checkWhitelist))
	}

	if conf.BotCredentialFile == "" {
		return nil, ErrBotCredentialFileNotSet
	}

	bs, err := os.ReadFile(conf.BotCredentialFile)
	if err != nil {
		return nil, err
	}

	var botCredential BotCredential
	err = yaml.Unmarshal(bs, &botCredential)
	if err != nil {
		return nil, err
	}

	tgbot, err := bot.New(botCredential.Token, opts...)
	if err != nil {
		return nil, err
	}

	tgbot.RegisterHandler(
		bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, app.handleHelp)

	moonCommand := &command.MoonCommand{}
	moonCommand.Register(tgbot)

	tarotCommand := &command.TarotCommand{TarotDeck: tarotDeck}
	tarotCommand.Register(tgbot)

	app.bot = tgbot

	return app, nil
}

func (app *App) Start(ctx context.Context) {
	app.bot.Start(ctx)
}

func (app *App) checkWhitelist(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		username := ""
		if update.Message != nil {
			username = update.Message.From.Username
		} else if update.CallbackQuery != nil {
			username = update.CallbackQuery.From.Username
		}

		if len(app.UsernameWhitelist) > 0 {
			if username == "" || !app.UsernameWhitelist[username] {
				log.Printf("Unauthorized access from %s", username)
				return
			}
		}

		next(ctx, b, update)
	}
}

func (app *App) handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	app.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   HELP_MESSAGE,
	})
}
