package config

import (
	"flag"

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
	flag.StringVar(&c.NatsAddress, "n", "nats://127.0.0.1:4222", "nats address")
}
