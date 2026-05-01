package domain

type Watchlist struct {
	ID     int    `json:"id"`
	Active bool   `json:"active"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

type WatchlistWithStocks struct {
	ID     int               `json:"id"`
	Active bool              `json:"active"`
	Name   string            `json:"name"`
	UserID int               `json:"user_id"`
	Stocks []*WatchlistStock `json:"stocks"`
}

type CreateWatchlistRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateWatchlistRequest struct {
	Name string `json:"name" binding:"required"`
}
