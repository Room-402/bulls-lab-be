package domain

import "time"

type StockHolding struct {
	ID              int       `json:"id" db:"id"`
	UserId          int       `json:"user_id" db:"user_id"`
	StockTicker     string    `json:"stock_ticker" db:"stock_ticker"`
	CurrentQuantity int       `json:"current_quantity" db:"current_quantity"`
	OrderId         int       `json:"order_id" db:"order_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type HoldingSellEntry struct {
	ID             int       `json:"id" db:"id"`
	OrderId        int       `json:"order_id" db:"order_id"`
	StockHoldingId int       `json:"stock_holding_id" db:"stock_holding_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
