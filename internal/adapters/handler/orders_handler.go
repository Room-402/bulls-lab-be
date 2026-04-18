package handler

import (
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service ports.OrderService
}

func NewOrderHandler(service ports.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})

}
