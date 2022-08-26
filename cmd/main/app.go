package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-storage/internal/handler"
	"test-storage/internal/server"
	"test-storage/pkg/config"
	"test-storage/pkg/storage"
	"time"
)

const configPath = "configs/main.json"

func main() {
	// Config init
	var c config.Config
	if err := c.Init(configPath); err != nil {
		log.Fatalf("[CONFIG INIT] | [ERROR]: %s", err.Error())
	}

	// Empty storage init
	store := storage.Storage{}
	store.Init()

	// HTTP Server init
	handler := handler.NewHandler(store)
	app := server.NewServer(&c, handler.Init())

	// Run app
	go func() {
		if err := app.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[SERVER START] || [FAILED]: %s", err.Error())
		}
	}()
	log.Printf("Application started: \n[PORT]: %s\n", c.HTTP.Port)

	// Gracefull shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("[SERVER STOP] || [FAILED]: %s", err.Error())
	}

	log.Println("Application stopped")
}
