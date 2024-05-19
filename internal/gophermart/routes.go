package gophermart

import (
	"github.com/Sadere/gophermart/internal/handler"
	"github.com/Sadere/gophermart/internal/middleware"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (g *GopherMart) SetupRoutes(r *gin.Engine, db *sqlx.DB) {

	userRepo := repository.NewPgUserRepository(db)
	userHandler := handler.NewAuthHandler(userRepo, g.config)

	orderRepo := repository.NewPgOrderRepository(db)
	orderHandler := handler.NewOrderHandler(orderRepo, g.config)

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
		// TODO: implement orders, withdrawals
	}
}
