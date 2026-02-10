package repository

import (
	"errors"
	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type memRepo struct {
	users map[int]*domain.User
}

func NewMemRepo() ports.UserRepository {
	return &memRepo{users: make(map[int]*domain.User)}
}

func (r *memRepo) GetByID(id int) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *memRepo) Save(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}
