package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bulls-lab-be/internal/core/ports"
)

type PortfolioHandler struct {
	service ports.PortfolioService
}

func NewPortfolioHandler(service ports.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// GET /api/v1/portfolio/holdings
func (h *PortfolioHandler) GetHoldings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	holdings, err := h.service.GetPortfolioHoldings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"holdings": holdings})
}
