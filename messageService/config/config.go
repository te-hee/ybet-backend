package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	Env           *string
	CustomBuffer  *int
	NoAuth        *bool
	ServiceApiKey string
	NATSAddress   string
)

func LoadConfig() {
	Env = flag.String("env", "", "")
	CustomBuffer = flag.Int("buffer", 100, "set message buffer size")
	NoAuth = flag.Bool("noauth", false, "disable auth")
	flag.Parse()

	loadEnvs()
}

func loadEnvs() {
	envEnv := os.Getenv("ENV")
	bufferEnv := os.Getenv("BUFFER_SIZE")
	noAuthEnv := os.Getenv("NO_AUTH")
	authKey := os.Getenv("SERVICE_API_KEY")
	natsAddress := os.Getenv("NATS_ADDRESS")

	switch envEnv {
	case "dev":
		*Env = "dev"
	}

	if bufferEnv != "" {
		size, err := strconv.Atoi(bufferEnv)
		if err != nil {
			log.Printf("can't convert %s to number. setting buffer to 100", bufferEnv)
			*CustomBuffer = 100
		} else {
			*CustomBuffer = size
		}
	}

	switch noAuthEnv {
	case "true":
		*NoAuth = true
	case "false":
		*NoAuth = false
		if authKey == "" {
			panic("NO_AUTH is set to false but auth key is not provided")
		}
		ServiceApiKey = authKey
	}

	if natsAddress == "" {
		NATSAddress = "localhost:4222"
	} else {
		NATSAddress = natsAddress
	}
}
