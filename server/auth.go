package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	UsersFile     = "./users.json"
	JWTSecret     = "your-secret-key-change-this" // Should be loaded from env/config
	TokenDuration = 24 * time.Hour
)

// User represents a registered user
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
}

// UserStore manages user data
type UserStore struct {
	mu    sync.RWMutex
	Users map[string]*User // key: username
}

var userStore = &UserStore{
	Users: make(map[string]*User),
}

// Claims for JWT token
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Initialize user store
func initUserStore() error {
	file, err := os.Open(UsersFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create empty users file
			return saveUsers()
		}
		return err
	}
	defer file.Close()

	var users map[string]*User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return err
	}

	userStore.mu.Lock()
	userStore.Users = users
	userStore.mu.Unlock()

	return nil
}

// Save users to disk
func saveUsers() error {
	userStore.mu.RLock()
	defer userStore.mu.RUnlock()

	file, err := os.Create(UsersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(userStore.Users)
}

// Register a new user
func registerUser(username, email, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password required")
	}

	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	// Check if username exists
	if _, exists := userStore.Users[username]; exists {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	for _, u := range userStore.Users {
		if u.Email == email && email != "" {
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate user ID
	userID := generateID()

	user := &User{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}

	userStore.Users[username] = user

	// Save to disk
	if err := saveUsersUnlocked(); err != nil {
		return nil, err
	}

	// Create user's data directory
	userDataDir := getUserDataDir(userID)
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %v", err)
	}

	return user, nil
}

// Authenticate user and return JWT token
func authenticateUser(username, password string) (string, *User, error) {
	userStore.mu.RLock()
	user, exists := userStore.Users[username]
	userStore.mu.RUnlock()

	if !exists {
		return "", nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Update last login
	userStore.mu.Lock()
	user.LastLogin = time.Now()
	userStore.mu.Unlock()
	saveUsers()

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Generate JWT token
func generateToken(user *User) (string, error) {
	expirationTime := time.Now().Add(TokenDuration)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// Verify JWT token and return claims
func verifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Get user by ID
func getUserByID(userID string) (*User, error) {
	userStore.mu.RLock()
	defer userStore.mu.RUnlock()

	for _, user := range userStore.Users {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Get user data directory
func getUserDataDir(userID string) string {
	return fmt.Sprintf("%s/users/%s", config.DataDir, userID)
}

// Get user metadata file
func getUserMetadataFile(userID string) string {
	return fmt.Sprintf("%s/metadata.json", getUserDataDir(userID))
}

// Generate random ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Save users without locking (internal use when already locked)
func saveUsersUnlocked() error {
	file, err := os.Create(UsersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(userStore.Users)
}
