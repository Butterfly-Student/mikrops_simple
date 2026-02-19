package mikrotik

import (
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type QueueService struct {
	client     *MikroTikClient
	routerRepo repositories.RouterRepository
}

func NewQueueService(client *MikroTikClient, routerRepo repositories.RouterRepository) *QueueService {
	return &QueueService{
		client:     client,
		routerRepo: routerRepo,
	}
}
