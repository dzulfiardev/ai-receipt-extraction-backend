package service

import (
	"fmt"

	"github.com/dzulfiardev/receipt-extraction-backend/internal/domain"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/repository"
	"github.com/dzulfiardev/receipt-extraction-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req domain.CreateUserRequest) (*domain.User, error)
	Login(req domain.LoginRequest) (string, *domain.User, error)
	GetUserByID(id int) (*domain.User, error)
}

type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpireHours int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpireHours int) AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

// Register registers a new user
func (s *authService) Register(req domain.CreateUserRequest) (*domain.User, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns JWT token
func (s *authService) Login(req domain.LoginRequest) (string, *domain.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)

	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, s.jwtSecret, s.jwtExpireHours)

	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetUserByID gets user by ID
func (s *authService) GetUserByID(id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
