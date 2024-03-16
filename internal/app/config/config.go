package config

import (
	"flag"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	NatsAddress   string
	DBAddress     string
	ServerAddress string
	Pool          *pgxpool.Pool
}

func (c *Config) InitByFlags() {
	flag.StringVar(&c.ServerAddress, "s", "localhost:8080", "run address")
	flag.StringVar(&c.DBAddress, "d", "postgresql://test:test@localhost:5432/test", "database URI")
	flag.StringVar(&c.NatsAddress, "n", "nats://localhost:4222", "nats address")
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
}
