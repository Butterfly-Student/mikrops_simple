package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
	"github.com/alijayanet/gembok-backend/pkg/utils"
)

func newTestAuthUsecase(t *testing.T) (*AuthUsecase, *mocks.AdminRepository) {
	repo := mocks.NewAdminRepository(t)
	uc := NewAuthUsecase(repo, "test-secret", 24*time.Hour)
	return uc, repo
}

func TestLogin_Success(t *testing.T) {
	uc, repo := newTestAuthUsecase(t)

	hashed, err := utils.HashPassword("password123")
	assert.NoError(t, err)

	admin := &entities.AdminUser{
		ID:       1,
		Username: "admin",
		Password: hashed,
		Role:     "admin",
		Status:   "active",
	}

	repo.On("FindByUsername", "admin").Return(admin, nil)

	resp, err := uc.Login("admin", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, uint(1), resp.UserID)
	assert.Equal(t, "admin", resp.Username)
	assert.Equal(t, "admin", resp.Role)
	assert.Equal(t, int64(86400), resp.ExpiresIn)
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, repo := newTestAuthUsecase(t)

	repo.On("FindByUsername", "nonexistent").Return(nil, errors.New("record not found"))

	resp, err := uc.Login("nonexistent", "password123")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid username or password")
}

func TestLogin_InactiveAccount(t *testing.T) {
	uc, repo := newTestAuthUsecase(t)

	hashed, _ := utils.HashPassword("password123")
	admin := &entities.AdminUser{
		ID:       1,
		Username: "admin",
		Password: hashed,
		Role:     "admin",
		Status:   "inactive",
	}

	repo.On("FindByUsername", "admin").Return(admin, nil)

	resp, err := uc.Login("admin", "password123")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "account is not active")
}

func TestLogin_WrongPassword(t *testing.T) {
	uc, repo := newTestAuthUsecase(t)

	hashed, _ := utils.HashPassword("correct-password")
	admin := &entities.AdminUser{
		ID:       1,
		Username: "admin",
		Password: hashed,
		Role:     "admin",
		Status:   "active",
	}

	repo.On("FindByUsername", "admin").Return(admin, nil)

	resp, err := uc.Login("admin", "wrong-password")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid username or password")
}
