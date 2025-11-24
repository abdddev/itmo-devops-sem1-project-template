package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"project_sem/internal/config"
	"project_sem/internal/migrator"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load(".env")

	dbCfg := config.LoadDBConfig()

	db, err := sql.Open("postgres", dbCfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}

	log.Println("Connected to database")

	migrationsDir := config.MigrationsDir()
	m := migrator.NewMigrator(db, migrationsDir)

	if err := m.Up(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Migrations applied")
}
