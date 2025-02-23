package koakumachan

import (
	"bytes"
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

	tgbot.RegisterHandler(
		bot.HandlerTypeMessageText, "/tarot", bot.MatchTypeExact, app.handleTarot)

	tgbot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData, "tarot_type", bot.MatchTypePrefix, app.handleTarotType)

	tgbot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData, "tarot_card", bot.MatchTypePrefix, app.handleTarotCard)

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

func (app *App) handleTarot(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if !app.checkWhitelist(update.Message.From.Username) {
		log.Printf("Unauthorized access from %s", update.Message.From.Username)
		return
	}

	// reply inline button
	app.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "請選擇類別",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "大阿卡那", CallbackData: "tarot_type.major_arcana"},
					{Text: "小阿卡那", CallbackData: "tarot_type.minor_arcana"},
					{Text: "宮廷", CallbackData: "tarot_type.court"},
				},
			},
		},
	})
}

func (app *App) handleTarotType(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	if !app.checkWhitelist(update.CallbackQuery.From.Username) {
		log.Printf("Unauthorized access from %s", update.CallbackQuery.From.Username)
		return
	}

	query := update.CallbackQuery.Data
	tarotType := query[len("tarot_type."):]

	log.Println("tarot type:", tarotType)

	var buttons [][]models.InlineKeyboardButton

	switch tarotType {
	case "major_arcana":
		for range 4 {
			buttons = append(buttons, []models.InlineKeyboardButton{})
		}

		for i, card := range app.tarotDeck.MajorArcana {
			lineIdx := i / 7
			buttons[lineIdx] = append(buttons[lineIdx], models.InlineKeyboardButton{
				Text:         card.Title,
				CallbackData: "tarot_card." + card.Name,
			})
		}

	case "minor_arcana":
		for range 10 {
			buttons = append(buttons, []models.InlineKeyboardButton{})
		}

		for i, card := range app.tarotDeck.MinorArcana {
			lineIdx := i % 10
			buttons[lineIdx] = append(buttons[lineIdx], models.InlineKeyboardButton{
				Text:         card.Title,
				CallbackData: "tarot_card." + card.Name,
			})
		}

	case "court":
		for range 4 {
			buttons = append(buttons, []models.InlineKeyboardButton{})
		}

		for i, card := range app.tarotDeck.Court {
			lineIdx := i % 4
			buttons[lineIdx] = append(buttons[lineIdx], models.InlineKeyboardButton{
				Text:         card.Title,
				CallbackData: "tarot_card." + card.Name,
			})
		}

	default:
		log.Printf("Unknown tarot type: %s", tarotType)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "請選擇卡",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})
}

func (app *App) handleTarotCard(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	if !app.checkWhitelist(update.CallbackQuery.From.Username) {
		log.Printf("Unauthorized access from %s", update.CallbackQuery.From.Username)
		return
	}

	query := update.CallbackQuery.Data
	cardName := query[len("tarot_card."):]

	log.Println("tarot card:", cardName)

	card := app.tarotDeck.Find(cardName)
	if card == nil {
		log.Printf("Unknown card: %s", cardName)
		return
	}

	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Photo: &models.InputFileUpload{
			Data: bytes.NewReader(card.ImageData),
		},
		Caption: card.Description,
	})
}
