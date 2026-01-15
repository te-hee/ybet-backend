package config

import (
	"flag"
	"os"
)

var (
	NoAuth     *bool
	ServerAddr *string
	NatsAddr   *string
)

func InitFlags() {
	NoAuth = flag.Bool("noAuth", false, "true or false")
	NatsAddr = flag.String("addr", "localhost:4222", "address of NATS")

	flag.Parse()

	InitEnv()
}

func InitEnv() {
	if val := os.Getenv("NO_AUTH"); val != "" {
		switch val {
		case "true":
			*NoAuth = true
		case "false":
			*NoAuth = false

		}
	}

	if val := os.Getenv("NATS_ADDRESS"); val != "" {
		*NatsAddr = val
	}
}
