package jobs

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LimitOrderJob struct {
	orderService ports.OrderService
}

type MarketPriceResponse struct {
	StockDetails struct {
		Price float64 `json:"price"`
	} `json:"stock_details"`
}

func NewLimitOrderJob(orderService ports.OrderService) *LimitOrderJob {
	return &LimitOrderJob{
		orderService: orderService,
	}
}

func (j *LimitOrderJob) Name() string {
	return "LimitOrderExecutor"
}

func (j *LimitOrderJob) Run() {
	ctx := context.Background()
	log.Println("⏳ [LimitOrderJob] Starting limit order execution cycle...")

	// 1. Get all pending limit orders
	orders, err := j.orderService.GetPendingLimitOrders(ctx)
	if err != nil {
		log.Printf("❌ [LimitOrderJob] Error fetching pending orders: %v", err)
		return
	}

	if len(orders) == 0 {
		log.Println("✅ [LimitOrderJob] No pending limit orders to process.")
		return
	}

	// 2. Group by ticker to avoid duplicate API calls
	tickerPrices := make(map[string]float64)

	for _, order := range orders {
		if _, ok := tickerPrices[order.StockTicker]; !ok {
			price, err := j.fetchCurrentPrice(order.StockTicker)
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

		// 3. Check execution logic
		shouldExecute := false
		if order.OrderType == constants.ORDER_TYPE_BUY {
			if currentPrice <= order.Price {
				shouldExecute = true
			}
		} else if order.OrderType == constants.ORDER_TYPE_SELL {
			if currentPrice >= order.Price {
				shouldExecute = true
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

func (j *LimitOrderJob) fetchCurrentPrice(ticker string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/stocks/get_stock_details?symbol=%s", constants.MARKET_SERVICE_URL, ticker)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("market service returned status %d", resp.StatusCode)
	}

	var data MarketPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.StockDetails.Price, nil
}
