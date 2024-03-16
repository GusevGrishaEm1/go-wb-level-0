package usecase

import (
	"context"
	"level0/internal/app/config"
	"level0/internal/app/models"
	"os"
	"path/filepath"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

const MaxSizeArray int = 1000

type OrderRepository interface {
	SaveAll(ctx context.Context, orders []*models.Order) ([]*models.Order, error)
}

type OrderCache interface {
	Put(key int, val string)
	Get(key int) string
}

type orderService struct {
	config       *config.Config
	ch           chan *models.Order
	schemaLoader gojsonschema.JSONLoader
	OrderCache
	OrderRepository
}

func NewOrderService(config *config.Config, storage OrderCache, repository OrderRepository) (*orderService, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	schemaPath := filepath.Join(currentDir, config.SchemaPath)
	schemaLoader := gojsonschema.NewReferenceLoader("file:///"+schemaPath)
	return &orderService{
		config,
		make(chan *models.Order, 1024),
		schemaLoader,
		storage,
		repository,
	}, nil
}

func (s *orderService) GetOrder(id int) string {
	return s.Get(id)
}

func (s *orderService) ProduceOrder(orderJSON string) {
	loader := gojsonschema.NewStringLoader(orderJSON)
	result, err := gojsonschema.Validate(s.schemaLoader, loader)
	if err != nil {
		s.ch <- &models.Order{
			ErrorMesage:   err.Error(),
			OrderInfoJSON: orderJSON,
		}
		return
	}
	if result.Valid() {
		s.ch <- &models.Order{
			OrderInfoJSON: orderJSON,
		}
		return
	}
	msg := getMessage(result)
	s.ch <- &models.Order{
		ErrorMesage:   msg,
		OrderInfoJSON: orderJSON,
	}
}

func getMessage(result *gojsonschema.Result) string {
	var msg string
	for _, desc := range result.Errors() {
		msg += desc.String() + "\n"
	}
	return msg
}

func (s *orderService) SaveOrders(ctx context.Context) {
	orders := make([]*models.Order, 0)
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		if len(orders) > 0 {
			s.SaveAll(ctx, orders)
			orders = orders[:0]
		}
	}()
	for {
		select {
		case val := <-s.ch:
			orders = append(orders, val)
			if len(orders) > MaxSizeArray {
				savedOrders, err := s.SaveAll(ctx, orders)
				if err != nil {
					continue
				}
				s.putAllToCache(savedOrders)
				orders = orders[:0]
			}
		case <-ticker.C:
			if len(orders) > 0 {
				savedOrders, err := s.SaveAll(ctx, orders)
				if err != nil {
					continue
				}
				s.putAllToCache(savedOrders)
				orders = orders[:0]
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *orderService) putAllToCache(orders []*models.Order) {
	for _, el := range orders {
		s.Put(el.ID, el.OrderInfoJSON)
	}
}
