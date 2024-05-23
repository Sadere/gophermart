package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {

	r := gin.New()
	r.Use(JSON())
	r.Use(gin.Recovery())

	r.GET("/example", func(c *gin.Context) {})

	request := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, request)

	result := w.Result()

	defer result.Body.Close()

	assert.Contains(t, result.Header.Get("Content-Type"), "application/json")
}
