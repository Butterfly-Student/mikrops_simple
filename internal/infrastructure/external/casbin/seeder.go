package casbin

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/pkg/config"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// SeedDefaultSuperadmin creates the superadmin user from config if it does not
// exist yet.  If the user already exists, it updates the password to match the
// current config value so that changing config.yaml is always reflected.
func SeedDefaultSuperadmin(db *gorm.DB, cfg *config.DefaultSuperAdminConfig) error {
	hashedPassword, err := utils.HashPassword(cfg.Password)
	if err != nil {
		return err
	}

	var existing entities.AdminUser
	err = db.Where("username = ? AND role = ?", cfg.Username, "superadmin").
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// ── First run: create the superadmin ─────────────────────────────────
		superadmin := entities.AdminUser{
			Username:  cfg.Username,
			Password:  hashedPassword,
			Email:     cfg.Email,
			Role:      "superadmin",
			Status:    "active",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&superadmin).Error; err != nil {
			return err
		}
		logger.Warn("Default superadmin dibuat",
			zap.String("username", cfg.Username),
			zap.String("email", cfg.Email),
			zap.String("note", "Segera ganti password default ini!"))
		return nil
	}

	if err != nil {
		return err
	}

	// ── Already exists: sync password from config ─────────────────────────
	if err := db.Model(&existing).Updates(map[string]interface{}{
		"password":   hashedPassword,
		"status":     "active",
		"is_active":  true,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	logger.Info("Superadmin password disinkronkan dari konfigurasi",
		zap.String("username", cfg.Username),
	)
	return nil
}
