package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	pricesV0API "project_sem/internal/api/prices/v0"
	"project_sem/internal/config"
	"project_sem/internal/infrastructure/database"
	pricesRepository "project_sem/internal/repository/prices"
	pricesService "project_sem/internal/service/prices"
)

const (
	serverPortEnv     = "APP_PORT"
	defaultServerPort = "8080"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
	}()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	dbCfg := config.LoadDBConfig()
	db, err := database.New(dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("database.New: %w", err)
	}
	defer db.Close()

	repo := pricesRepository.NewRepository(db)
	service := pricesService.NewService(repo)
	api := pricesV0API.NewAPI(service)

	router := newRouter(api)

	port := getEnv(serverPortEnv, defaultServerPort)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("🚀 REST API server listening on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to serve: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("🛑 Shutting down REST API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	log.Println("✅ Server stopped")
	return nil
}

func newRouter(api pricesV0API.PricesAPI) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/v0/prices", api.UploadPrices).Methods(http.MethodPost)
	r.HandleFunc("/api/v0/prices", api.DownloadPrices).Methods(http.MethodGet)
	return r
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
