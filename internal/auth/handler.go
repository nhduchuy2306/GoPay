package auth

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService Service
}

func NewHandler(s Service) *Handler {
	return &Handler{authService: s}
}

func (h *Handler) RegisterRouters(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
}

func (h *Handler) Login(c *gin.Context) {

}

func (h *Handler) Register(c *gin.Context) {

}
