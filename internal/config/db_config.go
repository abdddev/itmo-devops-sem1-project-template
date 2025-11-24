package config

import (
	"fmt"
	"os"
)

const (
	databaseHostEnv     = "APP_DB_HOST"
	databaseNameEnv     = "APP_DB_NAME"
	databasePasswordEnv = "APP_DB_PASSWORD"
	databasePortEnv     = "APP_DB_PORT"
	databaseUserEnv     = "APP_DB_USER"

	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "validator"
	defaultDBPassword = "val1dat0r"
	defaultDBName     = "project-sem-1"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv(databaseHostEnv, defaultDBHost),
		Port:     getEnv(databasePortEnv, defaultDBPort),
		User:     getEnv(databaseUserEnv, defaultDBUser),
		Password: getEnv(databasePasswordEnv, defaultDBPassword),
		Name:     getEnv(databaseNameEnv, defaultDBName),
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
