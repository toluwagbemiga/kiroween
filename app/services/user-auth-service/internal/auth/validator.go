package auth

import (
	"fmt"
	"regexp"
	"unicode"
)

var (
	// RFC 5322 compliant email regex (simplified)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidateEmail validates an email address according to RFC 5322
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(email) > 255 {
		return fmt.Errorf("email is too long (max 255 characters)")
	}
	
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidatePassword validates a password for strength
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	if len(password) > 128 {
		return fmt.Errorf("password is too long (max 128 characters)")
	}
	
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	
	return nil
}

// ValidateName validates a user's name
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}
	
	if len(name) > 255 {
		return fmt.Errorf("name is too long (max 255 characters)")
	}
	
	return nil
}
