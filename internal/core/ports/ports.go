package ports

import "bulls-lab-be/internal/core/domain"

type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	Save(user *domain.User) error
	GetAll() ([]domain.User, error)
}

type UserService interface {
	GetUser(id int) (*domain.User, error)
	CreateUser(user *domain.User) error
	ListUsers() ([]domain.User, error)
}
