package koakumachan

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	BotToken string `yaml:"bot_token"`
}

func readConfig(filename string) (*AppConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("readConfig: %w", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("readConfig: %w", err)
	}

	return &cfg, nil
}

func Main(confFilename string) error {
	conf, err := readConfig(confFilename)
	if err != nil {
		return fmt.Errorf("Main: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(conf.BotToken, opts...)
	if err != nil {
		return fmt.Errorf("Main: %w", err)
	}

	b.Start(ctx)
	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
