package config

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	NatsAddress   string
	DBAddress     string
	ServerAddress string
	SchemaPath 	  string
	Pool          *pgxpool.Pool
}

func (c *Config) InitDefault() {
	c.ServerAddress = "localhost:8080"
	c.DBAddress = "postgresql://test:test@localhost:5432/test"
	c.NatsAddress = "nats://localhost:4222"
	c.SchemaPath = "schema.json"
}

func (c *Config) InitByEnv() {
	if val := os.Getenv("SERVER_ADDRESS"); val != "" {
		c.ServerAddress = val
	}
	if val := os.Getenv("DB_ADDRESS"); val != "" {
		c.DBAddress = val
	}
	if val := os.Getenv("NATS_ADDRESS"); val != "" {
		c.NatsAddress = val
	}
	if val := os.Getenv("SCHEMA_PATH"); val != "" {
		c.SchemaPath = val
	}
}
