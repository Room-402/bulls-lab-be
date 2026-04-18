package services

import (
	"context"
	"errors"
	"fmt"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type WatchlistService struct {
	repo      ports.WatchlistRepository
	stockRepo ports.WatchlistStockRepository
}

func NewWatchlistService(repo ports.WatchlistRepository, stockRepo ports.WatchlistStockRepository) *WatchlistService {
	return &WatchlistService{
		repo:      repo,
		stockRepo: stockRepo,
	}
}

func (s *WatchlistService) CreateWatchlist(ctx context.Context, userID int, req *domain.CreateWatchlistRequest) (*domain.Watchlist, error) {
	watchlist := &domain.Watchlist{
		Name:   req.Name,
		UserID: userID,
		Active: true,
	}
	if err := s.repo.Create(ctx, watchlist); err != nil {
		return nil, fmt.Errorf("failed to create watchlist: %w", err)
	}
	return watchlist, nil
}

func (s *WatchlistService) GetWatchlists(ctx context.Context, userID int) ([]*domain.WatchlistWithStocks, error) {
	watchlists, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlists: %w", err)
	}

	result := make([]*domain.WatchlistWithStocks, len(watchlists))
	for i, w := range watchlists {
		stocks, err := s.stockRepo.GetByWatchlistID(ctx, w.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stocks for watchlist %d: %w", w.ID, err)
		}
		result[i] = &domain.WatchlistWithStocks{
			ID:     w.ID,
			Active: w.Active,
			Name:   w.Name,
			UserID: w.UserID,
			Stocks: stocks,
		}
	}
	return result, nil
}

func (s *WatchlistService) UpdateWatchlist(ctx context.Context, id int, userID int, req *domain.UpdateWatchlistRequest) (*domain.Watchlist, error) {
	watchlist, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}
	if watchlist.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	watchlist.Name = req.Name
	if err := s.repo.Update(ctx, watchlist); err != nil {
		return nil, fmt.Errorf("failed to update watchlist: %w", err)
	}
	return watchlist, nil
}

func (s *WatchlistService) DeleteWatchlist(ctx context.Context, id int, userID int) error {
	watchlist, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get watchlist: %w", err)
	}
	if watchlist.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete watchlist: %w", err)
	}
	return nil
}

func (s *WatchlistService) AddStock(ctx context.Context, watchlistID int, userID int, req *domain.AddStockRequest) (*domain.WatchlistStock, error) {
	watchlist, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watchlist: %w", err)
	}
	if watchlist.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	existing, err := s.stockRepo.GetByTickerAndWatchlist(ctx, watchlistID, req.StockTicker)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing stock: %w", err)
	}

	if existing != nil {
		if existing.Active {
			return existing, nil
		}
		if err := s.stockRepo.UpdateActive(ctx, existing.ID, true); err != nil {
			return nil, fmt.Errorf("failed to re-activate stock: %w", err)
		}
		existing.Active = true
		return existing, nil
	}

	stock := &domain.WatchlistStock{
		StockTicker: req.StockTicker,
		WatchlistID: watchlistID,
		Active:      true,
	}
	if err := s.stockRepo.Add(ctx, stock); err != nil {
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}
	return stock, nil
}

func (s *WatchlistService) RemoveStock(ctx context.Context, watchlistID int, userID int, req *domain.RemoveStockRequest) error {
	watchlist, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return fmt.Errorf("failed to get watchlist: %w", err)
	}
	if watchlist.UserID != userID {
		return errors.New("unauthorized")
	}

	existing, err := s.stockRepo.GetByTickerAndWatchlist(ctx, watchlistID, req.StockTicker)
	if err != nil {
		return fmt.Errorf("failed to find stock: %w", err)
	}
	if existing == nil || !existing.Active {
		return errors.New("stock not in watchlist")
	}

	if err := s.stockRepo.UpdateActive(ctx, existing.ID, false); err != nil {
		return fmt.Errorf("failed to remove stock: %w", err)
	}
	return nil
}
