package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type HotspotService struct {
	client     *MikroTikClient
	routerRepo repositories.RouterRepository
}

func NewHotspotService(client *MikroTikClient, routerRepo repositories.RouterRepository) *HotspotService {
	return &HotspotService{
		client:     client,
		routerRepo: routerRepo,
	}
}

func (s *HotspotService) GetActiveSessions() ([]ActiveSession, error) {
	_, routerID, err := s.client.GetActiveRouter()
	if err != nil {
		return nil, err
	}

	return s.client.GetActiveSessions(routerID)
}

func (s *HotspotService) GetActiveSessionsByRouter(routerID uint) ([]ActiveSession, error) {
	return s.client.GetActiveSessions(routerID)
}

func (s *HotspotService) DisconnectUser(username string, routerID uint) error {
	return s.client.DisconnectUser(routerID, username)
}
