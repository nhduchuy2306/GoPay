package wallet

import (
	"gopay/internal/auth/token"
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
	g.GET("/me", h.Me)
	g.GET("/user/:user_id", h.GetByUserID)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PATCH("/:id", h.Update)
	g.POST("/:id/deposit", h.Deposit)
	g.POST("/:id/withdraw", h.Withdraw)
	g.POST("/transfer", h.Transfer)
	g.GET("/:id/history", h.History)
	g.DELETE("/:id", h.Delete)
}

func (h *Handler) GetAll(c *gin.Context) {
	wallets, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallets)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	wallet, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	wallet, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) Me(c *gin.Context) {
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
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in token"})
		return
	}
	wallet, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) Create(c *gin.Context) {
	var body CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet, err := h.service.Create(Wallet{UserID: body.UserID, Currency: body.Currency, Balance: body.Balance})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wallet)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var body CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.service.Update(id, Wallet{UserID: body.UserID, Currency: body.Currency, Balance: body.Balance})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) Deposit(c *gin.Context) {
	id := c.Param("id")
	var body DepositRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet, err := h.service.Deposit(id, body.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) Withdraw(c *gin.Context) {
	id := c.Param("id")
	var body WithdrawRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wallet, err := h.service.Withdraw(id, body.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) Transfer(c *gin.Context) {
	var body TransferRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.Transfer(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) History(c *gin.Context) {
	id := c.Param("id")
	history, err := h.service.History(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
