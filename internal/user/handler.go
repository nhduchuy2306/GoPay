package user

import (
	"gopay/internal/auth/token"
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
	g.GET("/me", h.Me)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", middleware.AuthMiddleWare(), h.Update)
	g.DELETE("/:id", middleware.AuthMiddleWare(), h.Delete)
}

func (h *Handler) GetAll(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	usr, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, usr)
}

func (h *Handler) Me(c *gin.Context) {
	tokenStr, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	claims, err := token.ParseToken(tokenStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	email, _ := claims["email"].(string)
	usr, err := h.service.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, usr)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var body UpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.service.Update(id, User{Email: body.Email, Password: body.Password, Role: body.Role})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
