package repository

import (
	"bulls-lab-be/internal/core/domain"
	"context"

	"gorm.io/gorm"
)

type StockHoldingRepository struct {
	db *gorm.DB
}

func NewStockHoldingRepository(db *gorm.DB) *StockHoldingRepository {
	return &StockHoldingRepository{
		db: db,
	}
}

func (r *StockHoldingRepository) CreateHolding(ctx context.Context, holding *domain.StockHolding) error {
	return r.db.WithContext(ctx).Create(holding).Error
}

func (r *StockHoldingRepository) GetHoldingsByUserAndTicker(ctx context.Context, userID int, ticker string) ([]*domain.StockHolding, error) {
	var holdings []*domain.StockHolding
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND stock_ticker = ? AND current_quantity > 0", userID, ticker).
		Order("created_at asc").
		Find(&holdings).Error
	return holdings, err
}

func (r *StockHoldingRepository) UpdateHoldingQuantity(ctx context.Context, holdingID int, newQuantity int) error {
	return r.db.WithContext(ctx).
		Model(&domain.StockHolding{}).
		Where("id = ?", holdingID).
		Update("current_quantity", newQuantity).Error
}

func (r *StockHoldingRepository) CreateSellEntry(ctx context.Context, entry *domain.HoldingSellEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *StockHoldingRepository) GetHoldingsSummaryByUser(ctx context.Context, userID int) ([]*domain.PortfolioHolding, error) {
	var holdings []*domain.PortfolioHolding

	query := `
		SELECT 
			sh.stock_ticker, 
			SUM(sh.current_quantity) as total_quantity,
			SUM(o.price * sh.current_quantity) / SUM(sh.current_quantity) as avg_buy_price
		FROM 
			stock_holdings sh
		JOIN 
			orders o ON sh.order_id = o.id
		WHERE 
			sh.user_id = ? AND sh.current_quantity > 0
		GROUP BY 
			sh.stock_ticker
	`

	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&holdings).Error
	return holdings, err
}
