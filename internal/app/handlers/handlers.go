package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderService interface {
	GetOrder(id int) string
}

type orderHandler struct {
	OrderService
}

func NewOrderHandler(s OrderService) *orderHandler {
	return &orderHandler{s}
}

func (h *orderHandler) GetOrderHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	order := h.GetOrder(int(id))
	return c.JSONBlob(http.StatusOK, []byte(order))
}
