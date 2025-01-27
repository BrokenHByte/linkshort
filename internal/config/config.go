package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddr      string
	BaseURL         string
	PathDataStorage string
}

func GetConfig() *ServerConfig {
	var config ServerConfig
	serverAddr := flag.String("a", ":8080", "address and port to run server")
	shortReplaceAddr := flag.String("b", "http://localhost:8080", "address short link")
	pathDataStorage := flag.String("f", "test.dat", "path storage")
	flag.Parse()

	if config.ServerAddr = os.Getenv("SERVER_ADDRESS"); config.ServerAddr == "" {
		config.ServerAddr = *serverAddr
	}

	if config.BaseURL = os.Getenv("BASE_URL"); config.BaseURL == "" {
		config.BaseURL = *shortReplaceAddr
	}

	if config.PathDataStorage = os.Getenv("FILE_STORAGE_PATH"); config.PathDataStorage == "" {
		config.PathDataStorage = *pathDataStorage
	}
	return &config
}
