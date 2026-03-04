package usererror

import (
	"fmt"
	"strings"
)

// User-friendly error messages
var (
	// Authentication errors
	ErrNotLoggedIn        = New("You're not logged in", "Please run 'pman login' first to access your vault")
	ErrInvalidPassword    = New("Incorrect master password", "Please try again with the correct password")
	ErrSessionExpired     = New("Your session has expired", "Please login again with 'pman login'")
	ErrInvalidEmail       = New("Invalid email address", "Please enter a valid email like user@example.com")
	ErrEmailInUse         = New("Email already registered", "Try logging in instead, or use a different email")
	ErrInvalidCredentials = New("Invalid email or password", "Please check your credentials and try again")

	// Agent errors
	ErrAgentLocked     = New("Your vault is locked", "Enter your master password to unlock")
	ErrAgentNotRunning = New("Password agent is not running", "The agent should start automatically. Try running 'pman' again")
	ErrAgentConnection = New("Cannot connect to password agent", "Make sure no other instance is running and try again")

	// Encryption errors
	ErrDecryptFailed = New("Could not decrypt password", "This may happen if your master password changed. Try logging in again")
	ErrEncryptFailed = New("Could not encrypt password", "Please try again. If the problem persists, restart the agent")

	// Storage errors
	ErrNoSaltFound   = New("Account setup incomplete", "Please login again to complete setup")
	ErrKeyringAccess = New("Cannot access secure storage", "Make sure your system keyring is unlocked and accessible")

	// Network errors
	ErrServerUnreachable = New("Cannot reach server", "Check your internet connection and try again")
	ErrServerError       = New("Server error occurred", "The server is having issues. Please try again later")
	ErrTimeout           = New("Server is starting up", "The server was sleeping and is waking up. Please try again in 10-30 seconds")

	// Password entry errors
	ErrPasswordNotFound = New("Password not found", "The password entry may have been deleted")
	ErrEmptyFields      = New("All fields are required", "Please fill in all the required fields")
)

// UserError provides user-friendly error messages with hints
type UserError struct {
	Message string
	Hint    string
}

func (e *UserError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s\n  → %s", e.Message, e.Hint)
	}
	return e.Message
}

// New creates a new UserError
func New(message, hint string) *UserError {
	return &UserError{Message: message, Hint: hint}
}

// Wrap wraps an error with user-friendly context
func Wrap(userErr *UserError, cause error) error {
	if cause == nil {
		return userErr
	}
	return &wrappedError{
		UserError: userErr,
		cause:     cause,
	}
}

type wrappedError struct {
	*UserError
	cause error
}

func (e *wrappedError) Error() string {
	return e.UserError.Error()
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

// FromAPIError converts API errors to user-friendly messages
func FromAPIError(statusCode int, message string) *UserError {
	msg := strings.ToLower(message)

	switch {
	case statusCode == 401:
		if strings.Contains(msg, "expired") {
			return ErrSessionExpired
		}
		return ErrInvalidCredentials

	case statusCode == 403:
		return ErrSessionExpired

	case statusCode == 404:
		return ErrPasswordNotFound

	case statusCode == 409:
		if strings.Contains(msg, "email") {
			return ErrEmailInUse
		}
		return New("Conflict", message)

	case statusCode >= 500:
		return ErrServerError

	case strings.Contains(msg, "timeout"):
		return ErrTimeout

	case strings.Contains(msg, "connection"):
		return ErrServerUnreachable

	default:
		// Return the original message if we can't categorize it
		return New(message, "")
	}
}

// FromAgentError converts agent errors to user-friendly messages
func FromAgentError(err error) *UserError {
	if err == nil {
		return nil
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "agent locked"):
		return ErrAgentLocked

	case strings.Contains(msg, "decrypt failed"):
		return ErrDecryptFailed

	case strings.Contains(msg, "encrypt failed"):
		return ErrEncryptFailed

	case strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "no such file"),
		strings.Contains(msg, "the system cannot find"):
		return ErrAgentNotRunning

	case strings.Contains(msg, "bad key"):
		return ErrInvalidPassword

	default:
		return New("Operation failed", err.Error())
	}
}
