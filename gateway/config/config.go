package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

var (
	MessageServiceAddr *string
	GatewayPort        *string
	MessageServiceKey  *string
	NoAuth             *bool
)

func InitConfig() {

	MessageServiceAddr = flag.String("maddr", "localhost:50051", "The server address in the format of host:port")
	GatewayPort = flag.String("port", "8080", "port on which gateway should run")
	MessageServiceKey = flag.String("mkey", "cute", "message service auth key")
	NoAuth = flag.Bool("noauth", false, "disable useer authorization")

	InitEnv()
	flag.Parse()
}

func InitEnv() {
	_ = godotenv.Load()

	if msgServiceAddr, ok := os.LookupEnv("MESSAGE_SERVICE_IP"); ok {
		*MessageServiceAddr = msgServiceAddr
	} else {
		*MessageServiceAddr = "localhost:50051"
	}

	if gatewayPort, ok := os.LookupEnv("GATEWAY_PORT"); ok {
		*GatewayPort = gatewayPort
	} else {
		*GatewayPort = "8080"
	}

	if msgServiceApiKey, ok := os.LookupEnv("MESSAGE_SERVICE_API_KEY"); ok {
		*MessageServiceKey = msgServiceApiKey
	} else {
		*MessageServiceKey = "cute"
	}
	if noauth, ok := os.LookupEnv("NO_AUTH"); ok {
		switch noauth {
		case "true":
			*NoAuth = true
		case "false":
			*NoAuth = false
		}
	} else {
		*MessageServiceKey = "cute"
	}
}
