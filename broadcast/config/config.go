package config

import (
	"log"

	"github.com/spf13/viper"
)

type AuthConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	JwtSecret string `mapstructure:"jwt_secret"`
}

type NatsConfig struct {
	Address string `mapstructure:"address"`
}

type Config struct {
	Auth AuthConfig `mapstructure:"auth"`
	Nats NatsConfig `mapstructure:"nats"`
}

var Cfg Config

func Load() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("auth.enabled", true)
	viper.SetDefault("auth.jwt_secret", "")
	viper.SetDefault("nats.address", "localhost:4222")

	viper.BindEnv("auth.enabled", "AUTH")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("nats.address", "NATS_ADDRESS")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading config file: %v", err)
		}
		log.Println("no config file found, using defaults and environment variables")
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}
}
