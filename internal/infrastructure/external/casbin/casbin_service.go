package casbin

import (
	"context"
	"fmt"
	"github.com/alijayanet/gembok-backend/internal/domain/events"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/audit"
	eventbus "github.com/alijayanet/gembok-backend/internal/infrastructure/external/eventbus"
	"github.com/alijayanet/gembok-backend/pkg/config"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type CasbinService struct {
	enforcer     *casbin.Enforcer
	eventBus     eventbus.EventBus
	auditService *audit.AuditService
	mu           sync.RWMutex
}

func NewCasbinService(
	db *gorm.DB,
	eventBus eventbus.EventBus,
	auditService *audit.AuditService,
	config *config.CasbinConfig,
) (*CasbinService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(config.ModelFile, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	cs := &CasbinService{
		enforcer:     enforcer,
		eventBus:     eventBus,
		auditService: auditService,
	}

	eventBus.Subscribe(context.Background(), events.EventTypePolicyReload, cs.handlePolicyReloadEvent)

	if config.AutoSeedPolicies {
		if err := cs.seedInitialPolicies(); err != nil {
			logger.Error("Failed to seed initial policies", zap.Error(err))
		}
	}

	return cs, nil
}

func (s *CasbinService) handlePolicyReloadEvent(ctx context.Context, event events.Event) error {
	payload := event.GetPayload().(events.PolicyReloadEventPayload)

	logger.Info("Handling policy reload event",
		zap.String("triggered_by", payload.TriggeredBy),
		zap.String("reason", payload.Reason))

	return s.ReloadPolicies()
}

func (s *CasbinService) ReloadPolicies() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.enforcer.LoadPolicy(); err != nil {
		return err
	}

	logger.Info("Casbin policies reloaded successfully")
	return nil
}

func (s *CasbinService) Enforce(sub, obj, act, owner string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enforcer.Enforce(sub, obj, act, owner)
}

func (s *CasbinService) EnforceWithOwner(sub, obj, act, owner string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enforcer.Enforce(sub, obj, act, owner)
}

func (s *CasbinService) GetAllPolicies() ([][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (s *CasbinService) AddPolicy(ctx context.Context, sub, obj, act, owner string, performedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	added, err := s.enforcer.AddPolicy(sub, obj, act, owner)
	if err != nil {
		return err
	}
	if !added {
		return fmt.Errorf("policy already exists")
	}

	event := events.NewPolicyCreatedEvent(sub, obj, act, owner, performedBy)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		logger.Error("Failed to publish policy created event", zap.Error(err))
	}

	s.auditService.LogAction(ctx, audit.AuditEntry{
		UserID:       getPerformedByID(performedBy),
		Username:     performedBy,
		Action:       "policy.created",
		ResourceType: "policy",
		Details: map[string]interface{}{
			"subject": sub,
			"object":  obj,
			"action":  act,
			"owner":   owner,
		},
		Status: "success",
	})

	return nil
}

func (s *CasbinService) RemovePolicy(ctx context.Context, sub, obj, act, owner string, performedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	removed, err := s.enforcer.RemovePolicy(sub, obj, act, owner)
	if err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("policy not found")
	}

	event := events.NewPolicyDeletedEvent(sub, obj, act, owner, performedBy)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		logger.Error("Failed to publish policy deleted event", zap.Error(err))
	}

	s.auditService.LogAction(ctx, audit.AuditEntry{
		UserID:       getPerformedByID(performedBy),
		Username:     performedBy,
		Action:       "policy.deleted",
		ResourceType: "policy",
		Details: map[string]interface{}{
			"subject": sub,
			"object":  obj,
			"action":  act,
			"owner":   owner,
		},
		Status: "success",
	})

	return nil
}

func (s *CasbinService) GetRolesForUser(username string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.GetRolesForUser(username)
}

func (s *CasbinService) AddRoleForUser(username, role string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.AddRoleForUser(username, role)
}

func (s *CasbinService) DeleteRoleForUser(username, role string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.DeleteRoleForUser(username, role)
}

func (s *CasbinService) seedInitialPolicies() error {
	policies, _ := s.enforcer.GetPolicy()
	if len(policies) > 0 {
		return nil
	}

	defaultPolicies := [][]string{
		{"superadmin", "/api/*", "*", "*"},
		{"superadmin", "/api/casbin/*", "*", "*"},
		{"admin", "/api/routers", "*", "*"},
		{"admin", "/api/settings", "*", "*"},
		{"admin", "/api/admin/users", "*", "*"},
		{"admin", "/api/customers", "GET", "*"},
		{"admin", "/api/customers", "POST", "*"},
		{"admin", "/api/customers", "PUT", "*"},
		{"admin", "/api/invoices", "GET", "*"},
		{"admin", "/api/invoices", "POST", "*"},
		{"admin", "/api/mikrotik", "*", "*"},
		{"admin", "/api/hotspot", "*", "*"},
		{"operator", "/api/customers", "GET", "*"},
		{"operator", "/api/customers", "POST", "*"},
		{"operator", "/api/customers", "PUT", "*"},
		{"operator", "/api/invoices", "GET", "*"},
		{"operator", "/api/invoices", "POST", "*"},
		{"operator", "/api/mikrotik/ppp", "GET", "*"},
		{"operator", "/api/mikrotik/ppp", "POST", "*"},
		{"operator", "/api/hotspot/profiles", "GET", "*"},
		{"operator", "/api/hotspot/profiles", "POST", "*"},
		{"readonly", "/api/dashboard", "GET", "*"},
		{"readonly", "/api/customers", "GET", "*"},
		{"readonly", "/api/invoices", "GET", "*"},
	}

	for _, policy := range defaultPolicies {
		if _, err := s.enforcer.AddPolicy(policy); err != nil {
			return err
		}
	}

	logger.Info("Casbin initial policies seeded", zap.Int("count", len(defaultPolicies)))
	return nil
}

func getPerformedByID(performedBy string) int64 {
	if performedBy == "" || performedBy == "system" {
		return 0
	}
	return 0
}
