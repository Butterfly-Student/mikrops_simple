package usecase

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type MikroTikUsecase interface {
	GetRouters() ([]*mikrotik.RouterStatus, error)
	GetRouter(id uint) (*dto.RouterDetail, error)
	CreateRouter(router *dto.RouterCreate) error
	UpdateRouter(id uint, router *dto.RouterUpdate) error
	DeleteRouter(id uint) error
	TestRouterConnection(id uint) (*dto.ConnectionTestResult, error)
	ActivateRouter(id uint) error
	GetRouterStatus(id uint) (*dto.RouterStatus, error)
	GetAllRoutersStatus() ([]*dto.RouterStatus, error)

	GetPPPUsers(routerID uint) (*dto.PPPUsersResponse, error)
	GetActiveSessions(routerID uint) (*dto.ActiveSessionsResponse, error)
	GetPPPProfiles(routerID uint) (*dto.ProfilesResponse, error)
	AddPPPUser(req *dto.AddPPPUserRequest) error
	RemovePPPUser(username string, routerID uint) error
	UpdatePPPUser(username string, routerID uint, params map[string]interface{}) error
	DisconnectUser(username string, routerID uint) error

	IsolateCustomer(id uint) error
	ActivateCustomer(id uint) error
	BulkIsolate(ids []uint) error
	BulkActivate(ids []uint) error
	SyncCustomer(id uint) error
	SyncAllCustomers() error
}

type mikrotikUsecase struct {
	mikrotikService *mikrotik.MikroTikService
}

func NewMikroTikUsecase(mikrotikService *mikrotik.MikroTikService) MikroTikUsecase {
	return &mikrotikUsecase{
		mikrotikService: mikrotikService,
	}
}

func (u *mikrotikUsecase) GetRouters() ([]*mikrotik.RouterStatus, error) {
	statuses, err := u.mikrotikService.GetAllRoutersStatus()
	if err != nil {
		return nil, err
	}

	result := make([]*mikrotik.RouterStatus, len(statuses))
	for i := range statuses {
		result[i] = &statuses[i]
	}

	return result, nil
}

func (u *mikrotikUsecase) GetRouter(id uint) (*dto.RouterDetail, error) {
	return nil, nil
}

func (u *mikrotikUsecase) CreateRouter(router *dto.RouterCreate) error {
	return nil
}

func (u *mikrotikUsecase) UpdateRouter(id uint, router *dto.RouterUpdate) error {
	return nil
}

func (u *mikrotikUsecase) DeleteRouter(id uint) error {
	return nil
}

func (u *mikrotikUsecase) TestRouterConnection(id uint) (*dto.ConnectionTestResult, error) {
	err := u.mikrotikService.TestConnection(id)
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

func (u *mikrotikUsecase) ActivateRouter(id uint) error {
	return nil
}

func (u *mikrotikUsecase) GetRouterStatus(id uint) (*dto.RouterStatus, error) {
	return nil, nil
}

func (u *mikrotikUsecase) GetAllRoutersStatus() ([]*dto.RouterStatus, error) {
	statuses, err := u.mikrotikService.GetAllRoutersStatus()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.RouterStatus, len(statuses))
	for i, status := range statuses {
		result[i] = &dto.RouterStatus{
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

	return result, nil
}

func (u *mikrotikUsecase) GetPPPUsers(routerID uint) (*dto.PPPUsersResponse, error) {
	users, err := u.mikrotikService.GetPPPUsersByRouter(routerID)
	if err != nil {
		return nil, err
	}

	pppUsers := make([]dto.PPPUser, len(users))
	for i, user := range users {
		pppUsers[i] = dto.PPPUser{
			Name:      user.Name,
			Service:   user.Service,
			Profile:   user.Profile,
			CallerID:  user.CallerID,
			Disabled:  user.Disabled,
			LastLogin: user.LastLogin,
		}
	}

	return &dto.PPPUsersResponse{
		Users: pppUsers,
		Total: len(users),
	}, nil
}

func (u *mikrotikUsecase) GetActiveSessions(routerID uint) (*dto.ActiveSessionsResponse, error) {
	sessions, err := u.mikrotikService.GetActiveSessionsByRouter(routerID)
	if err != nil {
		return nil, err
	}

	activeSessions := make([]dto.ActiveSession, len(sessions))
	for i, session := range sessions {
		activeSessions[i] = dto.ActiveSession{
			Name:     session.Name,
			CallerID: session.CallerID,
			Address:  session.Address,
			Uptime:   session.Uptime,
			BytesIn:  session.BytesIn,
			BytesOut: session.BytesOut,
			Encoding: session.Encoding,
		}
	}

	return &dto.ActiveSessionsResponse{
		Sessions: activeSessions,
		Total:    len(sessions),
	}, nil
}

func (u *mikrotikUsecase) GetPPPProfiles(routerID uint) (*dto.ProfilesResponse, error) {
	profiles, err := u.mikrotikService.GetPPPProfilesByRouter(routerID)
	if err != nil {
		return nil, err
	}

	profileDTOs := make([]dto.Profile, len(profiles))
	for i, profile := range profiles {
		profileDTOs[i] = dto.Profile{
			Name:         profile.Name,
			RateLimit:    profile.RateLimit,
			LocalAddress: profile.LocalAddress,
			OnlyOne:      profile.OnlyOne,
		}
	}

	return &dto.ProfilesResponse{
		Profiles: profileDTOs,
		Total:    len(profiles),
	}, nil
}

func (u *mikrotikUsecase) AddPPPUser(req *dto.AddPPPUserRequest) error {
	return u.mikrotikService.AddPPPUser(req.Username, req.Password, req.Profile, req.RouterID)
}

func (u *mikrotikUsecase) RemovePPPUser(username string, routerID uint) error {
	return u.mikrotikService.RemovePPPUser(username, routerID)
}

func (u *mikrotikUsecase) UpdatePPPUser(username string, routerID uint, params map[string]interface{}) error {
	return u.mikrotikService.UpdatePPPUser(username, routerID, params)
}

func (u *mikrotikUsecase) DisconnectUser(username string, routerID uint) error {
	return u.mikrotikService.DisconnectUser(username, routerID)
}

func (u *mikrotikUsecase) IsolateCustomer(id uint) error {
	customer, err := u.mikrotikService.GetCustomerByID(id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	return u.mikrotikService.IsolateCustomer(customer)
}

func (u *mikrotikUsecase) ActivateCustomer(id uint) error {
	customer, err := u.mikrotikService.GetCustomerByID(id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	return u.mikrotikService.ActivateCustomer(customer)
}

func (u *mikrotikUsecase) BulkIsolate(ids []uint) error {
	for _, id := range ids {
		if customer, err := u.mikrotikService.GetCustomerByID(id); err == nil {
			_ = u.mikrotikService.IsolateCustomer(customer)
		}
	}
	return nil
}

func (u *mikrotikUsecase) BulkActivate(ids []uint) error {
	for _, id := range ids {
		if customer, err := u.mikrotikService.GetCustomerByID(id); err == nil {
			_ = u.mikrotikService.ActivateCustomer(customer)
		}
	}
	return nil
}

func (u *mikrotikUsecase) SyncCustomer(id uint) error {
	customer, err := u.mikrotikService.GetCustomerByID(id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	return u.mikrotikService.SyncCustomerToMikroTik(customer)
}

func (u *mikrotikUsecase) SyncAllCustomers() error {
	return u.mikrotikService.SyncAllCustomers()
}
