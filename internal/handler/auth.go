package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-json-experiment/json"
)

type AuthHandler struct {
	userService *service.UserService
	config      config.Config
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewAuthHandler(userService *service.UserService, config config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      config,
	}
}

func authRequest(c *gin.Context) (*AuthRequest, error) {
	request := &AuthRequest{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	err = json.Unmarshal(
		body,
		request,
	
		json.RejectUnknownMembers(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	if len(request.Login) == 0 || len(request.Password) == 0 {
		return nil, errors.New("login and password can't be empty")
	}

	return request, nil
}

func (u *AuthHandler) Register(c *gin.Context) {
	request, err := authRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Регистрируем юзера
	newUser, err := u.userService.RegisterUser(request.Login, request.Password)

	// Проверяем существует ли юзер с таким логином
	if errors.Is(err, &service.ErrUserExists{Login: request.Login}) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Обработка остальных ошибок
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Успешная аутентификация нового юзера
	u.authUser(newUser.ID, c)
}

func (u *AuthHandler) Login(c *gin.Context) {
	request, err := authRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Пытаемся залогиниться
	user, err := u.userService.LoginUser(request.Login, request.Password)

	// Неверные данные для авторизации
	if errors.Is(err, service.ErrBadCredentials) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Остальные ошибки
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Успешная аутентификация
	u.authUser(user.ID, c)
}

func (u *AuthHandler) authUser(userID uint64, c *gin.Context) {
	// Возвращаем токен авторизации
	token, err := auth.CreateToken(userID, time.Now().Add(time.Hour*24), []byte(u.config.SecretKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
}

func getCurrentUser(c *gin.Context) (model.User, error) {
	var currentUser model.User

	errCurrentUser := errors.New("failed to retrieve current user")

	user, ok := c.Get("user")

	if !ok {
		return currentUser, errCurrentUser
	}

	currentUser, ok = user.(model.User)

	if !ok {
		return currentUser, errCurrentUser
	}

	return currentUser, nil
}
