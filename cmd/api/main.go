package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jyanza11/qeko-backend/internal/auth"
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

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error while creating postgres connection: ", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Error while migrating database: ", err)
	}

	rdb, err := cache.NewRedis(context.Background(), cfg.RedisAddr)
	if err != nil {
		log.Fatal("Error while connecting to redis: ", err)
	}
	defer rdb.Close()

	tokenManager := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTExpires)
	authRepo := auth.NewRepository(db.DB)
	authService := auth.NewService(authRepo, tokenManager)
	authHandler := auth.NewHandler(authService)

	router := server.NewRouter(
		server.Handlers{
			Health: health.NewHandler(db, rdb),
			Auth:   authHandler,
		},
		server.Middleware{
			Authenticate: auth.Middleware(tokenManager),
		},
	)

	log.Printf("Server is running on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
