package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
	config      config.Config
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(userService *service.UserService, config config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      config,
	}
}

func (u *AuthHandler) Register(c *gin.Context) {
	request := AuthRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to parse request: %v", err),
		})
		return
	}

	// Регистрируем юзера
	newUser, err := u.userService.RegisterUser(request.Login, request.Password)

	// Проверяем существует ли юзер с таким логином
	if errors.Is(err, &service.ErrUserExists{Login: request.Login}) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
	}

	// Обработка остальных ошибок
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Успешная аутентификация нового юзера
	u.authUser(newUser.ID, c)
}

func (u *AuthHandler) Login(c *gin.Context) {
	request := AuthRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to parse request: %v", err),
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
