package audit

import (
	"context"
	"encoding/json"
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/pkg/config"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type AuditService struct {
	db     *gorm.DB
	config *config.AuditConfig
	mu     sync.Mutex
}

func NewAuditService(db *gorm.DB, cfg *config.AuditConfig) *AuditService {
	return &AuditService{
		db:     db,
		config: cfg,
	}
}

type AuditEntry struct {
	UserID       int64                  `json:"user_id"`
	Username     string                 `json:"username"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Status       string                 `json:"status"`
}

func (s *AuditService) LogAction(ctx context.Context, entry AuditEntry) error {
	if !s.config.Enabled {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	detailsJSON, err := json.Marshal(entry.Details)
	if err != nil {
		logger.Error("Failed to marshal audit details", zap.Error(err))
		detailsJSON = []byte("{}")
	}

	auditLog := entities.AuditLog{
		UserID:       entry.UserID,
		Username:     entry.Username,
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		ResourceID:   entry.ResourceID,
		Details:      string(detailsJSON),
		IPAddress:    entry.IPAddress,
		UserAgent:    entry.UserAgent,
		Status:       entry.Status,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		logger.Error("Failed to create audit log", zap.Error(err))
		return err
	}

	logger.Debug("Audit log created",
		zap.String("action", entry.Action),
		zap.Int64("user_id", entry.UserID))

	return nil
}
