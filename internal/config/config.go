package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddr string
	BaseURL    string
}

func NewServerConfig() *ServerConfig {
	var c ServerConfig
	c.ServerAddr = ""
	c.BaseURL = ""
	return &c
}

func ParseConfig(config *ServerConfig) {
	serverAddr := flag.String("a", ":8080", "address and port to run server")
	shortReplaceAddr := flag.String("b", "http://localhost:8080", "address short link")
	flag.Parse()

	if config.ServerAddr = os.Getenv("SERVER_ADDRESS"); config.ServerAddr == "" {
		config.ServerAddr = *serverAddr
	}

	if config.BaseURL = os.Getenv("BASE_URL"); config.BaseURL == "" {
		config.BaseURL = *shortReplaceAddr
	}
}
