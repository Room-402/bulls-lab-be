package domain

import "time"

type Order struct {
	ID            int       `json:"id" db:"id"`
	UserId        int       `json:"user_id" db:"user_id"`
	StockTicker   string    `json:"stock_ticker" db:"stock_ticker"`
	OrderType     string    `json:"order_type" db:"order_type"`
	OrderCategory string    `json:"order_category" db:"order_category"`
	ProductType   string    `json:"product_type" db:"product_type"`
	Quantity      int       `json:"quantity" db:"quantity"`
	ExecutionType string    `json:"execution_type" db:"execution_type"`
	Price         float64   `json:"price" db:"price"`
	TriggerPrice  float64   `json:"trigger_price" db:"trigger_price"`
	OrderStatus   string    `json:"order_status" db:"order_status"`
	Active        bool      `json:"active" db:"active"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateOrderRequest struct {
	UserId        int     `json:"user_id"`
	StockTicker   string  `json:"stock_ticker" binding:"required"`
	OrderType     string  `json:"order_type" binding:"required"`
	OrderCategory string  `json:"order_category" binding:"required"`
	ProductType   string  `json:"product_type" binding:"required"`
	ExecutionType string  `json:"execution_type" db:"execution_type"`
	Quantity      int     `json:"quantity" binding:"required,gt=0"`
	Price         float64 `json:"price"`
	TriggerPrice  float64 `json:"trigger_price"`
	StopLossPrice float64 `json:"stop_loss_price"`
}
