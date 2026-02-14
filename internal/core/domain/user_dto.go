package domain

import "time"

// CreateUserRequest represents data for creating a new user
type CreateUserRequest struct {
	FirstName    string `json:"first_name" binding:"required,min=1,max=100"`
	LastName     string `json:"last_name" binding:"required,min=1,max=100"`
	Email        string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10"`
	DateOfBirth  string `json:"date_of_birth" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
}

// ParseDateOfBirth parses the date of birth string
func (r *CreateUserRequest) ParseDateOfBirth() (time.Time, error) {
	return time.Parse("2006-01-02", r.DateOfBirth)
}

// UpdateUserRequest represents data for updating user profile
type UpdateUserRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID           int `json:"id"`
	Active       bool   `json:"active"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth  string `json:"date_of_birth"`
	Age          int    `json:"age"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Active:       u.Active,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		FullName:     u.FullName(),
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		DateOfBirth:  u.DateOfBirth.Format("2006-01-02"),
		Age:          u.Age(),
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}
