package main

import (
	"log"
	"ocean-ranch/config"
	"ocean-ranch/database"
	"ocean-ranch/handlers"
	"ocean-ranch/repository"
	"ocean-ranch/router"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)

	cageRepo := repository.NewCageRepository()
	feedRepo := repository.NewFeedRepository()

	cageHandler := handlers.NewCageHandler(cageRepo)
	feedHandler := handlers.NewFeedHandler(cageRepo, feedRepo)

	r := router.Setup(cageHandler, feedHandler)

	log.Printf("server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
