package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DATA_SOURCE_URL   string `mapstructure:"DATA_SOURCE_URL"`
	RABBIT_SOURCE_URL string `mapstructure:"RABBIT_SOURCE_URL"`
	MONGO_SOURCE_URL  string `mapstructure:"MONGO_SOURCE_URL"`
	APPLICATION_PORT  string `mapstructure:"APPLICATION_PORT"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
