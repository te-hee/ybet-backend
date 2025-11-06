package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	Env          *string
	CustomBuffer *int
	NoAuth       *bool
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

	switch envEnv {
	case "dev":
		*Env = "dev"
	}

	if bufferEnv != "" {
		size, err := strconv.Atoi(bufferEnv)
		if err == nil {
			*CustomBuffer = size
		}
		log.Printf("can't convert %s to number. setting buffer to 100", bufferEnv)
	}

	switch noAuthEnv {
	case "true":
		*NoAuth = true
	case "false":
		*NoAuth = false
	}
}
