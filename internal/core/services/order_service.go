package services

import (
	"bulls-lab-be/constants"
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
	"context"
	"errors"
	"time"
)

type OrderService struct {
	repo          ports.OrderRepository
	marketService ports.MarketService
}

func NewOrderService(repo ports.OrderRepository, marketService ports.MarketService) *OrderService {
	return &OrderService{
		repo:          repo,
		marketService: marketService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	// 1. Accept only DELIVERY for now
	if req.ProductType != constants.PRODUCT_TYPE_DELIVERY {
		return nil, errors.New("only DELIVERY product type is supported currently")
	}

	now := time.Now()
	var expiryDate time.Time

	if req.OrderCategory == constants.GTT_LOSS_ORDER_CATEGORY {
		expiryDate = time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, now.Location()).AddDate(1, 0, 0)
	} else {
		expiryDate = time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, now.Location())
	}

	// 2. Prepare SL Order and Validate BEFORE creating any order in DB
	var slOrder *domain.Order
	if req.OrderCategory == constants.STOP_LOSS_ORDER_CATEGORY {
		// Fetch Current Price for validation
		cmp, err := s.marketService.GetStockPrice(ctx, req.StockTicker)
		if err != nil {
			return nil, errors.New("failed to fetch current market price for validation")
		}

		slOrderType := constants.ORDER_TYPE_SELL
		if req.OrderType == constants.ORDER_TYPE_SELL {
			slOrderType = constants.ORDER_TYPE_BUY
			// Stop Loss Buy must be above Current Price (protecting a short)
			if req.TriggerPrice <= cmp {
				return nil, errors.New("STOP_LOSS trigger price must be greater than current market price for a SELL order")
			}
		} else {
			// Stop Loss Sell must be below Current Price (protecting a long)
			if req.TriggerPrice >= cmp {
				return nil, errors.New("STOP_LOSS trigger price must be less than current market price for a BUY order")
			}
		}

		// Fallback for execution price
		executionPrice := req.StopLossPrice
		if executionPrice == 0 {
			executionPrice = req.Price
		}

		slOrder = &domain.Order{
			UserId:        req.UserId,
			StockTicker:   req.StockTicker,
			OrderType:     slOrderType,
			OrderCategory: constants.STOP_LOSS_ORDER_CATEGORY,
			ProductType:   req.ProductType,
			Quantity:      req.Quantity,
			Price:         executionPrice,
			TriggerPrice:  req.TriggerPrice,
			OrderStatus:   constants.ORDER_STATUS_PENDING,
			Active:        true,
			ExecutionType: req.ExecutionType,
			ExpiresAt:     expiryDate,
		}
	}

	// 3. Create Main Order
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

	// 4. Create Linked SL Order if it was prepared
	if slOrder != nil {
		slOrder.ParentOrderId = &order.ID
		if err := s.repo.Create(ctx, slOrder); err != nil {
			return order, nil // Main order is created, but SL failed
		}
	}

	// 5. Market Execution (for non-SL regular orders)
	if order.ExecutionType == constants.EXECUTION_TYPE_MARKET && order.OrderCategory != constants.STOP_LOSS_ORDER_CATEGORY {
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

func (s *OrderService) GetOrdersByTab(ctx context.Context, userID int, tab string) ([]*domain.OrderWithChild, error) {
	baseOrders, err := s.repo.GetOrdersByTab(ctx, userID, tab)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.OrderWithChild, 0)
	parentMap := make(map[int]*domain.OrderWithChild)

	orderIds := make([]int, 0)
	for _, o := range baseOrders {
		if o.ParentOrderId == nil {
			owc := &domain.OrderWithChild{Order: *o}
			result = append(result, owc)
			parentMap[o.ID] = owc
			orderIds = append(orderIds, o.ID)
		}
	}

	if len(orderIds) > 0 {
		allUserOrders, err := s.repo.GetOrdersByTab(ctx, userID, "all")
		if err == nil {
			for _, o := range allUserOrders {
				if o.ParentOrderId != nil {
					if parent, ok := parentMap[*o.ParentOrderId]; ok {
						parent.ChildOrder = o
					}
				}
			}
		}
	}

	return result, nil
}
