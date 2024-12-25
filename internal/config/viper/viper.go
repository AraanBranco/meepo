package viper

import (
	"fmt"
	"strings"

	"github.com/AraanBranco/meepo/internal/config"
	"github.com/spf13/viper"
)

func NewViperConfig(configPath string) (config.Config, error) {
	config := viper.New()
	config.SetEnvPrefix("meepo")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	config.SetConfigType("yaml")
	config.SetConfigFile(configPath)
	err := config.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}
