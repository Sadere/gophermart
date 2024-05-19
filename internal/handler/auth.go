package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	config   config.Config
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(userRepo repository.UserRepository, config config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		config:   config,
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

	_, err := u.userRepo.GetUserByLogin(context.Background(), request.Login)

	// Провреяем существует ли пользователь с таким логином
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("user with login '%s' has already registered", request.Login)})
		return
	}

	// Создаем юзера
	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate password hash"})
		return
	}

	newUser := model.User{
		Login:        request.Login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	var newUserID uint64
	newUserID, err = u.userRepo.CreateUser(context.Background(), newUser)

	if err != nil {
		slog.Error("failed to create user: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	u.authUser(newUserID, c)
}

func (u *AuthHandler) Login(c *gin.Context) {
	request := AuthRequest{}
	badCredMsg := "bad credentials"

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to parse request: %v", err),
		})
		return
	}

	user, err := u.userRepo.GetUserByLogin(context.Background(), request.Login)

	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": badCredMsg,
		})
		return
	}

	if !auth.CheckPassword(user.PasswordHash, request.Password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": badCredMsg,
		})
		return
	}

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
