package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/api"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/config"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/cache"
	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/db"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	pg, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()

	r := cache.NewRedis(cfg.RedisURL)
	defer r.Close()

	router := api.NewRouter(cfg, pg, r)

	addr := ":" + cfg.Port
	log.Printf("API server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
