package user

import (
	"gopay/internal/middleware"
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
	g := rg.Group("/users")
	g.GET("", middleware.AuthMiddleWare(), h.GetAll)
	g.GET("/:id", h.GetByID)
}

func (h *Handler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetUsers())
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, ok := h.service.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
