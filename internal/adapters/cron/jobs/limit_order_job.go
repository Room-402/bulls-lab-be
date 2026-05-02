package jobs

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/ports"
	"context"
	"log"
)

type LimitOrderJob struct {
	orderService  ports.OrderService
	marketService ports.MarketService
}

func NewLimitOrderJob(orderService ports.OrderService, marketService ports.MarketService) *LimitOrderJob {
	return &LimitOrderJob{
		orderService:  orderService,
		marketService: marketService,
	}
}

func (j *LimitOrderJob) Name() string {
	return "LimitOrderExecutor"
}

func (j *LimitOrderJob) Run() {
	ctx := context.Background()
	log.Println("⏳ [LimitOrderJob] Starting execution cycle...")

	// 1. Get all pending (STOP_LOSS) and placed (LIMIT) orders
	orders, err := j.orderService.GetPendingLimitOrders(ctx)
	if err != nil {
		log.Printf("❌ [LimitOrderJob] Error fetching orders: %v", err)
		return
	}

	if len(orders) == 0 {
		log.Println("✅ [LimitOrderJob] No orders to process.")
		return
	}

	// 2. Group by ticker to avoid duplicate API calls
	tickerPrices := make(map[string]float64)

	for _, order := range orders {
		if _, ok := tickerPrices[order.StockTicker]; !ok {
			price, err := j.marketService.GetStockPrice(ctx, order.StockTicker)
			if err != nil {
				log.Printf("⚠️ [LimitOrderJob] Failed to fetch price for %s: %v", order.StockTicker, err)
				continue
			}
			tickerPrices[order.StockTicker] = price
		}

		currentPrice := tickerPrices[order.StockTicker]
		if currentPrice == 0 {
			continue
		}

		// 3. Handle STOP_LOSS trigger (PENDING status)
		if order.OrderStatus == string(constants.OrderStatusPending) && order.OrderCategory == string(constants.StopLossOrderCategory) {
			triggered := false
			if order.OrderType == string(constants.OrderTypeBuy) {
				if currentPrice >= order.TriggerPrice {
					triggered = true
				}
			} else if order.OrderType == string(constants.OrderTypeSell) {
				if currentPrice <= order.TriggerPrice {
					triggered = true
				}
			}

			if triggered {
				log.Printf("🔔 [LimitOrderJob] STOP_LOSS triggered for order %d (%s at %.2f)",
					order.ID, order.StockTicker, currentPrice)

				order.OrderStatus = string(constants.OrderStatusPlaced)
			} else {
				continue
			}
		}

		// 4. Handle Execution (PLACED status)
		if order.OrderStatus == string(constants.OrderStatusPlaced) {
			shouldExecute := false

			// Market orders that were triggered or placed
			if order.ExecutionType == string(constants.ExecutionTypeMarket) {
				shouldExecute = true
			} else if order.ExecutionType == string(constants.ExecutionTypeLimit) {
				// Limit price check
				if order.OrderType == string(constants.OrderTypeBuy) {
					if currentPrice <= order.Price {
						shouldExecute = true
					}
				} else if order.OrderType == string(constants.OrderTypeSell) {
					if currentPrice >= order.Price {
						shouldExecute = true
					}
				}
			}

			if shouldExecute {
				log.Printf("🚀 [LimitOrderJob] Executing order %d for %s (Limit: %.2f, Current: %.2f)",
					order.ID, order.StockTicker, order.Price, currentPrice)

				if err := j.orderService.ExecuteOrder(ctx, order); err != nil {
					log.Printf("❌ [LimitOrderJob] Error executing order %d: %v", order.ID, err)
				} else {
					log.Printf("✅ [LimitOrderJob] Order %d executed successfully.", order.ID)
				}
			}
		}
	}
}
