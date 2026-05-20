package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key"
	middleware := AuthMiddleware(secret)

	// Create a valid token
	claims := &UserClaims{
		UserID:         "user-123",
		OrganizationID: "org-123",
		Role:           "owner",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte(secret))

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that will be wrapped
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with middleware
			handler := middleware(testHandler)

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	// Test without user in context
	_, ok := GetUserFromContext(req)
	if ok {
		t.Error("Expected ok=false when no user in context")
	}

	// Test with user in context
	claims := &UserClaims{
		UserID:         "user-123",
		OrganizationID: "org-123",
		Role:           "owner",
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, UserContextKey, claims)
	req = req.WithContext(ctx)

	retrievedClaims, ok := GetUserFromContext(req)
	if !ok {
		t.Error("Expected ok=true when user in context")
	}

	if retrievedClaims.UserID != claims.UserID {
		t.Errorf("Expected UserID %s but got %s", claims.UserID, retrievedClaims.UserID)
	}
}
