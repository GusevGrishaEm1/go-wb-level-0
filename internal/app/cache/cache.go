package cache

import (
	"context"
	"level0/internal/app/models"
)

type OrderRepository interface {
	FindAll(ctx context.Context) ([]*models.Order, error)
}

type cacheForOrders struct {
	data map[int]string
}

func NewOrderCache(ctx context.Context, r OrderRepository) (*cacheForOrders, error) {
	cache := &cacheForOrders{
		data: make(map[int]string),
	}
	orders, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, el := range orders {
		cache.Put(el.ID, el.OrderInfoJSON)
	}
	return cache, nil
}

func (c *cacheForOrders) Put(key int, val string) {
	c.data[key] = val
}

func (c *cacheForOrders) Get(key int) string {
	return c.data[key]
}
