package repository

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) ExecuteOrder(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Model(order).Updates(map[string]interface{}{
		"order_status": string(constants.OrderStatusExecuted),
		"active":       false,
		"updated_at":   time.Now(),
	}).Error
}

func (r *OrderRepository) GetPendingLimitOrders(ctx context.Context) ([]*domain.Order, error) {
	var orders []*domain.Order
	// Fetch both PLACED (to check price match) and PENDING (to check trigger match)
	err := r.db.WithContext(ctx).Where(
		"(order_status = ? OR order_status = ?) AND active = ? AND expires_at > ?",
		string(constants.OrderStatusPlaced),
		string(constants.OrderStatusPending),
		true,
		time.Now(),
	).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) BatchCancelExpiredOrders(ctx context.Context, status string, cancelledStatus string, now time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).
		Where("order_status = ? AND expires_at < ?", status, now).
		Updates(map[string]interface{}{
			"order_status": cancelledStatus,
			"active":       false,
			"updated_at":   time.Now(),
		})
	
	return result.RowsAffected, result.Error
}

func (r *OrderRepository) GetOrdersByTab(ctx context.Context, userID int, tab string) ([]*domain.Order, error) {
	var orders []*domain.Order
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch tab {
	case "open":
		query = query.Where("created_at >= ? AND order_status = ? AND order_category != ?", 
			todayStart, string(constants.OrderStatusPlaced), string(constants.GTTOrderCategory))
	case "history":
		query = query.Where("created_at >= ? AND order_status IN (?)", 
			todayStart, []string{string(constants.OrderStatusExecuted), string(constants.OrderStatusCancelled)})
	case "gtt":
		query = query.Where("order_status = ? AND order_category = ?", 
			string(constants.OrderStatusPlaced), string(constants.GTTOrderCategory))
	case "positions":
		query = query.Where("created_at >= ? AND order_status = ?", 
			todayStart, string(constants.OrderStatusExecuted))
	default:
		// Return all if no tab specified? Or error? Let's return today's orders
		query = query.Where("created_at >= ?", todayStart)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}
