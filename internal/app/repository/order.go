package repository

import (
	"context"
	"level0/internal/app/config"
	"level0/internal/app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(config *config.Config) *orderRepository {
	return &orderRepository{config.Pool}
}

func (r *orderRepository) SaveAll(ctx context.Context, orders []*models.Order) ([]*models.Order, error) {
	queryInsWithErr := `
		insert into "order" (order_info_json, error_message) values($1, $2) returning 0 as id, '' as order_info_json
	`
	queryInsWithoutErr := `
		insert into "order" (order_info_json) values($1) returning id as id, cast(order_info_json as varchar) as order_info_json
	`
	batch := &pgx.Batch{}
	for _, el := range orders {
		if el.ErrorMesage == "" {
			batch.Queue(queryInsWithoutErr, el.OrderInfoJSON)
		} else {
			batch.Queue(queryInsWithErr, el.OrderInfoJSON, el.ErrorMesage)
		}
	}
	res := r.pool.SendBatch(ctx, batch)
	savedValidOrders := make([]*models.Order, 0)
	var i = 0
	for i < batch.Len() {
		row := res.QueryRow()
		order := &models.Order{}
		err := row.Scan(&order.ID, &order.OrderInfoJSON)
		if err != nil {
			return nil, err
		}
		if order.ID != 0 && order.OrderInfoJSON != "" {
			savedValidOrders = append(savedValidOrders, order)
		}
		i++
	}
	err := res.Close()
	if err != nil {
		return nil, err
	}
	return savedValidOrders, nil
}

func (r *orderRepository) FindAll(ctx context.Context) ([]*models.Order, error) {
	query := `
		select id, cast(order_info_json as varchar) as order_info_json from "order" where error_message is null
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	orders := make([]*models.Order, 0)
	for rows.Next() {
		order := &models.Order{}
		err = rows.Scan(&order.ID, &order.OrderInfoJSON)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
