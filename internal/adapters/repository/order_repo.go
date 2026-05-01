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
	return r.db.WithContext(ctx).Model(order).Update("order_status", constants.ORDER_STATUS_EXECUTED).
		Error
}

func (r *OrderRepository) BatchCancelExpiredOrders(ctx context.Context, status string, cancelledStatus string, now time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).
		Where("order_status = ? AND expires_at < ?", status, now).
		Update("order_status", cancelledStatus)
	
	return result.RowsAffected, result.Error
}
