package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

// mockMikroTikClient implements MikroTikClientInterface for testing
type mockMikroTikClient struct {
	mock.Mock
}

func (m *mockMikroTikClient) HealthCheck(routerID uint) error {
	args := m.Called(routerID)
	return args.Error(0)
}

func (m *mockMikroTikClient) GetRouterStatus(routerID uint) (*mikrotik.RouterStatus, error) {
	args := m.Called(routerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mikrotik.RouterStatus), args.Error(1)
}

func (m *mockMikroTikClient) GetAllRoutersStatus() ([]mikrotik.RouterStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]mikrotik.RouterStatus), args.Error(1)
}

func newTestRouterUsecase(t *testing.T) (RouterUsecase, *mocks.RouterRepository, *mockMikroTikClient) {
	routerRepo := mocks.NewRouterRepository(t)
	mtClient := new(mockMikroTikClient)
	uc := NewRouterUsecase(routerRepo, mtClient)
	return uc, routerRepo, mtClient
}

func TestRouterGetAll_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)

	now := time.Now()
	routerRepo.On("FindAll").Return([]*entities.Router{
		{ID: 1, Name: "Router-1", Host: "10.0.0.1", Username: "admin", Port: 8728, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Router-2", Host: "10.0.0.2", Username: "admin", Port: 8728, IsActive: false, CreatedAt: now, UpdatedAt: now},
	}, nil)

	result, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Router-1", result[0].Name)
	assert.Equal(t, true, result[0].IsActive)
}

func TestRouterGetAll_RepoError(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)
	routerRepo.On("FindAll").Return(nil, errors.New("db error"))

	result, err := uc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRouterGetByID_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)

	now := time.Now()
	routerRepo.On("FindByID", uint(1)).Return(&entities.Router{
		ID: 1, Name: "Router-1", Host: "10.0.0.1", Username: "admin", Port: 8728, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}, nil)

	result, err := uc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Router-1", result.Name)
	assert.Equal(t, "10.0.0.1", result.Host)
}

func TestRouterGetByID_NotFound(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)
	routerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := uc.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRouterCreate_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)
	routerRepo.On("Create", mock.AnythingOfType("*entities.Router")).Return(nil)

	err := uc.Create(&dto.RouterCreate{
		Name: "Router-1", Host: "10.0.0.1", Username: "admin", Password: "secret", Port: 8728, IsActive: true,
	})

	assert.NoError(t, err)
}

func TestRouterUpdate_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)

	now := time.Now()
	existing := &entities.Router{ID: 1, Name: "Old", Host: "10.0.0.1", Username: "admin", Port: 8728, CreatedAt: now, UpdatedAt: now}
	routerRepo.On("FindByID", uint(1)).Return(existing, nil)
	routerRepo.On("Update", mock.AnythingOfType("*entities.Router")).Return(nil)

	err := uc.Update(1, &dto.RouterUpdate{Name: "New"})

	assert.NoError(t, err)
	assert.Equal(t, "New", existing.Name)
}

func TestRouterUpdate_NotFound(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)
	routerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := uc.Update(999, &dto.RouterUpdate{Name: "X"})

	assert.Error(t, err)
}

func TestRouterDelete_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)
	routerRepo.On("Delete", uint(1)).Return(nil)

	err := uc.Delete(1)

	assert.NoError(t, err)
}

func TestRouterTestConnection_Success(t *testing.T) {
	uc, _, mtClient := newTestRouterUsecase(t)
	mtClient.On("HealthCheck", uint(1)).Return(nil)

	result, err := uc.TestConnection(1)

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Connection successful", result.Message)
}

func TestRouterTestConnection_Failure(t *testing.T) {
	uc, _, mtClient := newTestRouterUsecase(t)
	mtClient.On("HealthCheck", uint(1)).Return(errors.New("connection refused"))

	result, err := uc.TestConnection(1)

	assert.NoError(t, err) // No error returned; failure is in the result
	assert.False(t, result.Success)
	assert.Contains(t, result.Message, "connection refused")
}

func TestRouterSetActive_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)

	now := time.Now()
	routers := []*entities.Router{
		{ID: 1, Name: "R1", IsActive: false, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "R2", IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
	routerRepo.On("FindAll").Return(routers, nil)
	routerRepo.On("Update", mock.AnythingOfType("*entities.Router")).Return(nil)

	err := uc.SetActive(1)

	assert.NoError(t, err)
}

func TestRouterGetActive_Success(t *testing.T) {
	uc, routerRepo, _ := newTestRouterUsecase(t)

	now := time.Now()
	routerRepo.On("FindActive").Return(&entities.Router{
		ID: 1, Name: "R1", Host: "10.0.0.1", Username: "admin", Port: 8728, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}, nil)

	result, err := uc.GetActive()

	assert.NoError(t, err)
	assert.Equal(t, "R1", result.Name)
	assert.True(t, result.IsActive)
}

func TestRouterGetStatus_Success(t *testing.T) {
	uc, _, mtClient := newTestRouterUsecase(t)

	now := time.Now()
	mtClient.On("GetRouterStatus", uint(1)).Return(&mikrotik.RouterStatus{
		RouterID: 1, Name: "R1", Host: "10.0.0.1", Status: "online",
		LastCheck: now, ActiveUsers: 10, CPU: 25.0, Memory: 512, Uptime: "3d12h",
	}, nil)

	result, err := uc.GetStatus(1)

	assert.NoError(t, err)
	assert.Equal(t, "online", result.Status)
	assert.Equal(t, 10, result.ActiveUsers)
}

func TestRouterGetAllStatus_Success(t *testing.T) {
	uc, _, mtClient := newTestRouterUsecase(t)

	now := time.Now()
	mtClient.On("GetAllRoutersStatus").Return([]mikrotik.RouterStatus{
		{RouterID: 1, Name: "R1", Host: "10.0.0.1", Status: "online", LastCheck: now},
		{RouterID: 2, Name: "R2", Host: "10.0.0.2", Status: "offline", LastCheck: now},
	}, nil)

	result, err := uc.GetAllStatus()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "online", result[0].Status)
	assert.Equal(t, "offline", result[1].Status)
}
