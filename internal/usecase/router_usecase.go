package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type RouterUsecase interface {
	GetAll() ([]*dto.RouterDetail, error)
	GetByID(id uint) (*dto.RouterDetail, error)
	Create(req *dto.RouterCreate) error
	Update(id uint, req *dto.RouterUpdate) error
	Delete(id uint) error
	TestConnection(id uint) (*dto.ConnectionTestResult, error)
	SetActive(id uint) error
	GetActive() (*dto.RouterDetail, error)
	GetStatus(id uint) (*dto.RouterStatus, error)
	GetAllStatus() ([]*dto.RouterStatus, error)
}

type routerUsecase struct {
	routerRepo     repositories.RouterRepository
	mikrotikClient MikroTikClientInterface
}

type MikroTikClientInterface interface {
	HealthCheck(routerID uint) error
	GetRouterStatus(routerID uint) (*mikrotik.RouterStatus, error)
	GetAllRoutersStatus() ([]mikrotik.RouterStatus, error)
}

func NewRouterUsecase(routerRepo repositories.RouterRepository, mikrotikClient MikroTikClientInterface) RouterUsecase {
	return &routerUsecase{
		routerRepo:     routerRepo,
		mikrotikClient: mikrotikClient,
	}
}

func (u *routerUsecase) GetAll() ([]*dto.RouterDetail, error) {
	routers, err := u.routerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.RouterDetail, len(routers))
	for i, router := range routers {
		dtos[i] = &dto.RouterDetail{
			ID:        router.ID,
			Name:      router.Name,
			Host:      router.Host,
			Username:  router.Username,
			Port:      router.Port,
			IsActive:  router.IsActive,
			CreatedAt: router.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: router.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return dtos, nil
}

func (u *routerUsecase) GetByID(id uint) (*dto.RouterDetail, error) {
	router, err := u.routerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.RouterDetail{
		ID:        router.ID,
		Name:      router.Name,
		Host:      router.Host,
		Username:  router.Username,
		Port:      router.Port,
		IsActive:  router.IsActive,
		CreatedAt: router.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: router.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (u *routerUsecase) Create(req *dto.RouterCreate) error {
	router := &entities.Router{
		Name:     req.Name,
		Host:     req.Host,
		Username: req.Username,
		Password: req.Password,
		Port:     req.Port,
		IsActive: req.IsActive,
	}

	return u.routerRepo.Create(router)
}

func (u *routerUsecase) Update(id uint, req *dto.RouterUpdate) error {
	router, err := u.routerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		router.Name = req.Name
	}
	if req.Host != "" {
		router.Host = req.Host
	}
	if req.Username != "" {
		router.Username = req.Username
	}
	if req.Password != "" {
		router.Password = req.Password
	}
	if req.Port != nil {
		router.Port = *req.Port
	}
	if req.IsActive != nil {
		router.IsActive = *req.IsActive
	}

	return u.routerRepo.Update(router)
}

func (u *routerUsecase) Delete(id uint) error {
	return u.routerRepo.Delete(id)
}

func (u *routerUsecase) TestConnection(id uint) (*dto.ConnectionTestResult, error) {
	err := u.mikrotikClient.HealthCheck(id)
	if err != nil {
		return &dto.ConnectionTestResult{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &dto.ConnectionTestResult{
		Success: true,
		Message: "Connection successful",
	}, nil
}

func (u *routerUsecase) SetActive(id uint) error {
	routers, err := u.routerRepo.FindAll()
	if err != nil {
		return err
	}

	for _, router := range routers {
		if router.ID == id {
			router.IsActive = true
		} else {
			router.IsActive = false
		}
		u.routerRepo.Update(router)
	}

	return nil
}

func (u *routerUsecase) GetActive() (*dto.RouterDetail, error) {
	router, err := u.routerRepo.FindActive()
	if err != nil {
		return nil, err
	}

	return &dto.RouterDetail{
		ID:        router.ID,
		Name:      router.Name,
		Host:      router.Host,
		Username:  router.Username,
		Port:      router.Port,
		IsActive:  router.IsActive,
		CreatedAt: router.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: router.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (u *routerUsecase) GetStatus(id uint) (*dto.RouterStatus, error) {
	status, err := u.mikrotikClient.GetRouterStatus(id)
	if err != nil {
		return nil, err
	}

	return &dto.RouterStatus{
		ID:          status.RouterID,
		Name:        status.Name,
		Host:        status.Host,
		Status:      status.Status,
		LastCheck:   status.LastCheck.Format("2006-01-02 15:04:05"),
		ActiveUsers: status.ActiveUsers,
		CPU:         status.CPU,
		Memory:      status.Memory,
		Uptime:      status.Uptime,
		Error:       status.Error,
	}, nil
}

func (u *routerUsecase) GetAllStatus() ([]*dto.RouterStatus, error) {
	statuses, err := u.mikrotikClient.GetAllRoutersStatus()
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.RouterStatus, len(statuses))
	for i, status := range statuses {
		dtos[i] = &dto.RouterStatus{
			ID:          status.RouterID,
			Name:        status.Name,
			Host:        status.Host,
			Status:      status.Status,
			LastCheck:   status.LastCheck.Format("2006-01-02 15:04:05"),
			ActiveUsers: status.ActiveUsers,
			CPU:         status.CPU,
			Memory:      status.Memory,
			Uptime:      status.Uptime,
			Error:       status.Error,
		}
	}

	return dtos, nil
}
