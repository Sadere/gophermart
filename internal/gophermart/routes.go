package gophermart

import (
	"github.com/Sadere/gophermart/internal/handler"
	"github.com/Sadere/gophermart/internal/middleware"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (g *GopherMart) SetupRoutes(r *gin.Engine, db *sqlx.DB) {

	userRepo := repository.NewPgUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewAuthHandler(userService, g.config)

	orderRepo := repository.NewPgOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	apiMiddleware := middleware.NewMiddleware(userRepo)

	api := r.Group("/api")
	{
		api.POST("/user/register", userHandler.Register)
		api.POST("/user/login", userHandler.Login)
	}

	authRoutes := api.Group("")

	authRoutes.Use(apiMiddleware.AuthCheck([]byte(g.config.SecretKey)))
	{
		authRoutes.POST("/user/orders", orderHandler.SaveOrder)
		authRoutes.GET("/user/orders", middleware.JSON(), orderHandler.ListOrders)
	}
}
