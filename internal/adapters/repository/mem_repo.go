package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type memRepo struct {
	users map[string]*domain.User // Use string for UUID key
	mu    sync.RWMutex
}

func NewMemRepo() ports.UserRepository {
	return &memRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *memRepo) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID.String()] = user
	return nil
}

func (r *memRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id.String()]
	if !exists || !user.Active {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *memRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email && user.Active {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *memRepo) GetByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.MobileNumber == mobile && user.Active {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *memRepo) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID.String()]; !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID.String()] = user
	return nil
}

func (r *memRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id.String()]
	if !exists {
		return errors.New("user not found")
	}

	user.SoftDelete()
	return nil
}

func (r *memRepo) Restore(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id.String()]
	if !exists {
		return errors.New("user not found")
	}

	user.Restore()
	return nil
}

func (r *memRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	activeUsers := []*domain.User{}
	for _, user := range r.users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}

	// Simple pagination
	start := offset
	end := offset + limit

	if start > len(activeUsers) {
		return []*domain.User{}, nil
	}

	if end > len(activeUsers) {
		end = len(activeUsers)
	}

	return activeUsers[start:end], nil
}

func (r *memRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

func (r *memRepo) MobileExists(ctx context.Context, mobile string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.MobileNumber == mobile {
			return true, nil
		}
	}

	return false, nil
}