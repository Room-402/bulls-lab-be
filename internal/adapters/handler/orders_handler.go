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

	// Override userId from token for security
	userId, exists := c.Get("userId")
	if exists {
		req.UserId = userId.(int)
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

func (h *OrderHandler) GetOrdersByTab(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tab := c.Query("tab")
	if tab == "" {
		tab = "open"
	}

	orders, err := h.service.GetOrdersByTab(c.Request.Context(), userId.(int), tab)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}
