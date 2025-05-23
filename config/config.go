package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port           int      `mapstructure:"port"`
		Host           string   `mapstructure:"host"`
		TrustedProxies []string `mapstructure:"trusted_proxies"`
	} `mapstructure:"server"`
	Storage struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"storage"`
	Domain string `mapstructure:"domain"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	viper.SetDefault("server.port", 6060)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("storage.path", "./uploads")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
