package validation

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/subrat-dwi/passman-cli/internal/usererror"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if the email is valid
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return usererror.New("Email is required", "Enter your email address")
	}
	if len(email) > 254 {
		return usererror.New("Email is too long", "Maximum 254 characters allowed")
	}
	if !emailRegex.MatchString(email) {
		return usererror.New("Invalid email format", "Use format: user@example.com")
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
		return usererror.New("Password is required", "Enter your master password")
	}
	if len(password) < 8 {
		return usererror.New("Password too short", "Use at least 8 characters")
	}
	if len(password) > 128 {
		return usererror.New("Password too long", "Maximum 128 characters allowed")
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
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character (!@#$...)")
	}

	if len(missing) > 2 {
		return usererror.New("Password too weak", "Add: "+strings.Join(missing, ", "))
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
		return usererror.New("Service name is required", "Enter a name like 'Gmail' or 'Netflix'")
	}
	if len(name) < 1 {
		return usererror.New("Service name too short", "Enter at least 1 character")
	}
	if len(name) > 64 {
		return usererror.New("Service name too long", "Maximum 64 characters allowed")
	}
	return nil
}

// ValidateUsername validates the username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return usererror.New("Username is required", "Enter your username or email for this service")
	}
	if len(username) > 128 {
		return usererror.New("Username too long", "Maximum 128 characters allowed")
	}
	return nil
}

// ValidatePassword validates the password for entries (not master password)
func ValidatePassword(password string) error {
	if password == "" {
		return usererror.New("Password is required", "Enter the password for this service")
	}
	if len(password) > 256 {
		return usererror.New("Password too long", "Maximum 256 characters allowed")
	}
	return nil
}
