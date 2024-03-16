package server

import (
	"context"
	"level0/internal/app/cache"
	"level0/internal/app/config"
	"level0/internal/app/handlers"
	"level0/internal/app/repository"
	"level0/internal/app/usecase"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

type OrderHandler interface {
	GetOrderHandler(c echo.Context) error
}

type OrderService interface {
	ProduceOrder(orderJSON string)
}

func StartServer(ctx context.Context, config *config.Config) error {
	pool, err := getPool(ctx, config)
	if err != nil {
		return err
	}
	config.Pool = pool
	err = initTables(ctx, config)
	if err != nil {
		return err
	}

	repository := repository.NewOrderRepository(config)
	cache, err := cache.NewOrderCache(ctx, repository)
	if err != nil {
		return err
	}
	service, err := usecase.NewOrderService(config, cache, repository)
	if err != nil {
		return err
	}
	handler := handlers.NewOrderHandler(service)

	err = initNats(config, service)
	if err != nil {
		return err
	}
	
	go service.SaveOrders(ctx)
	e := echo.New()
	e.GET("/api/order/:id", handler.GetOrderHandler)
	return e.Start(config.ServerAddress)
}

func initNats(config *config.Config, service OrderService) error {
	nc, err := nats.Connect(config.NatsAddress)
	if err != nil {
		return err
	}
	_, err = nc.Subscribe("orders", func(m *nats.Msg) {
		service.ProduceOrder(string(m.Data))
	})
	return err
}
