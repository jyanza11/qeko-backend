package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jyanza11/qeko-backend/internal/health"
	"github.com/jyanza11/qeko-backend/internal/platform/cache"
	"github.com/jyanza11/qeko-backend/internal/platform/config"
	"github.com/jyanza11/qeko-backend/internal/platform/database"
	"github.com/jyanza11/qeko-backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error while loading config: ", err)
	}

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Error while creating postgres pool: ", err)
	}
	defer db.Close()

	rdb, err := cache.NewRedis(ctx, cfg.RedisAddr)

	if err != nil {
		log.Fatal("Error while connecting to redis: ", err)
	}
	defer rdb.Close()

	router := server.NewRouter(server.Handlers{
		Health: health.NewHandler(db, rdb),
	})

	log.Printf("Server is running on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
