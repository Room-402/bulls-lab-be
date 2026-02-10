package services

import (
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) CreateUser(user *domain.User) error {
	return s.repo.Save(user)
}
