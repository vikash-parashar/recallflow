package services
package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/recallflow/backend/internal/middleware"
	"github.com/recallflow/backend/internal/models"
	"github.com/recallflow/backend/internal/repositories"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(orgName, email, password, firstName, lastName, phone string) (string, *models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return "", nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create organization and user (this would be a transaction in production)
	orgID := uuid.New().String()
	userID := uuid.New().String()

	user := &models.User{
		ID:             userID,
		OrganizationID: orgID,
		Email:          email,
		PasswordHash:   string(hashedPassword),
		FirstName:      firstName,
		LastName:       lastName,
		Role:           "owner",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := &middleware.UserClaims{
		UserID:         user.ID,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
