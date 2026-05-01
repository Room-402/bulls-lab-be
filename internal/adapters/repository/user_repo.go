package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type memRepo struct {
	users  map[int]*domain.User // Changed from map[string] to map[int]
	nextID int                  // Added for auto-increment simulation
	mu     sync.RWMutex
}

func NewMemRepo() ports.UserRepository {
	return &memRepo{
		users:  make(map[int]*domain.User),
		nextID: 1, // Start IDs from 1
	}
}

func (r *memRepo) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Auto-generate ID (simulating SERIAL)
	user.ID = r.nextID
	r.nextID++

	r.users[user.ID] = user
	return nil
}

func (r *memRepo) GetByID(ctx context.Context, id int) (*domain.User, error) { // Changed from uuid.UUID to int
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
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

func (r *memRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) { // Renamed from GetByPhone
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.PhoneNumber == phone && user.Active { // Changed from PhoneNumber to PhoneNumber
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *memRepo) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists { // No longer need .String()
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *memRepo) SoftDelete(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	user.SoftDelete()
	return nil
}

func (r *memRepo) Restore(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
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

func (r *memRepo) PhoneExists(ctx context.Context, phone string) (bool, error) { // Renamed from PhoneExists
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.PhoneNumber == phone { // Changed from PhoneNumber to PhoneNumber
			return true, nil
		}
	}

	return false, nil
}
