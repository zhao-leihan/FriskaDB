package auth

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
)

// User represents a database user
type User struct {
	Username     string
	PasswordHash string
}

// Authenticator handles user authentication
type Authenticator struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		users: make(map[string]*User),
	}
}

// AddUser adds a new user with hashed password
func (a *Authenticator) AddUser(username, password string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.users[username]; exists {
		return ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.users[username] = &User{
		Username:     username,
		PasswordHash: string(hash),
	}

	return nil
}

// Authenticate verifies username and password
func (a *Authenticator) Authenticate(username, password string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	user, exists := a.users[username]
	if !exists {
		return ErrUserNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

// RemoveUser removes a user
func (a *Authenticator) RemoveUser(username string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.users[username]; !exists {
		return ErrUserNotFound
	}

	delete(a.users, username)
	return nil
}

// UserExists checks if a user exists
func (a *Authenticator) UserExists(username string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, exists := a.users[username]
	return exists
}
