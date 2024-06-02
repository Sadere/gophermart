package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSaveOrder(t *testing.T) {
	repo := &repository.TestOrderRepository{}
	orderService := service.NewOrderService(repo)
	orderHandler := NewOrderHandler(orderService)

	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		c.Set("user", model.User{
			ID:    111,
			Login: "registered_user",
		})
	}

	r.POST("/api/user/orders", authMiddleware, orderHandler.SaveOrder)
	r.POST("/api/user/orders/unauth", orderHandler.SaveOrder)

	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
		body    []byte
	}{
		{
			name:    "successful save order",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`84913`),
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:    "order already uploaded",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`56317`),
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "empty order number",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(``),
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "order number contains non-numeric",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`test_order`),
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "invalid order number",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`54362`),
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:    "order added by another user",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`24844`),
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name:    "unexpected order error",
			request: "/api/user/orders",
			method:  http.MethodPost,
			body:    []byte(`43513`),
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:    "unauthenticated attempt",
			request: "/api/user/orders/unauth",
			method:  http.MethodPost,
			body:    []byte(`43513`),
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestListOrders(t *testing.T) {
	repo := &repository.TestOrderRepository{}
	orderService := service.NewOrderService(repo)
	orderHandler := NewOrderHandler(orderService)

	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		uid := c.Query("user_id")

		if len(uid) > 0 {
			userID, _ := strconv.Atoi(uid)

			c.Set("user", model.User{
				ID:    uint64(userID),
				Login: "registered_user",
			})
		}
	}

	r.GET("/api/user/orders", authMiddleware, orderHandler.ListOrders)

	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "successful list orders",
			request: "/api/user/orders?user_id=111",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusOK,
				body:       `[{"uploaded_at":"0001-01-01T00:00:00Z","number":"111","status":"NEW"}]`,
			},
		},
		{
			name:    "empty order list",
			request: "/api/user/orders?user_id=555",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name:    "unexpected error",
			request: "/api/user/orders?user_id=222",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name:    "unauthorized request",
			request: "/api/user/orders",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if len(tt.want.body) > 0 {
				resultBody, err := io.ReadAll(result.Body)
				assert.NoError(t, err)

				assert.Equal(t, tt.want.body, string(resultBody))
			}
		})
	}
}
