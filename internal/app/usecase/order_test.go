package usecase

import (
	"context"
	"encoding/json"
	"level0/internal/app/cache"
	"level0/internal/app/config"
	"level0/internal/app/models"
	"strconv"
	"time"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repositoryMock struct {
	counter int
}

func (r *repositoryMock) SaveAll(ctx context.Context, orders []*models.Order) ([]*models.Order, error) {
	savedOrders := make([]*models.Order, 0)
	for _, el := range orders {
		order := &models.Order{}
		order.ID = r.counter
		order.OrderInfoJSON = el.OrderInfoJSON
		r.counter++
		savedOrders = append(savedOrders, order)
	}
	return savedOrders, nil
}

func (r *repositoryMock) FindAll(ctx context.Context) ([]*models.Order, error) {
	return make([]*models.Order, 0), nil
}

func TestProduceOrder(t *testing.T) {
	config := &config.Config{}
	ctx := context.Background()
	repository := &repositoryMock{1}
	cache, err := cache.NewOrderCache(ctx, repository)
	require.NoError(t, err)
	service, err := NewOrderService(config, cache, repository)
	require.NoError(t, err)
	go service.SaveOrders(ctx)
	fileContent, err := os.ReadFile("./testdata/orders_test.json")
	require.NoError(t, err)
	var orders []map[string]interface{}
	err = json.Unmarshal(fileContent, &orders)
	require.NoError(t, err)
	type test struct {
		json string
		id   int
	}
	var i int = 1
	tests := make([]*test, 0)
	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		require.NoError(t, err)
		service.ProduceOrder(string(orderJSON))
		tests = append(tests, &test{string(orderJSON), i})
		i++
	}
	time.Sleep(time.Millisecond * 2000)
	for _, test := range tests {
		t.Run("test#"+strconv.Itoa(test.id), func(t *testing.T) {
			val := service.Get(test.id)
			assert.Equal(t, test.json, val)
		})
	}
}
