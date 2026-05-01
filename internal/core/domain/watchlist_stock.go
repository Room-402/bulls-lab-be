package domain

type WatchlistStock struct {
	ID          int    `json:"id"`
	Active      bool   `json:"active"`
	StockTicker string `json:"stock_ticker"`
	WatchlistID int    `json:"watchlist_id"`
}

type AddStockRequest struct {
	StockTicker string `json:"stock_ticker" binding:"required"`
}

type RemoveStockRequest struct {
	StockTicker string `json:"stock_ticker" binding:"required"`
}
