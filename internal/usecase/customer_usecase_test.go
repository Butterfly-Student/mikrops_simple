package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

func newTestCustomerUsecase(t *testing.T) (CustomerUsecase, *mocks.CustomerRepository) {
	customerRepo := mocks.NewCustomerRepository(t)
	// Pass nil for mikrotikService and whatsappService.
	// Tests that exercise MikroTik/WhatsApp paths are skipped.
	uc := NewCustomerUsecase(customerRepo, nil, nil)
	return uc, customerRepo
}

func TestCustomerGetAll_Success(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	repo.On("FindAll", 1, 20, "").Return([]*entities.Customer{
		{ID: 1, Name: "John", Phone: "08123", Status: "active", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Jane", Phone: "08456", Status: "active", CreatedAt: now, UpdatedAt: now},
	}, int64(2), nil)

	result, err := uc.GetCustomers(1, 20, "")

	assert.NoError(t, err)
	assert.Len(t, result.Customers, 2)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 1, result.TotalPages)
}

func TestCustomerGetAll_DefaultPagination(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	repo.On("FindAll", 1, 20, "").Return([]*entities.Customer{
		{ID: 1, Name: "John", CreatedAt: now, UpdatedAt: now},
	}, int64(1), nil)

	// page=0 and perPage=0 corrected to defaults
	result, err := uc.GetCustomers(0, 0, "")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PerPage)
}

func TestCustomerGetAll_PerPageOverLimit(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	repo.On("FindAll", 1, 20, "").Return([]*entities.Customer{
		{ID: 1, Name: "John", CreatedAt: now, UpdatedAt: now},
	}, int64(1), nil)

	// perPage > 100 should be corrected to 20
	result, err := uc.GetCustomers(1, 200, "")

	assert.NoError(t, err)
	assert.Equal(t, 20, result.PerPage)
}

func TestCustomerGetAll_WithSearch(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	repo.On("FindAll", 1, 20, "john").Return([]*entities.Customer{
		{ID: 1, Name: "John", CreatedAt: now, UpdatedAt: now},
	}, int64(1), nil)

	result, err := uc.GetCustomers(1, 20, "john")

	assert.NoError(t, err)
	assert.Len(t, result.Customers, 1)
}

func TestCustomerGetAll_RepoError(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("FindAll", 1, 20, "").Return(nil, int64(0), errors.New("db error"))

	result, err := uc.GetCustomers(1, 20, "")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCustomerGetByID_Success(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	repo.On("FindByID", uint(1)).Return(&entities.Customer{
		ID: 1, Name: "John", Phone: "08123", Status: "active",
		CreatedAt: now, UpdatedAt: now,
	}, nil)

	result, err := uc.GetCustomerByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "", result.PPPoEPassword) // password hidden in DTO
}

func TestCustomerGetByID_NotFound(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := uc.GetCustomerByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCustomerCreate_SuccessWithoutPPPoE(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("Create", mock.AnythingOfType("*entities.Customer")).Return(nil)

	// Without PPPoE credentials, MikroTik is not called
	err := uc.CreateCustomer(&dto.CustomerDetail{
		Name: "John", Phone: "08123", PackageID: 1,
	})

	assert.NoError(t, err)
}

func TestCustomerCreate_RepoError(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("Create", mock.AnythingOfType("*entities.Customer")).Return(errors.New("duplicate"))

	err := uc.CreateCustomer(&dto.CustomerDetail{
		Name: "John", Phone: "08123",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestCustomerUpdate_Success(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	existing := &entities.Customer{
		ID: 1, Name: "Old", Phone: "08123", Status: "active",
		CreatedAt: now, UpdatedAt: now,
	}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*entities.Customer")).Return(nil)

	err := uc.UpdateCustomer(1, &dto.CustomerDetail{
		Name: "New", Phone: "08456", Email: "new@test.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, "New", existing.Name)
	assert.Equal(t, "08456", existing.Phone)
}

func TestCustomerUpdate_NotFound(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := uc.UpdateCustomer(999, &dto.CustomerDetail{Name: "X"})

	assert.Error(t, err)
}

func TestCustomerDelete_Success(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("Delete", uint(1)).Return(nil)

	err := uc.DeleteCustomer(1)

	assert.NoError(t, err)
}

func TestCustomerDelete_RepoError(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)
	repo.On("Delete", uint(999)).Return(errors.New("not found"))

	err := uc.DeleteCustomer(999)

	assert.Error(t, err)
}

func TestCustomerEntityToDTO_WithPackage(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	isoDate := now.Add(-24 * time.Hour)
	repo.On("FindByID", uint(1)).Return(&entities.Customer{
		ID: 1, Name: "John", Phone: "08123",
		Package:       &entities.Package{Name: "Premium", Price: 250000},
		IsolationDate: &isoDate,
		CreatedAt:     now, UpdatedAt: now,
	}, nil)

	result, err := uc.GetCustomerByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Premium", result.PackageName)
	assert.Equal(t, float64(250000), result.PackagePrice)
	assert.NotNil(t, result.IsolationDate)
}

func TestCustomerGetAll_PaginationCalc(t *testing.T) {
	uc, repo := newTestCustomerUsecase(t)

	now := time.Now()
	customers := make([]*entities.Customer, 0)
	for i := 1; i <= 10; i++ {
		customers = append(customers, &entities.Customer{
			ID: uint(i), Name: "Customer", CreatedAt: now, UpdatedAt: now,
		})
	}
	repo.On("FindAll", 1, 3, "").Return(customers[:3], int64(10), nil)

	result, err := uc.GetCustomers(1, 3, "")

	assert.NoError(t, err)
	assert.Len(t, result.Customers, 3)
	assert.Equal(t, int64(10), result.Total)
	assert.Equal(t, 4, result.TotalPages) // ceil(10/3) = 4
}
