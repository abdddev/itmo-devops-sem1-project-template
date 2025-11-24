package config

import "os"

const MigrationsDirEnv = "MIGRATIONS_DIR"

func MigrationsDir() string {
	if v := os.Getenv(MigrationsDirEnv); v != "" {
		return v
	}
	return "migrations"
}
