package market

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MarketClient struct {
	baseURL string
	client  *http.Client
}

type MarketPriceResponse struct {
	StockDetails struct {
		Price float64 `json:"price"`
	} `json:"stock_details"`
}

func NewMarketClient(baseURL string) *MarketClient {
	return &MarketClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *MarketClient) GetStockPrice(ctx context.Context, ticker string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/stocks/get_stock_details?symbol=%s", c.baseURL, ticker)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Do(req)
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
