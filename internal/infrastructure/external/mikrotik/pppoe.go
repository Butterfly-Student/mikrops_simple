package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type PPPoEService struct {
	client       *MikroTikClient
	customerRepo repositories.CustomerRepository
	routerRepo   repositories.RouterRepository
}

func NewPPPoEService(client *MikroTikClient, customerRepo repositories.CustomerRepository, routerRepo repositories.RouterRepository) *PPPoEService {
	return &PPPoEService{
		client:       client,
		customerRepo: customerRepo,
		routerRepo:   routerRepo,
	}
}

func (s *PPPoEService) GetPPPUsers() ([]PPPoEUser, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetAllUsers(routerID)
}

func (s *PPPoEService) GetPPPUsersByRouter(routerID uint) ([]PPPoEUser, error) {
	return s.client.GetAllUsers(routerID)
}

func (s *PPPoEService) AddPPPUser(username, password, profile string, routerID uint) error {
	return s.client.AddUser(routerID, username, password, profile)
}

func (s *PPPoEService) RemovePPPUser(username string, routerID uint) error {
	return s.client.RemoveUser(routerID, username)
}

func (s *PPPoEService) UpdatePPPUser(username string, routerID uint, args ...string) error {
	return s.client.UpdateUser(routerID, username, args...)
}

func (s *PPPoEService) GetPPPProfiles() ([]Profile, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetAllProfiles(routerID)
}

func (s *PPPoEService) GetPPPProfilesByRouter(routerID uint) ([]Profile, error) {
	return s.client.GetAllProfiles(routerID)
}
