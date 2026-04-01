package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	JwtSecret string `mapstructure:"jwt_secret"`
	TokenLifespan uint `mapstructure:"token_lifespan"` // in minutes 
}

type ServiceEndpoint struct {
	Address string `mapstructure:"address"`
	ApiKey  string `mapstructure:"api_key"`
}

type ServicesConfig struct {
	Message ServiceEndpoint `mapstructure:"message"`
	Room    ServiceEndpoint `mapstructure:"room"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Services ServicesConfig `mapstructure:"services"`
}

var Cfg Config

func Load() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("auth.enabled", true)
	viper.SetDefault("auth.jwt_secret", "")
	viper.SetDefault("auth.token_lifespan", 10)
	viper.SetDefault("services.message.address", "message-service:50051")
	viper.SetDefault("services.message.api_key", "cute")
	viper.SetDefault("services.room.address", "localhost:50052")
	viper.SetDefault("services.room.api_key", "cute")

	viper.BindEnv("server.port", "GATEWAY_PORT")
	viper.BindEnv("auth.enabled", "AUTH")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("services.message.address", "MESSAGE_SERVICE_IP")
	viper.BindEnv("services.message.api_key", "MESSAGE_SERVICE_API_KEY")
	viper.BindEnv("services.room.address", "ROOM_SERVICE_IP")
	viper.BindEnv("services.room.api_key", "ROOM_SERVICE_API_KEY")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("error reading config file: %v", err)
		}
		log.Println("no config file found, using defaults and environment variable")
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("unable to unmarshal config: %v", err)
	}
}
