package main

import (
	"context"
	"level0/internal/app/config"
	"level0/internal/app/server"
)

func main() {
	context := context.Background()
	config := &config.Config{}
	config.InitByEnv()
	err := server.StartServer(context, config)
	if err != nil {
		panic(err)
	}
}
