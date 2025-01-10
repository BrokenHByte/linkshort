package flags

import (
	"flag"

	"github.com/BrokenHByte/linkshort/internal/config"
)

func ParseFlags(config *config.ServerConfig) {
	flag.StringVar(&config.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.PrefixLink, "b", "http://localhost:8080", "prefix address")
	flag.Parse()
}
