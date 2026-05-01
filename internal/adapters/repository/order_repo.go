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
		"order_status": constants.ORDER_STATUS_EXECUTED,
		"active":       false,
		"updated_at":   time.Now(),
	}).Error
}

func (r *OrderRepository) GetPendingLimitOrders(ctx context.Context) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).Where(
		"execution_type = ? AND order_status = ? AND active = ? AND expires_at > ?",
		constants.EXECUTION_TYPE_LIMIT,
		constants.ORDER_STATUS_PLACED,
		true,
		time.Now(),
	).Find(&orders).Error
	return orders, err
}
