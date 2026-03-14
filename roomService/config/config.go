package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	ServiceApiKey string `mapstructure:"service_api_key"`
	JwtSecret     string `mapstructure:"jwt_secret"`
}

type Config struct {
	Env    string       `mapstructure:"env"`
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

var Cfg Config

func Load() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("env", "dev")
	viper.SetDefault("server.port", ":50052")
	viper.SetDefault("auth.enabled", true)
	viper.SetDefault("auth.service_api_key", "")
	viper.SetDefault("auth.jwt_secret", "")

	viper.BindEnv("env", "ENV")
	viper.BindEnv("server.port", "LISTEN_PORT")
	viper.BindEnv("auth.enabled", "AUTH")
	viper.BindEnv("auth.service_api_key", "SERVICE_API_KEY")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading config file: %v", err)
		}
		log.Println("no config file found, using defaults and environment variables")
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}

	if Cfg.Auth.Enabled && Cfg.Auth.ServiceApiKey == "" {
		panic("auth is enabled but service api key is not provided")
	}
	if Cfg.Auth.Enabled && Cfg.Auth.JwtSecret == "" {
		panic("auth is enabled but jwt secret is not provided")
	}
}
