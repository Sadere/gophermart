package handler

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestUserRepository struct {
	registeredUserPwHash string
}

func (tu *TestUserRepository) Create(ctx context.Context, user model.User) (uint64, error) {
	if user.Login == "invalid" {
		return 0, errors.New("test error")
	}

	return 1000, nil
}

func (tu *TestUserRepository) GetUserByID(ctx context.Context, ID uint64) (model.User, error) {
	var user model.User

	if ID == 0 {
		return user, sql.ErrNoRows
	}

	return model.User{
		ID:    111,
		Login: "test_user",
	}, nil
}

func (tu *TestUserRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	if login == "registered_user" {
		return model.User{
			ID:           111,
			Login:        "registered_user",
			PasswordHash: tu.registeredUserPwHash,
		}, nil
	}

	return user, sql.ErrNoRows
}

func TestAuthHandlers(t *testing.T) {
	testPassword := "test_password_123"
	registeredUserPwHash, err := auth.HashPassword(testPassword)

	assert.NoError(t, err, "Failed to generate test password")

	repo := &TestUserRepository{
		registeredUserPwHash: registeredUserPwHash,
	}
	service := service.NewUserService(repo)
	authHandler := NewAuthHandler(service, config.Config{})

	r := gin.New()
	r.POST("/api/user/register", authHandler.Register)
	r.POST("/api/user/login", authHandler.Login)

	type want struct {
		statusCode          int
		authorizationHeader bool
	}
	tests := []struct {
		name    string
		request string
		method  string
		want    want
		body    []byte
	}{
		{
			name:    "successful register",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"login":"test_user1","password":"pw_test"}`),
			want: want{
				authorizationHeader: true,
				statusCode:          http.StatusOK,
			},
		},
		{
			name:    "register without login",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"password":"pw_test"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusBadRequest,
			},
		},
		{
			name:    "register without password",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"login":"test_user1"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusBadRequest,
			},
		},
		{
			name:    "register duplicate user",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"login":"registered_user","password":"pw_test"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusConflict,
			},
		},
		{
			name:    "register with very long password",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"login":"test_user1","password":"` + strings.Repeat("AA", 100) + `"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusInternalServerError,
			},
		},
		{
			name:    "register erroneous user",
			request: "/api/user/register",
			method:  http.MethodPost,
			body:    []byte(`{"login":"invalid","password":"pw_test"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusInternalServerError,
			},
		},
		{
			name:    "successful login",
			request: "/api/user/login",
			method:  http.MethodPost,
			body:    []byte(`{"login":"registered_user","password":"` + testPassword + `"}`),
			want: want{
				authorizationHeader: true,
				statusCode:          http.StatusOK,
			},
		},
		{
			name:    "login without login",
			request: "/api/user/login",
			method:  http.MethodPost,
			body:    []byte(`{"password":"test_pw"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusBadRequest,
			},
		},
		{
			name:    "login unknown user",
			request: "/api/user/login",
			method:  http.MethodPost,
			body:    []byte(`{"login":"test_user333","password":"test_pw"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusUnauthorized,
			},
		},
		{
			name:    "login with wrong password",
			request: "/api/user/login",
			method:  http.MethodPost,
			body:    []byte(`{"login":"registered_user","password":"test_pw"}`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusUnauthorized,
			},
		},
		{
			name:    "malformed json",
			request: "/api/user/login",
			method:  http.MethodPost,
			body:    []byte(`"login":"registered_user","password":"test_pw"`),
			want: want{
				authorizationHeader: false,
				statusCode:          http.StatusBadRequest,
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

			if tt.want.authorizationHeader {
				assert.NotEmpty(t, result.Header.Get("Authorization"))
			} else {
				assert.Empty(t, result.Header.Get("Authorization"))
			}
		})
	}
}
