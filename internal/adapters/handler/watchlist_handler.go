package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type WatchlistHandler struct {
	service ports.WatchlistService
}

func NewWatchlistHandler(service ports.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{service: service}
}

func getUserID(c *gin.Context) (int, bool) {
	val, exists := c.Get("userId")
	if !exists {
		return 0, false
	}
	id, ok := val.(int)
	return id, ok
}

// POST /api/v1/watchlists
func (h *WatchlistHandler) CreateWatchlist(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	watchlist, err := h.service.CreateWatchlist(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, watchlist)
}

// GET /api/v1/watchlists
func (h *WatchlistHandler) GetWatchlists(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	watchlists, err := h.service.GetWatchlists(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"watchlists": watchlists})
}

// PUT /api/v1/watchlists/:id
func (h *WatchlistHandler) UpdateWatchlist(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	var req domain.UpdateWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	watchlist, err := h.service.UpdateWatchlist(c.Request.Context(), id, userID, &req)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, watchlist)
}

// DELETE /api/v1/watchlists/:id
func (h *WatchlistHandler) DeleteWatchlist(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	if err := h.service.DeleteWatchlist(c.Request.Context(), id, userID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "watchlist deleted"})
}

// POST /api/v1/watchlists/:id/stocks
func (h *WatchlistHandler) AddStock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	watchlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	var req domain.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stock, err := h.service.AddStock(c.Request.Context(), watchlistID, userID, &req)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

// DELETE /api/v1/watchlists/:id/stocks
func (h *WatchlistHandler) RemoveStock(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	watchlistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid watchlist ID"})
		return
	}

	var req domain.RemoveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveStock(c.Request.Context(), watchlistID, userID, &req); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock removed"})
}
