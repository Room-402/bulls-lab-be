package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Active       bool      `json:"active" db:"active"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password"`
	MobileNumber string    `json:"mobile_number" db:"mobile_number"`
	DateOfBirth  time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user instance
func NewUser(firstName, lastName, email, mobileNumber string, dob time.Time, passwordHash string) (*User, error) {
	if err := ValidateUserData(firstName, lastName, email, mobileNumber, dob); err != nil {
		return nil, err
	}

	return &User{
		ID:           uuid.New(),
		Active:       true,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		MobileNumber: mobileNumber,
		DateOfBirth:  dob,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Age calculates the user's age
func (u *User) Age() int {
	now := time.Now()
	age := now.Year() - u.DateOfBirth.Year()

	if now.Month() < u.DateOfBirth.Month() ||
		(now.Month() == u.DateOfBirth.Month() && now.Day() < u.DateOfBirth.Day()) {
		age--
	}

	return age
}

// SoftDelete marks the user as inactive
func (u *User) SoftDelete() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

// Restore reactivates a soft-deleted user
func (u *User) Restore() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

// IsEligibleForTrading checks if user can trade (18+ and active)
func (u *User) IsEligibleForTrading() bool {
	return u.Age() >= 18 && u.Active
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return errors.New("first name and last name cannot be empty")
	}

	if len(firstName) > 100 || len(lastName) > 100 {
		return errors.New("name is too long (max 100 characters)")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()

	return nil
}

// UpdatePassword updates user password hash
func (u *User) UpdatePassword(passwordHash string) {
	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Active
}