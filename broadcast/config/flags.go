package config

import "flag"

var (
	NoAuth     *bool
	ServerAddr *string
)

func InitFlags() {
	NoAuth = flag.Bool("noAuth", false, "true or false")
	ServerAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
}
