package services

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists
	exists, err = s.repo.PhoneExists(ctx, req.PhoneNumber) // Changed from PhoneExists to PhoneExists, PhoneNumber to PhoneNumber
	if err != nil {
		return nil, fmt.Errorf("failed to check phone existence: %w", err)
	}
	if exists {
		return nil, errors.New("phone number already registered")
	}

	// Parse date of birth
	dob, err := req.ParseDateOfBirth()
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth format: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := domain.NewUser(
		req.FirstName,
		req.LastName,
		req.Email,
		req.PhoneNumber, // Changed from PhoneNumber to PhoneNumber
		dob,
		string(hashedPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, req *domain.LoginUserRequest) (*domain.User, error) {
	var user *domain.User
	var err error

	if req.Email != "" {
		user, err = s.repo.GetByEmail(ctx, req.Email)
	} else if req.PhoneNumber != "" {
		user, err = s.repo.GetByPhone(ctx, req.PhoneNumber)
	} else {
		return nil, errors.New("email or phone number required")
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive() {
		return nil, errors.New("user account is inactive")
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id int) (*domain.User, error) { // Changed from uuid.UUID to int
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, id int, req *domain.UpdateUserRequest) error { // Changed from uuid.UUID to int
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.UpdateProfile(req.FirstName, req.LastName); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// RestoreUser restores a soft-deleted user
func (s *UserService) RestoreUser(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	if err := s.repo.Restore(ctx, id); err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}

	return nil
}

// ListUsers retrieves all active users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// VerifyPassword checks if the provided password matches the user's password
func (s *UserService) VerifyPassword(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
