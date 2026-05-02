package domain

type PortfolioHolding struct {
	StockTicker   string  `json:"stock_ticker" db:"stock_ticker"`
	TotalQuantity int     `json:"total_quantity" db:"total_quantity"`
	AvgBuyPrice   float64 `json:"avg_buy_price" db:"avg_buy_price"`
}
