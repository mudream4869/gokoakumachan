package command

import (
	"bytes"
	"context"
	"log"

	"github.com/mudream4869/gokoakumachan/koakumachan/tarotdata"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TarotCommand struct {
	TarotDeck *tarotdata.Deck
}

func (c *TarotCommand) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText,
		"/tarot", bot.MatchTypeExact, c.handleTarot)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData,
		"tarot_type", bot.MatchTypePrefix, c.handleTarotType)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData,
		"tarot_card", bot.MatchTypePrefix, c.handleTarotCard)
}

func (c *TarotCommand) handleTarot(ctx context.Context, b *bot.Bot, update *models.Update) {
	// reply inline button
	b.SendMessage(ctx, &bot.SendMessageParams{
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

func (c *TarotCommand) handleTarotType(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	query := update.CallbackQuery.Data
	tarotType := query[len("tarot_type."):]

	log.Println("tarot type:", tarotType)

	var buttons [][]models.InlineKeyboardButton

	switch tarotType {
	case "major_arcana":
		for range 4 {
			buttons = append(buttons, []models.InlineKeyboardButton{})
		}

		for i, card := range c.TarotDeck.MajorArcana {
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

		for i, card := range c.TarotDeck.MinorArcana {
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

		for i, card := range c.TarotDeck.Court {
			lineIdx := i % 4
			buttons[lineIdx] = append(buttons[lineIdx], models.InlineKeyboardButton{
				Text:         card.Title,
				CallbackData: "tarot_card." + card.Name,
			})
		}

	default:
		log.Printf("Unknown tarot type: %s", tarotType)
	}

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      "請選擇卡",
	})

	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})
}

func (c *TarotCommand) handleTarotCard(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	query := update.CallbackQuery.Data
	cardName := query[len("tarot_card."):]

	log.Println("tarot card:", cardName)

	card := c.TarotDeck.Find(cardName)
	if card == nil {
		log.Printf("Unknown card: %s", cardName)
		return
	}

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})

	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Photo: &models.InputFileUpload{
			Data: bytes.NewReader(card.ImageData),
		},
		Caption: card.Description,
	})
}
