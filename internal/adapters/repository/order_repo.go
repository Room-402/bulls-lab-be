package repository

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/domain"
	"context"

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
