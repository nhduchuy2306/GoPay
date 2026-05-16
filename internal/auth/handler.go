package auth

import (
	"gopay/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRouters(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
	g.GET("/me", h.Me)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
}

func (h *Handler) Login(context *gin.Context) {
	var body user.LoginRequest
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.Login(body)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, result)
}

func (h *Handler) Register(context *gin.Context) {
	var body user.RegisterRequest
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.service.Register(body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, created)
}

func (h *Handler) Me(context *gin.Context) {
	tokenValue, exists := context.Get("token")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	usr, err := h.service.Me(tokenValue.(string))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, usr)
}

func (h *Handler) Refresh(context *gin.Context) {
	tokenValue, exists := context.Get("token")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	newToken, err := h.service.Refresh(tokenValue.(string))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user.LoginResponse{Token: newToken})
}

func (h *Handler) Logout(context *gin.Context) {
	tokenValue, exists := context.Get("token")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	if err := h.service.Logout(tokenValue.(string)); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
