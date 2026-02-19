package mikrotik

import (
	"fmt"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

type MikroTikService struct {
	client       *MikroTikClient
	customerRepo repositories.CustomerRepository
	packageRepo  repositories.PackageRepository
	routerRepo   repositories.RouterRepository
}

func NewMikroTikService(client *MikroTikClient, customerRepo repositories.CustomerRepository, packageRepo repositories.PackageRepository, routerRepo repositories.RouterRepository) *MikroTikService {
	return &MikroTikService{
		client:       client,
		customerRepo: customerRepo,
		packageRepo:  packageRepo,
		routerRepo:   routerRepo,
	}
}

func (s *MikroTikService) GetActiveRouter() (*entities.Router, error) {
	return s.routerRepo.FindActive()
}

func (s *MikroTikService) TestConnection(routerID uint) error {
	return s.client.HealthCheck(routerID)
}

func (s *MikroTikService) GetAllRoutersStatus() ([]RouterStatus, error) {
	return s.client.GetAllRoutersStatus()
}

func (s *MikroTikService) GetPPPUsers() ([]PPPoEUser, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetAllUsers(routerID)
}

func (s *MikroTikService) GetPPPUsersByRouter(routerID uint) ([]PPPoEUser, error) {
	_, err := s.client.GetClient(routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetAllUsers(routerID)
}

func (s *MikroTikService) AddPPPUser(username, password, profile string, routerID uint) error {
	return s.client.AddUser(routerID, username, password, profile)
}

func (s *MikroTikService) RemovePPPUser(username string, routerID uint) error {
	return s.client.RemoveUser(routerID, username)
}

func (s *MikroTikService) UpdatePPPUser(username string, routerID uint, params map[string]interface{}) error {
	args := make([]string, 0, len(params))
	for key, val := range params {
		args = append(args, fmt.Sprintf("=%s=%v", key, val))
	}
	return s.client.UpdateUser(routerID, username, args...)
}

func (s *MikroTikService) GetActiveSessions() ([]ActiveSession, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetActiveSessions(routerID)
}

func (s *MikroTikService) GetActiveSessionsByRouter(routerID uint) ([]ActiveSession, error) {
	return s.client.GetActiveSessions(routerID)
}

func (s *MikroTikService) DisconnectUser(username string, routerID uint) error {
	return s.client.DisconnectUser(routerID, username)
}

func (s *MikroTikService) GetPPPProfiles() ([]Profile, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetAllProfiles(routerID)
}

func (s *MikroTikService) GetPPPProfilesByRouter(routerID uint) ([]Profile, error) {
	return s.client.GetAllProfiles(routerID)
}

func (s *MikroTikService) CreateCustomerOnMikroTik(customer *entities.Customer) error {
	if customer.RouterID == 0 {
		router, err := s.routerRepo.FindActive()
		if err != nil {
			return fmt.Errorf("no router specified and no active router found: %w", err)
		}
		customer.RouterID = router.ID
	}

	if customer.PPPoEUsername == "" {
		return fmt.Errorf("PPPoE username is required")
	}

	if customer.PPPoEPassword == "" {
		return fmt.Errorf("PPPoE password is required")
	}

	err := s.client.AddUser(customer.RouterID, customer.PPPoEUsername, customer.PPPoEPassword, "default")
	if err != nil {
		return fmt.Errorf("failed to create PPPoE user on MikroTik: %w", err)
	}

	logger.Info("Created customer on MikroTik",
		zap.Uint("customer_id", customer.ID),
		zap.Uint("router_id", customer.RouterID),
		zap.String("username", customer.PPPoEUsername),
	)

	return nil
}

func (s *MikroTikService) DeleteCustomerFromMikroTik(customer *entities.Customer) error {
	if customer.RouterID == 0 {
		return fmt.Errorf("customer has no router assigned")
	}

	if customer.PPPoEUsername == "" {
		return fmt.Errorf("customer has no PPPoE username")
	}

	err := s.client.RemoveUser(customer.RouterID, customer.PPPoEUsername)
	if err != nil {
		return fmt.Errorf("failed to delete PPPoE user from MikroTik: %w", err)
	}

	logger.Info("Deleted customer from MikroTik",
		zap.Uint("customer_id", customer.ID),
		zap.Uint("router_id", customer.RouterID),
		zap.String("username", customer.PPPoEUsername),
	)

	return nil
}

func (s *MikroTikService) IsolateCustomer(customer *entities.Customer) error {
	if customer.RouterID == 0 {
		return fmt.Errorf("customer has no router assigned")
	}

	if customer.PPPoEUsername == "" {
		return fmt.Errorf("customer has no PPPoE username")
	}

	if customer.PackageID == 0 {
		return fmt.Errorf("customer has no package assigned")
	}

	pkg, err := s.packageRepo.FindByID(customer.PackageID)
	if err != nil {
		return fmt.Errorf("failed to get customer package: %w", err)
	}

	if pkg.ProfileIsolir == "" {
		return fmt.Errorf("package has no isolation profile configured")
	}

	err = s.client.SetActiveProfile(customer.RouterID, customer.PPPoEUsername, pkg.ProfileIsolir)
	if err != nil {
		return fmt.Errorf("failed to isolate customer on MikroTik: %w", err)
	}

	now := time.Now()
	customer.Status = "isolated"
	customer.IsolationDate = &now

	err = s.customerRepo.Update(customer)
	if err != nil {
		logger.Error("Failed to update customer status",
			zap.Uint("customer_id", customer.ID),
			zap.Error(err),
		)
	}

	logger.Info("Customer isolated on MikroTik",
		zap.Uint("customer_id", customer.ID),
		zap.Uint("router_id", customer.RouterID),
		zap.String("username", customer.PPPoEUsername),
		zap.String("profile", pkg.ProfileIsolir),
	)

	return nil
}

func (s *MikroTikService) ActivateCustomer(customer *entities.Customer) error {
	if customer.RouterID == 0 {
		return fmt.Errorf("customer has no router assigned")
	}

	if customer.PPPoEUsername == "" {
		return fmt.Errorf("customer has no PPPoE username")
	}

	if customer.PackageID == 0 {
		return fmt.Errorf("customer has no package assigned")
	}

	pkg, err := s.packageRepo.FindByID(customer.PackageID)
	if err != nil {
		return fmt.Errorf("failed to get customer package: %w", err)
	}

	if pkg.ProfileNormal == "" {
		return fmt.Errorf("package has no normal profile configured")
	}

	err = s.client.SetActiveProfile(customer.RouterID, customer.PPPoEUsername, pkg.ProfileNormal)
	if err != nil {
		return fmt.Errorf("failed to activate customer on MikroTik: %w", err)
	}

	now := time.Now()
	customer.Status = "active"
	customer.ActivationDate = &now
	customer.IsolationDate = nil

	err = s.customerRepo.Update(customer)
	if err != nil {
		logger.Error("Failed to update customer status",
			zap.Uint("customer_id", customer.ID),
			zap.Error(err),
		)
	}

	logger.Info("Customer activated on MikroTik",
		zap.Uint("customer_id", customer.ID),
		zap.Uint("router_id", customer.RouterID),
		zap.String("username", customer.PPPoEUsername),
		zap.String("profile", pkg.ProfileNormal),
	)

	return nil
}

func (s *MikroTikService) SyncCustomerToMikroTik(customer *entities.Customer) error {
	if customer.Status == "active" {
		return s.ActivateCustomer(customer)
	} else if customer.Status == "isolated" {
		return s.IsolateCustomer(customer)
	}
	return nil
}

func (s *MikroTikService) BulkSyncCustomers(customerIDs []uint) error {
	for _, id := range customerIDs {
		customer, err := s.customerRepo.FindByID(id)
		if err != nil {
			continue
		}
		_ = s.SyncCustomerToMikroTik(customer)
	}
	return nil
}

// GetCustomerByID fetches a customer entity by ID.
func (s *MikroTikService) GetCustomerByID(id uint) (*entities.Customer, error) {
	return s.customerRepo.FindByID(id)
}

// SyncAllCustomers syncs every customer to MikroTik.
func (s *MikroTikService) SyncAllCustomers() error {
	customers, _, err := s.customerRepo.FindAll(1, 1000, "")
	if err != nil {
		return err
	}
	for _, c := range customers {
		_ = s.SyncCustomerToMikroTik(c)
	}
	return nil
}
