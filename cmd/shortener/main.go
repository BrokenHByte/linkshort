package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	configServer = config.GetConfig()
	handlersServer := handlers.NewHandlers(configServer, storageServer)
	// Загрузку из файла, как частный случай источника данных
	handlersServer.LoadStorageFromFile(configServer.PathDataStorage)
	serverInst := server.RunServer(httpServerExitDone, handlersServer, configServer)

	// Ждём сигнала завершения, останавливаем сервер, сохраняем данные
	<-stop
	if err := serverInst.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
	httpServerExitDone.Wait()
	handlersServer.SaveStorageFromFile(configServer.PathDataStorage)
}
