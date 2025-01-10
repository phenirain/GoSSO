package auth

import (
	"context"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/domain/models"
	userRepo "github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/repository/user"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/services/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestAuth_Login(t *testing.T) {
	mockRepo := new(mocks.UsrRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	authService := New(time.Minute*15, log, mockRepo, []byte("secret"))

	const validPassword = "password"
	validPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(email string)
		expectError bool
	}{
		{
			name:     "valid login",
			email:    "...",
			password: validPassword,
			mockSetup: func(email string) {
				mockRepo.On("GetUser", mock.Anything, email).Return(&models.User{
					Id:           1,
					Email:        "...",
					PasswordHash: validPasswordHash,
				}, nil)
			},
			expectError: false,
		},
		{
			name:     "user not found",
			email:    "...",
			password: "...",
			mockSetup: func(email string) {
				mockRepo.On("GetUser", mock.Anything, email).Return(nil, userRepo.ErrUserNotFound)
			},
			expectError: true,
		},
		{
			name:     "invalid password hash",
			email:    "...",
			password: "...",
			mockSetup: func(email string) {
				mockRepo.On("GetUser", mock.Anything, email).Return(&models.User{
					Id:           1,
					Email:        "...",
					PasswordHash: validPasswordHash,
				}, nil)
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			test.mockSetup(test.email)
			accessToken, refreshToken, err := authService.Login(context.Background(), test.email, test.password)

			if test.expectError {
				assert.Error(t, err, "Expected error but got none")
				assert.Empty(t, accessToken, "Expected no access token")
				assert.Empty(t, refreshToken, "Expected no refresh token")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.NotEmpty(t, accessToken, "Expected access token but got none")
				log.Info("token", "access token", accessToken)
				assert.NotEmpty(t, refreshToken, "Expected refresh token but got none")
				log.Info("token", "refresh token", refreshToken)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuth_Register(t *testing.T) {
	mockRepo := new(mocks.UsrRepo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	authService := New(time.Minute*15, log, mockRepo, []byte("secret"))

	const validPassword = "password"
	validPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(email string)
		expectError bool
	}{
		{
			name:     "user already exist",
			email:    "...",
			password: validPassword,
			mockSetup: func(email string) {
				mockRepo.On("GetUser", mock.Anything, email).Return(&models.User{
					Id:           1,
					Email:        "...",
					PasswordHash: validPasswordHash,
				}, nil)
			},
			expectError: true,
		},
		{
			name:     "valid register",
			email:    "...",
			password: validPassword,
			mockSetup: func(email string) {
				mockRepo.On("GetUser", mock.Anything, email).Return(nil, userRepo.ErrUserNotFound)
				mockRepo.On("AddUser", mock.Anything, email, mock.Anything).Return(int64(1), nil)
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			test.mockSetup(test.email)
			uid, err := authService.Register(context.Background(), test.email, test.password)

			if test.expectError {
				assert.Error(t, err, "Expected error but got none")
				assert.Equal(t, int64(-1), uid, "Expected access token but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.NotEqual(t, uid, -1, "Expected access token but got none")
				log.Info("token", "access token", uid)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TODO: test for refresh

// TODO: test for validate
