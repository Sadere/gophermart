package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sadere/gophermart/internal/service"
	"github.com/Sadere/gophermart/internal/structs"
	"github.com/gin-gonic/gin"
)

type RegisterWithdrawRequest struct {
	Order string  `json:"order" binding:"required"`
	Sum   float64 `json:"sum" binding:"required,gt=0"`
}

type BalanceHandler struct {
	balanceService *service.BalanceService
}

func NewBalanceHandler(balanceService *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

func (h *BalanceHandler) RegisterWithdraw(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	request := RegisterWithdrawRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to parse request: %v", err),
		})
		return
	}

	err = h.balanceService.RegisterWithdraw(currentUser.ID, request.Order, request.Sum)

	// Невалидный номер заказа на вывод
	if errors.Is(err, service.ErrOrderInvalidNumber) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Недостаточно средств для вывода
	if errors.Is(err, service.ErrInsufficientFunds) {
		c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
		return
	}

	// Неизвестная ошибка
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type ListWithdrawalItem struct {
	Order       string          `json:"order"`
	Sum         float64         `json:"sum"`
	ProcessedAt structs.RFCTime `json:"processed_at"`
}

type ListWithdrawalsResponse []ListWithdrawalItem

func (h *BalanceHandler) ListUserWithdrawals(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	withdrawals, err := h.balanceService.ListUserWithdrawals(currentUser.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	response := ListWithdrawalsResponse{}

	for _, withdrawal := range withdrawals {
		response = append(response, ListWithdrawalItem{
			Order:       withdrawal.Number,
			Sum:         withdrawal.Amount,
			ProcessedAt: withdrawal.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BalanceHandler) GetUserBalance(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.balanceService.GetUserBalance(currentUser.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}
