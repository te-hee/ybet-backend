package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var(
	PasswordSalt *string
	JwtKey *string
	ListenPort *string
	AuthTokenDuration *uint // IN MINUTES
	RefreshTokenDuration *uint // IN MINUTES
	Issuer *string
)


func LoadConfig() error{
	PasswordSalt = flag.String("password-salt", "kosher", "Password salt used for hashing")
	JwtKey = flag.String("jwt-key", "kosher", "Jwt encode key")
	ListenPort = flag.String("port", ":6741", "Port On which server will listen for requests")
	AuthTokenDuration = flag.Uint("auth-duration", 10, "Duration IN MINUTES of auth token")
	RefreshTokenDuration = flag.Uint("refresh-duration", 60*24, "Duration IN MINUTES of refresh token")
	Issuer = flag.String("issuer", "loginService", "Isser name used in JWT tokens")

	err := loadEnvs()

	if err != nil{
		return err
	}

	flag.Parse()
	return nil
}

func loadEnvs() (error){
	err := godotenv.Load()
	if err != nil{
		return err
	}


	if value, ok := os.LookupEnv("PASSWORD_SALT"); ok{
		*PasswordSalt = value	
	}
	if value, ok := os.LookupEnv("JWT_KEY"); ok{
		*JwtKey = value	
	}
	if value, ok := os.LookupEnv("LISTEN_PORT"); ok{
		*ListenPort = value	
	}
	if value, ok := os.LookupEnv("AUTH_TOKEN_DURATION"); ok{
		duration, err := strconv.ParseUint(value, 10, 64)

		if err != nil{
			return err
		}

		*AuthTokenDuration = uint(duration)	
	}
	if value, ok := os.LookupEnv("REFRESH_TOKEN_DURATION"); ok{
		duration, err := strconv.ParseUint(value, 10, 64)

		if err != nil{
			return err
		}

		*RefreshTokenDuration = uint(duration)	
	}
	if value, ok := os.LookupEnv("ISSUER"); ok{
		*Issuer = value
	}

	log.Println("Password salt: ", *PasswordSalt)
	log.Println("Jwt key: ", *JwtKey)
	log.Println("Listen port: ", *ListenPort)
	log.Println("Auth token duration (in minutes): ", *AuthTokenDuration)
	log.Println("Refresh token duration (in minutes): ", *RefreshTokenDuration)
	log.Println("Issuer: ", *Issuer)

	return nil
}
