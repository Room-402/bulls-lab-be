package ports

import (
	"context"

	"bulls-lab-be/internal/core/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error) // Changed from uuid.UUID to int
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error) // Renamed from GetByPhone
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, id int) error // Changed from uuid.UUID to int
	Restore(ctx context.Context, id int) error    // Changed from uuid.UUID to int
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	PhoneExists(ctx context.Context, phone string) (bool, error) // Renamed from PhoneExists
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
}

type OrderService interface {
	CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error)
}
