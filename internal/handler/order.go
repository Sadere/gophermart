package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/Sadere/gophermart/internal/service"
	"github.com/Sadere/gophermart/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (o *OrderHandler) SaveOrder(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "empty order number"})
		return
	}

	orderNumber := string(body)

	err = utils.CheckOnlyDigits(orderNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "order number must contain only digits"})
		return
	}

	var isLoaded bool
	isLoaded, err = o.orderService.SaveOrderForUser(currentUser.ID, orderNumber)

	if errors.Is(err, service.ErrOrderInvalidNumber) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, service.ErrOrderExists) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	if isLoaded {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusAccepted)
	}

}

func (o *OrderHandler) ListOrders(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orders, err := o.orderService.GetOrdersByUser(currentUser.ID)
	if errors.Is(err, service.ErrOrdersNotAdded) {
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
