package command

import (
	_ "image/png"

	_ "golang.org/x/image/webp"

	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/mudream4869/gokoakumachan/koakumachan/kastate"
)

const COVER_SIZE_LIMIT = 300 * 1024 // 300KB

type CoverCommand struct {
}

func (c *CoverCommand) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText,
		"/cover_from_url", bot.MatchTypePrefix, c.handleCoverFromURL)
}

func (c *CoverCommand) handleCoverFromURL(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "請輸入圖片網址",
	})

	kastate.State.Set(chatID, c.handleCoverFromURL1)
}

func (c *CoverCommand) handleCoverFromURL1(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	defer kastate.State.Delete(chatID)

	url := update.Message.Text

	resp, err := http.Get(url)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "無法取得圖片: " + err.Error(),
		})
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	img, filename, err := image.Decode(io.LimitReader(resp.Body, COVER_SIZE_LIMIT))
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "無法解析圖片: " + err.Error(),
		})
		log.Println(err)
		return
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "無法編碼圖片成jpg: " + err.Error(),
		})
		log.Println(err)
		return
	}

	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.Message.Chat.ID,
		Photo: &models.InputFileUpload{
			Data:     bytes.NewReader(buf.Bytes()),
			Filename: filename,
		},
	})
}
