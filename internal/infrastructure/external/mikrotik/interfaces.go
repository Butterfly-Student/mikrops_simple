package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
)

type IMikroTikService interface {
	GetActiveRouter() (*entities.Router, error)
	TestConnection(routerID uint) error
	GetAllRoutersStatus() ([]RouterStatus, error)
	GetPPPUsers() ([]PPPoEUser, error)
	GetPPPUsersByRouter(routerID uint) ([]PPPoEUser, error)
	AddPPPUser(username, password, profile string, routerID uint) error
	RemovePPPUser(username string, routerID uint) error
	UpdatePPPUser(username string, routerID uint, params map[string]interface{}) error
	GetActiveSessions() ([]ActiveSession, error)
	GetActiveSessionsByRouter(routerID uint) ([]ActiveSession, error)
	DisconnectUser(username string, routerID uint) error
	GetPPPProfiles() ([]Profile, error)
	GetPPPProfilesByRouter(routerID uint) ([]Profile, error)
	CreateCustomerOnMikroTik(customer *entities.Customer) error
	DeleteCustomerFromMikroTik(customer *entities.Customer) error
	IsolateCustomer(customer *entities.Customer) error
	ActivateCustomer(customer *entities.Customer) error
	SyncCustomerToMikroTik(customer *entities.Customer) error
	BulkSyncCustomers(customerIDs []uint) error
	GetCustomerByID(id uint) (*entities.Customer, error)
	SyncAllCustomers() error
}
