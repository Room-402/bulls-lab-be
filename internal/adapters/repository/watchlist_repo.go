package repository

import (
	"context"
	"sync"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type watchlistRepo struct {
	Watchlists map[int]*domain.Watchlist // Changed from map[string] to map[int]
	nextID     int                       // Added for auto-increment simulation
	mu         sync.RWMutex
}

func NewWatchlistRepo() ports.WatchlistRepository {
	return &watchlistRepo{
		Watchlists: make(map[int]*domain.Watchlist),
		nextID:     1,
	}
}

func (r *watchlistRepo) Create(ctx context.Context, watchlist *domain.Watchlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	watchlist.ID = r.nextID
	r.nextID++
	r.Watchlists[watchlist.ID] = watchlist
	return nil
}

func (r *watchlistRepo) GetByID(ctx context.Context, id int) (*domain.Watchlist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	watchlist, exists := r.Watchlists[id]
	if !exists || !watchlist.Active {
		return nil, nil
	}
	return watchlist, nil
}

func (r *watchlistRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.Watchlist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var watchlists []*domain.Watchlist
	for _, w := range r.Watchlists {
		if w.UserID == userID && w.Active {
			watchlists = append(watchlists, w)
		}
	}
	return watchlists, nil
}

func (r *watchlistRepo) Update(ctx context.Context, watchlist *domain.Watchlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Watchlists[watchlist.ID]; !exists {
		return nil
	}
	r.Watchlists[watchlist.ID] = watchlist
	return nil
}

func (r *watchlistRepo) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if w, exists := r.Watchlists[id]; exists {
		w.Active = false
	}
	return nil
}
