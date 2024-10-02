package handler

import (
	"bytes"
	"fmt"
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

func setupBalanceHandler() *BalanceHandler {
	repo := &repository.TestBalanceRepository{}
	balanceService := service.NewBalanceService(repo)
	return NewBalanceHandler(balanceService)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Query("user_id")

		if len(uid) > 0 {
			userID, _ := strconv.Atoi(uid)

			c.Set("user", model.User{
				ID:    uint64(userID),
				Login: "registered_user",
			})
		}
	}
}

func TestRegisterWithdraw(t *testing.T) {
	balanceHandler := setupBalanceHandler()

	r := gin.New()

	r.Use(authMiddleware())

	r.POST("/api/user/balance/withdraw", balanceHandler.RegisterWithdraw)

	tests := []struct {
		name     string
		request  string
		userID   int
		method   string
		wantCode int
		body     []byte
	}{
		{
			name:     "successful register withdraw",
			request:  "/api/user/balance/withdraw",
			userID:   111,
			method:   http.MethodPost,
			body:     []byte(`{"order":"41004","sum":100}`),
			wantCode: http.StatusOK,
		},
		{
			name:     "unauthorized",
			request:  "/api/user/balance/withdraw",
			userID:   0,
			method:   http.MethodPost,
			body:     []byte(`{"order":"41004","sum":100}`),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "malformed request",
			request:  "/api/user/balance/withdraw",
			userID:   111,
			method:   http.MethodPost,
			body:     []byte(`"order":"41004","sum":100`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid withdrawal number",
			request:  "/api/user/balance/withdraw",
			userID:   111,
			method:   http.MethodPost,
			body:     []byte(`{"order":"111","sum":100}`),
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "no funds",
			request:  "/api/user/balance/withdraw",
			userID:   222,
			method:   http.MethodPost,
			body:     []byte(`{"order":"41004","sum":100}`),
			wantCode: http.StatusPaymentRequired,
		},
		{
			name:     "unexpected error",
			request:  "/api/user/balance/withdraw",
			userID:   555,
			method:   http.MethodPost,
			body:     []byte(`{"order":"41004","sum":100}`),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.request

			if tt.userID > 0 {
				target += fmt.Sprintf("?user_id=%d", tt.userID)
			}

			request := httptest.NewRequest(tt.method, target, bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.wantCode, result.StatusCode)
		})
	}
}

func TestListUserWithdrawals(t *testing.T) {
	balanceHandler := setupBalanceHandler()

	r := gin.New()
	r.Use(authMiddleware())

	r.GET("/api/user/withdrawals", balanceHandler.ListUserWithdrawals)

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name    string
		request string
		userID  int
		method  string
		want    want
	}{
		{
			name:    "success list withdrawals",
			request: "/api/user/withdrawals",
			userID:  111,
			method:  http.MethodGet,
			want: want{
				code: http.StatusOK,
				body: `[{"order":"78477","sum":200,"processed_at":"2024-01-01T00:00:00Z"}]`,
			},
		},
		{
			name:    "unauthorized",
			request: "/api/user/withdrawals",
			userID:  0,
			method:  http.MethodGet,
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:    "unexpected error",
			request: "/api/user/withdrawals",
			userID:  222,
			method:  http.MethodGet,
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name:    "empty withdrawals",
			request: "/api/user/withdrawals",
			userID:  333,
			method:  http.MethodGet,
			want: want{
				code: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.request

			if tt.userID > 0 {
				target += fmt.Sprintf("?user_id=%d", tt.userID)
			}

			request := httptest.NewRequest(tt.method, target, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)

			if len(tt.want.body) > 0 {
				resultBody, err := io.ReadAll(result.Body)
				assert.NoError(t, err)

				assert.Equal(t, tt.want.body, string(resultBody))
			}
		})
	}
}


func TestGetUserBalance(t *testing.T) {
	balanceHandler := setupBalanceHandler()

	r := gin.New()
	r.Use(authMiddleware())

	r.GET("/api/user/balance", balanceHandler.GetUserBalance)

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name    string
		request string
		userID  int
		method  string
		want    want
	}{
		{
			name:    "success get balance",
			request: "/api/user/balance",
			userID:  111,
			method:  http.MethodGet,
			want: want{
				code: http.StatusOK,
				body: `{"current":200,"withdrawn":200}`,
			},
		},
		{
			name:    "unauthorized",
			request: "/api/user/balance",
			userID:  0,
			method:  http.MethodGet,
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:    "unexpected error",
			request: "/api/user/balance",
			userID:  333,
			method:  http.MethodGet,
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.request

			if tt.userID > 0 {
				target += fmt.Sprintf("?user_id=%d", tt.userID)
			}

			request := httptest.NewRequest(tt.method, target, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)

			if len(tt.want.body) > 0 {
				resultBody, err := io.ReadAll(result.Body)
				assert.NoError(t, err)

				assert.Equal(t, tt.want.body, string(resultBody))
			}
		})
	}
}
