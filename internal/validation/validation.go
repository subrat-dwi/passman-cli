package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if the email is valid
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if len(email) > 254 {
		return fmt.Errorf("email is too long (max 254 characters)")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// PasswordStrength represents the strength of a password
type PasswordStrength int

const (
	PasswordWeak PasswordStrength = iota
	PasswordFair
	PasswordStrong
)

// ValidateMasterPassword validates the master password with detailed feedback
func ValidateMasterPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password is too long (max 128 characters)")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	missing := []string{}
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasDigit {
		missing = append(missing, "digit")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 2 {
		return fmt.Errorf("password needs: %s", strings.Join(missing, ", "))
	}

	return nil
}

// GetPasswordStrength returns the strength of the password
func GetPasswordStrength(password string) (PasswordStrength, string) {
	if len(password) == 0 {
		return PasswordWeak, ""
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	score := 0
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	switch {
	case score >= 5:
		return PasswordStrong, "Strong"
	case score >= 3:
		return PasswordFair, "Fair"
	default:
		return PasswordWeak, "Weak"
	}
}

// ValidateServiceName validates the service name for password entries
func ValidateServiceName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if len(name) < 1 {
		return fmt.Errorf("service name must be at least 1 character")
	}
	if len(name) > 64 {
		return fmt.Errorf("service name is too long (max 64 characters)")
	}
	return nil
}

// ValidateUsername validates the username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > 128 {
		return fmt.Errorf("username is too long (max 128 characters)")
	}
	return nil
}

// ValidatePassword validates the password for entries (not master password)
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) > 256 {
		return fmt.Errorf("password is too long (max 256 characters)")
	}
	return nil
}
