package services

import (
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
	"context"
)

type PortfolioService struct {
	holdingRepo ports.StockHoldingRepository
}

func NewPortfolioService(holdingRepo ports.StockHoldingRepository) *PortfolioService {
	return &PortfolioService{
		holdingRepo: holdingRepo,
	}
}

func (s *PortfolioService) GetPortfolioHoldings(ctx context.Context, userID int) ([]*domain.PortfolioHolding, error) {
	return s.holdingRepo.GetHoldingsSummaryByUser(ctx, userID)
}
