package ports

import (
	"context"
	"time"

	"bulls-lab-be/internal/core/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	PhoneExists(ctx context.Context, phone string) (bool, error)
}

// UserService defines business logic operations
type UserService interface {
	Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	UpdateProfile(ctx context.Context, id int, req *domain.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, req *domain.Order) error
	ExecuteOrder(ctx context.Context, req *domain.Order) error
	GetPendingLimitOrders(ctx context.Context) ([]*domain.Order, error)
	BatchCancelExpiredOrders(ctx context.Context, status string, cancelledStatus string, now time.Time) (int64, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error)
	GetPendingLimitOrders(ctx context.Context) ([]*domain.Order, error)
	ExecuteOrder(ctx context.Context, order *domain.Order) error
	CancelExpiredOrders(ctx context.Context) (int64, error)
}

type WatchlistRepository interface {
	Create(ctx context.Context, watchlist *domain.Watchlist) error
	GetByID(ctx context.Context, id int) (*domain.Watchlist, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Watchlist, error)
	Update(ctx context.Context, watchlist *domain.Watchlist) error
	Delete(ctx context.Context, id int) error
}

type WatchlistStockRepository interface {
	Add(ctx context.Context, stock *domain.WatchlistStock) error
	GetByWatchlistID(ctx context.Context, watchlistID int) ([]*domain.WatchlistStock, error)
	GetByTickerAndWatchlist(ctx context.Context, watchlistID int, ticker string) (*domain.WatchlistStock, error)
	UpdateActive(ctx context.Context, id int, active bool) error
}

type WatchlistService interface {
	CreateWatchlist(ctx context.Context, userID int, req *domain.CreateWatchlistRequest) (*domain.Watchlist, error)
	GetWatchlists(ctx context.Context, userID int) ([]*domain.WatchlistWithStocks, error)
	UpdateWatchlist(ctx context.Context, id int, userID int, req *domain.UpdateWatchlistRequest) (*domain.Watchlist, error)
	DeleteWatchlist(ctx context.Context, id int, userID int) error
	AddStock(ctx context.Context, watchlistID int, userID int, req *domain.AddStockRequest) (*domain.WatchlistStock, error)
	RemoveStock(ctx context.Context, watchlistID int, userID int, req *domain.RemoveStockRequest) error
}
