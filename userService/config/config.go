package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	AuthEnabled bool `mapstructure:"auth"`
	ServiceApiKey string `mapstructure:"service_api_key"`
}

type AuthConfig struct {
	PasswordSalt         string `mapstructure:"password_salt"`
	JwtKey               string `mapstructure:"jwt_secret"`
	AuthTokenDuration    uint   `mapstructure:"auth_token_duration"`
	RefreshTokenDuration uint   `mapstructure:"refresh_token_duration"`
	Issuer               string `mapstructure:"issuer"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

var Cfg Config

func Load() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", ":6741")
	viper.SetDefault("server.auth", true)
	viper.SetDefault("server.service_api_key", "cute")
	viper.SetDefault("auth.password_salt", "kosher")
	viper.SetDefault("auth.jwt_key", " ")
	viper.SetDefault("auth.auth_token_duration", 10)
	viper.SetDefault("auth.refresh_token_duration", 1440)
	viper.SetDefault("auth.issuer", "loginService")

	viper.BindEnv("server.port", "LISTEN_PORT")
	viper.BindEnv("server.service_api_key", "SERVICE_API_KEY")
	viper.BindEnv("server.enabled", "AUTH")
	viper.BindEnv("auth.password_salt", "PASSWORD_SALT")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("auth.auth_token_duration", "AUTH_TOKEN_DURATION")
	viper.BindEnv("auth.refresh_token_duration", "REFRESH_TOKEN_DURATION")
	viper.BindEnv("auth.issuer", "ISSUER")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading config file: %v", err)
		}
		log.Println("no config file found, using defaults and environment variables")
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}

	if Cfg.Auth.JwtKey == " " {
		log.Fatal("jwt key is not set. Try adding JWT_SECRET env variable")
	}
}
