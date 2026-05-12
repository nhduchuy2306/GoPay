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
}

func (h *Handler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, "123")
}
