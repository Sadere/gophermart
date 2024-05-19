package handler

import (
	"fmt"

	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
	config    config.Config
}

func NewOrderHandler(orderRepo repository.OrderRepository, config config.Config) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
		config:    config,
	}
}

func (o *OrderHandler) SaveOrder(c *gin.Context) {
	// TODO: implement
	user, _ := c.Get("user")
	fmt.Printf("%+v", user.(model.User))
}
