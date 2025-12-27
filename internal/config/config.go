package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Telegram struct {
		Token string `mapstructure:"token"`
		ChatIDs []int64 `mapstructure:"chat_ids"`
	} `mapstructure:"telegram"`
}

func Load(env string) (*Config, error) {
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yml")
	viper.AddConfigPath("../config/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Can't read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Can't parse сonfig: %w", err)
	}

	if config.Telegram.Token == "" {
		return nil, fmt.Errorf("Can't find Telegram bot token")
	}

	if len(config.Telegram.ChatIDs) == 0 {
		return nil, fmt.Errorf("Chat ids is empty")
	}

	log.Printf("Loaded configuration from: %s", viper.ConfigFileUsed())

	return &config, nil
}