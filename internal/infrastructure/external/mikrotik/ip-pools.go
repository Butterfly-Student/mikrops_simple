package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type IPPoolsService struct {
	client     *MikroTikClient
	routerRepo repositories.RouterRepository
}

func NewIPPoolsService(client *MikroTikClient, routerRepo repositories.RouterRepository) *IPPoolsService {
	return &IPPoolsService{
		client:     client,
		routerRepo: routerRepo,
	}
}
