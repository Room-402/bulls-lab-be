package constants

type OrderType string
type OrderCategory string
type OrderStatus string
type ExecutionType string

const (
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

const (
	RegularOrderCategory  OrderCategory = "REGULAR"
	StopLossOrderCategory OrderCategory = "STOP_LOSS"
	GTTOrderCategory      OrderCategory = "GTT"
)

const (
	OrderStatusPlaced    OrderStatus = "PLACED"
	OrderStatusExecuted  OrderStatus = "EXECUTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusPending   OrderStatus = "PENDING"
)

const (
	ExecutionTypeMarket ExecutionType = "MARKET"
	ExecutionTypeLimit  ExecutionType = "LIMIT"
)

const (
	PRODUCT_TYPE_DELIVERY  = "DELIVERY"
	PRODUCT_TYPE_INTRADAY  = "INTRADAY"

	MARKET_SERVICE_URL = "http://localhost:8000"
)
