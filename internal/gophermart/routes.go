package gophermart

import (
	"github.com/Sadere/gophermart/internal/handler"
	"github.com/Sadere/gophermart/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (g *GopherMart) SetupRoutes(r *gin.Engine, db *sqlx.DB) {
	userHandler := handler.NewAuthHandler(g.userService, g.config)
	orderHandler := handler.NewOrderHandler(g.orderService)
	balanceHandler := handler.NewBalanceHandler(g.balanceService)

	apiMiddleware := middleware.NewMiddleware(g.userRepo)

	api := r.Group("/api")
	{
		api.POST("/user/register", userHandler.Register)
		api.POST("/user/login", userHandler.Login)
	}

	// Методы, доступные только авторизованным пользователям
	apiAuthRoutes := api.Group("")

	apiAuthRoutes.Use(apiMiddleware.AuthCheck([]byte(g.config.SecretKey)))
	{
		// Orders
		apiAuthRoutes.POST("/user/orders", orderHandler.SaveOrder)
		apiAuthRoutes.GET("/user/orders", middleware.JSON(), orderHandler.ListOrders)

		// Balance
		apiAuthRoutes.POST("/user/balance/withdraw", balanceHandler.RegisterWithdraw)
		apiAuthRoutes.GET("/user/withdrawals", balanceHandler.ListUserWithdrawals)
		apiAuthRoutes.GET("/user/balance", balanceHandler.GetUserBalance)
	}
}
