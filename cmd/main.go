package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zhaisan-medods/api"
	"zhaisan-medods/db"
	"zhaisan-medods/utils"
)

func main() {
	utils.Logger.Info("Running server...")
	config, err := utils.LoadConfig(".")
	if err != nil {
		utils.Logger.WithError(err).Fatal("Cannot load configurations")
	}

	client := db.MongoDBClient(config.DBDriver)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			utils.Logger.WithError(err).Fatal("Cannot disconnect from db")
		}
	}()

	server := api.NewServer(client, &config)

	go func() {
		if err := server.Start(); err != nil {
			utils.Logger.WithError(err).Fatal("Cannot start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.WithError(err).Fatal("Server forced to shutdown")
	}

	utils.Logger.Info("Server exiting")
}