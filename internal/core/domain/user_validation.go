package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrInvalidMobile = errors.New("invalid Indian mobile number format (must be 10 digits starting with 6-9)")
	ErrInvalidAge    = errors.New("user must be 18 years or older")
	ErrInvalidName   = errors.New("name cannot be empty")
	ErrNameTooLong   = errors.New("name is too long (max 100 characters)")
)

// ValidateUserData validates user input data
func ValidateUserData(firstName, lastName, email, mobileNumber string, dob time.Time) error {
	if err := ValidateName(firstName); err != nil {
		return err
	}

	if err := ValidateName(lastName); err != nil {
		return err
	}

	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidateIndianMobile(mobileNumber); err != nil {
		return err
	}

	if err := ValidateAge(dob); err != nil {
		return err
	}

	return nil
}

// ValidateName validates first and last names
func ValidateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}

	if len(name) > 100 {
		return ErrNameTooLong
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)

	if !matched {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateIndianMobile validates Indian mobile number
func ValidateIndianMobile(mobile string) error {
	if mobile == "" {
		return ErrInvalidMobile
	}

	pattern := `^[6-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)

	if !matched {
		return ErrInvalidMobile
	}

	return nil
}

// ValidateAge checks if user is 18 or older
func ValidateAge(dob time.Time) error {
	now := time.Now()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() ||
		(now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	if age < 18 {
		return ErrInvalidAge
	}

	return nil
}
