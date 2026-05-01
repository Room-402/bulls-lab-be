package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type postgresWatchlistRepo struct {
	db *gorm.DB
}

func NewPostgresWatchlistRepo(db *gorm.DB) ports.WatchlistRepository {
	return &postgresWatchlistRepo{db: db}
}

func (r *postgresWatchlistRepo) Create(ctx context.Context, watchlist *domain.Watchlist) error {
	return r.db.WithContext(ctx).Create(watchlist).Error
}

func (r *postgresWatchlistRepo) GetByID(ctx context.Context, id int) (*domain.Watchlist, error) {
	var w domain.Watchlist
	err := r.db.WithContext(ctx).Where("id = ? AND active = ?", id, true).First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("watchlist not found")
		}
		return nil, err
	}
	return &w, nil
}

func (r *postgresWatchlistRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.Watchlist, error) {
	var watchlists []*domain.Watchlist
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND active = ?", userID, true).
		Order("id ASC").
		Find(&watchlists).Error
	return watchlists, err
}

func (r *postgresWatchlistRepo) Update(ctx context.Context, watchlist *domain.Watchlist) error {
	return r.db.WithContext(ctx).Save(watchlist).Error
}

func (r *postgresWatchlistRepo) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Watchlist{}).
		Where("id = ?", id).
		Update("active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("watchlist not found")
	}
	return nil
}

// ── Watchlist Stock Repo ──────────────────────────────────────────────────────

type postgresWatchlistStockRepo struct {
	db *gorm.DB
}

func NewPostgresWatchlistStockRepo(db *gorm.DB) ports.WatchlistStockRepository {
	return &postgresWatchlistStockRepo{db: db}
}

func (r *postgresWatchlistStockRepo) Add(ctx context.Context, stock *domain.WatchlistStock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

func (r *postgresWatchlistStockRepo) GetByWatchlistID(ctx context.Context, watchlistID int) ([]*domain.WatchlistStock, error) {
	var stocks []*domain.WatchlistStock
	err := r.db.WithContext(ctx).
		Where("watchlist_id = ?", watchlistID).
		Order("id ASC").
		Find(&stocks).Error
	return stocks, err
}

func (r *postgresWatchlistStockRepo) GetByTickerAndWatchlist(ctx context.Context, watchlistID int, ticker string) (*domain.WatchlistStock, error) {
	var stock domain.WatchlistStock
	err := r.db.WithContext(ctx).
		Where("watchlist_id = ? AND stock_ticker = ?", watchlistID, ticker).
		First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (r *postgresWatchlistStockRepo) UpdateActive(ctx context.Context, id int, active bool) error {
	return r.db.WithContext(ctx).
		Model(&domain.WatchlistStock{}).
		Where("id = ?", id).
		Update("active", active).Error
}
