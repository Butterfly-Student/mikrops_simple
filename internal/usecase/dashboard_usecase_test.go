package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
)

func newTestDashboardUsecase(t *testing.T) (DashboardUsecase, *mocks.CustomerRepository, *mocks.InvoiceRepository, *mocks.PackageRepository) {
	customerRepo := mocks.NewCustomerRepository(t)
	invoiceRepo := mocks.NewInvoiceRepository(t)
	packageRepo := mocks.NewPackageRepository(t)
	uc := NewDashboardUsecase(customerRepo, invoiceRepo, packageRepo)
	return uc, customerRepo, invoiceRepo, packageRepo
}

func TestDashboardGetStats_Success(t *testing.T) {
	uc, customerRepo, invoiceRepo, packageRepo := newTestDashboardUsecase(t)

	// Total customers (FindAll page 1, perPage 1)
	customerRepo.On("FindAll", 1, 1, "").Return(
		[]*entities.Customer{{ID: 1}}, int64(1), nil,
	).Once()

	// Active customers
	customerRepo.On("FindByStatus", "active", 1, 1).Return(
		[]*entities.Customer{{ID: 1}}, int64(1), nil,
	)

	// Isolated customers
	customerRepo.On("FindByStatus", "isolated", 1, 1).Return(
		[]*entities.Customer{}, int64(0), nil,
	)

	// Packages
	packageRepo.On("FindAll").Return(
		[]*entities.Package{{ID: 1}, {ID: 2}}, nil,
	)

	// Total invoices
	invoiceRepo.On("FindAll", 1, 1).Return(
		[]*entities.Invoice{{ID: 1}}, int64(1), nil,
	)

	// Paid invoices
	invoiceRepo.On("FindByStatus", "paid", 1, 100).Return(
		[]*entities.Invoice{{ID: 1, Amount: 100000}}, int64(1), nil,
	)

	// Unpaid invoices
	invoiceRepo.On("FindByStatus", "unpaid", 1, 100).Return(
		[]*entities.Invoice{}, int64(0), nil,
	)

	// Recent invoices (page 1, perPage 10)
	invoiceRepo.On("FindAll", 1, 10).Return(
		[]*entities.Invoice{{ID: 1, Number: "INV-001", Amount: 100000, Status: "paid"}}, int64(1), nil,
	)

	// Recent customers (page 1, perPage 5)
	customerRepo.On("FindAll", 1, 5, "").Return(
		[]*entities.Customer{{ID: 1, Name: "John", Phone: "08123", Status: "active"}}, int64(1), nil,
	)

	resp, err := uc.GetDashboardStats()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Stats.TotalCustomers)
	assert.Equal(t, int64(1), resp.Stats.ActiveCustomers)
	assert.Equal(t, int64(0), resp.Stats.IsolatedCustomers)
	assert.Equal(t, int64(2), resp.Stats.TotalPackages)
	assert.Equal(t, float64(100000), resp.Stats.TotalRevenue)
	assert.Len(t, resp.RecentInvoices, 1)
	assert.Len(t, resp.RecentCustomers, 1)
}

func TestDashboardGetStats_CustomerRepoError(t *testing.T) {
	uc, customerRepo, _, _ := newTestDashboardUsecase(t)

	customerRepo.On("FindAll", 1, 1, "").Return(
		nil, int64(0), errors.New("db connection error"),
	)

	resp, err := uc.GetDashboardStats()

	assert.Error(t, err)
	assert.Nil(t, resp)
}
