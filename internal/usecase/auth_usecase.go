package usecase

import (
	"fmt"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/pkg/utils"
)

type AuthUsecase struct {
	adminRepo repositories.AdminRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthUsecase(adminRepo repositories.AdminRepository, jwtSecret string, jwtExpiry time.Duration) *AuthUsecase {
	return &AuthUsecase{
		adminRepo: adminRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	ExpiresIn int64  `json:"expires_in"`
}

func (u *AuthUsecase) Login(username, password string) (*LoginResponse, error) {
	admin, err := u.adminRepo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if admin.Status != "active" {
		return nil, fmt.Errorf("account is not active")
	}

	if !utils.CheckPassword(password, admin.Password) {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role, u.jwtSecret, u.jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:     token,
		UserID:    admin.ID,
		Username:  admin.Username,
		Role:      admin.Role,
		ExpiresIn: int64(u.jwtExpiry.Seconds()),
	}, nil
}
