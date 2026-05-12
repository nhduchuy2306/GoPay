package wallet

import (
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
	g := rg.Group("/wallets")
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}

func (h *Handler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetAll())
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	wallet, ok := h.service.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) Create(c *gin.Context) {
	var body Wallet
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h.service.Create(body))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	wallet, ok := h.service.Delete(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}
