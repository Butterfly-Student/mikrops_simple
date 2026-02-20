package entities

import (
	"time"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	Username     string    `gorm:"not null" json:"username"`
	Action       string    `gorm:"not null;index" json:"action"`
	ResourceType string    `gorm:"not null" json:"resource_type"`
	ResourceID   string    `gorm:"index" json:"resource_id"`
	Details      string    `gorm:"type:json" json:"details"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Status       string    `gorm:"default:'success';size:20" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
