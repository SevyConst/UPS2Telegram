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

func Load(configFile string) (*Config, error) {

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Can't read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Can't parse config: %w", err)
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