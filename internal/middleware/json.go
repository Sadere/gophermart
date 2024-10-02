package middleware

import (
	"github.com/gin-gonic/gin"
)

// Устанавливаем заголовок Content-Type для json данных
func JSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		c.Next()
	}
}