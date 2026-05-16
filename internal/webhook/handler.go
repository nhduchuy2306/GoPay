package webhook

import (
	"gopay/internal/auth/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRouters(rg *gin.RouterGroup) {
	g := rg.Group("/webhooks")
	g.GET("/endpoints", h.ListEndpoints)
	g.GET("/endpoints/me", h.ListMyEndpoints)
	g.GET("/endpoints/:id", h.GetEndpoint)
	g.POST("/endpoints", h.CreateEndpoint)
	g.PATCH("/endpoints/:id", h.UpdateEndpoint)
	g.DELETE("/endpoints/:id", h.DeleteEndpoint)
	g.GET("/endpoints/:id/deliveries", h.ListDeliveries)
	g.GET("/deliveries/:id", h.GetDelivery)
}

func (h *Handler) ListEndpoints(c *gin.Context) {
	items, err := h.service.ListEndpoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) ListMyEndpoints(c *gin.Context) {
	tokenValue, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	claims, err := token.ParseToken(tokenValue.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID, _ := claims["user_id"].(string)
	items, err := h.service.ListMyEndpoints(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetEndpoint(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetEndpoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) CreateEndpoint(c *gin.Context) {
	tokenValue, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	claims, err := token.ParseToken(tokenValue.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID, _ := claims["user_id"].(string)
	var body CreateEndpointRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.service.CreateEndpoint(userID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) UpdateEndpoint(c *gin.Context) {
	id := c.Param("id")
	var body UpdateEndpointRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.service.UpdateEndpoint(id, body)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteEndpoint(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListDeliveries(c *gin.Context) {
	id := c.Param("id")
	items, err := h.service.ListDeliveries(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetDelivery(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetDelivery(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}
