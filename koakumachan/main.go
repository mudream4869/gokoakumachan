package koakumachan

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gopkg.in/yaml.v3"
)

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

	app, err := NewApp(conf)
	if err != nil {
		return fmt.Errorf("Main: %w", err)
	}

	app.Start(ctx)

	return nil
}
