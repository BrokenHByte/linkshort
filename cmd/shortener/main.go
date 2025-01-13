package main

import (
	"time"

	"github.com/BrokenHByte/linkshort/internal/config"
	"github.com/BrokenHByte/linkshort/internal/handlers"
	"github.com/BrokenHByte/linkshort/internal/linkstorage"
	"github.com/BrokenHByte/linkshort/internal/server"
	"golang.org/x/exp/rand"
)

var storageServer handlers.ActionsStorage = &linkstorage.LinkStorage{}
var configServer *config.ServerConfig

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	configServer = config.GetConfig()
	handlersServer := handlers.NewHandlers(configServer, storageServer)
	server.RunServer(handlersServer, configServer)
}
