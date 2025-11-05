package config

import (
	"flag"
	"os"
)

var (
	NoAuth     *bool
	ServerAddr *string
)

func InitFlags() {
	NoAuth = flag.Bool("noAuth", false, "true or false")
	ServerAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")

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

	if val := os.Getenv("ADDR"); val != "" {
		*ServerAddr = val
	}
}
