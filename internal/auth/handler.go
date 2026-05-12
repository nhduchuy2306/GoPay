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
}

func (h *Handler) Login(context *gin.Context) {
	var body user.LoginRequest
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, h.service.Login(body))
}

func (h *Handler) Register(context *gin.Context) {
	var body user.RegisterRequest
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created := h.service.Register(user.User{
		Email:    body.Email,
		Password: body.Password,
	})
	context.JSON(http.StatusCreated, created)
}
