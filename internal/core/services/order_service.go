package services

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
	"context"
	"time"
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	now := time.Now()
	var expiryDate time.Time

	if req.OrderCategory == constants.GTT_LOSS_ORDER_CATEGORY {
		expiryDate = time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, now.Location()).AddDate(1, 0, 0)
	} else {
		expiryDate = time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, now.Location())
	}

	order := &domain.Order{
		UserId:        req.UserId,
		StockTicker:   req.StockTicker,
		OrderType:     req.OrderType,
		OrderCategory: req.OrderCategory,
		ProductType:   req.ProductType,
		Quantity:      req.Quantity,
		Price:         req.Price,
		TriggerPrice:  req.TriggerPrice,
		OrderStatus:   constants.ORDER_STATUS_PLACED,
		Active:        true,
		ExecutionType: req.ExecutionType,
		ExpiresAt:     expiryDate,
	}
	
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	if order.ExecutionType == constants.EXECUTION_TYPE_MARKET {
		// Market orders execute immediately
		if err := s.repo.ExecuteOrder(ctx, order); err != nil {
			return nil, err
		}
	}
	
	return order, nil
}

func (s *OrderService) GetPendingLimitOrders(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.GetPendingLimitOrders(ctx)
}

func (s *OrderService) ExecuteOrder(ctx context.Context, order *domain.Order) error {
	return s.repo.ExecuteOrder(ctx, order)
}

func (s *OrderService) CancelExpiredOrders(ctx context.Context) (int64, error) {
	return s.repo.BatchCancelExpiredOrders(ctx, constants.ORDER_STATUS_PLACED, constants.ORDER_STATUS_CANCELLED, time.Now())
}
