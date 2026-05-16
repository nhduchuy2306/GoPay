package transaction

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
	g := rg.Group("/transactions")
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.GET("/wallet/:wallet_id", h.GetByWalletID)
	g.POST("", h.Create)
	g.PATCH("/:id/status", h.UpdateStatus)
}

func (h *Handler) GetAll(c *gin.Context) {
	txs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tx, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *Handler) GetByWalletID(c *gin.Context) {
	walletID := c.Param("wallet_id")
	txs, err := h.service.ListByWalletID(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) Create(c *gin.Context) {
	var body CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := h.service.Create(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var body StatusUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := h.service.UpdateStatus(id, body.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}
