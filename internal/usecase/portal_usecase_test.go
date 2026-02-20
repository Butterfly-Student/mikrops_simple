package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
	"github.com/alijayanet/gembok-backend/pkg/utils"
)

func newTestPortalUsecase(t *testing.T) (*PortalUsecase, *mocks.CustomerRepository, *mocks.InvoiceRepository, *mocks.TroubleTicketRepository) {
	customerRepo := mocks.NewCustomerRepository(t)
	invoiceRepo := mocks.NewInvoiceRepository(t)
	ticketRepo := mocks.NewTroubleTicketRepository(t)
	uc := NewPortalUsecase(customerRepo, invoiceRepo, ticketRepo, "test-portal-secret")
	return uc, customerRepo, invoiceRepo, ticketRepo
}

func TestPortalLogin_Success(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	hashed, _ := utils.HashPassword("mypassword")
	customer := &entities.Customer{
		ID: 1, Name: "John", Phone: "08123", Status: "active", PPPoEPassword: hashed,
	}
	customerRepo.On("FindByPhone", "08123").Return(customer, nil)

	resp, err := uc.Login(PortalLoginRequest{Phone: "08123", Password: "mypassword"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "", resp.Customer.PPPoEPassword) // sensitive data hidden
}

func TestPortalLogin_PlainTextPasswordCompat(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	// Simulate a customer with plain text PPPoE password (backward compat)
	customer := &entities.Customer{
		ID: 1, Name: "John", Phone: "08123", Status: "active", PPPoEPassword: "plainpass",
	}
	customerRepo.On("FindByPhone", "08123").Return(customer, nil)

	resp, err := uc.Login(PortalLoginRequest{Phone: "08123", Password: "plainpass"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
}

func TestPortalLogin_PhoneNotFound(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)
	customerRepo.On("FindByPhone", "00000").Return(nil, errors.New("not found"))

	resp, err := uc.Login(PortalLoginRequest{Phone: "00000", Password: "pass"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "phone or password incorrect")
}

func TestPortalLogin_WrongPassword(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	hashed, _ := utils.HashPassword("correct")
	customer := &entities.Customer{
		ID: 1, Name: "John", Phone: "08123", Status: "active", PPPoEPassword: hashed,
	}
	customerRepo.On("FindByPhone", "08123").Return(customer, nil)

	resp, err := uc.Login(PortalLoginRequest{Phone: "08123", Password: "wrong"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "phone or password incorrect")
}

func TestPortalLogin_InactiveAccount(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	customer := &entities.Customer{
		ID: 1, Name: "John", Phone: "08123", Status: "inactive", PPPoEPassword: "pass",
	}
	customerRepo.On("FindByPhone", "08123").Return(customer, nil)

	resp, err := uc.Login(PortalLoginRequest{Phone: "08123", Password: "pass"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "account is inactive")
}

func TestPortalGetProfile_Success(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	customerRepo.On("FindByID", uint(1)).Return(&entities.Customer{
		ID: 1, Name: "John", Phone: "08123", PPPoEPassword: "should-be-hidden",
	}, nil)

	customer, err := uc.GetProfile(1)

	assert.NoError(t, err)
	assert.Equal(t, "John", customer.Name)
	assert.Equal(t, "", customer.PPPoEPassword) // sensitive data hidden
}

func TestPortalGetProfile_NotFound(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)
	customerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	customer, err := uc.GetProfile(999)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.Contains(t, err.Error(), "customer not found")
}

func TestPortalChangePassword_Success(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)

	customerRepo.On("FindByID", uint(1)).Return(&entities.Customer{
		ID: 1, Name: "John", PPPoEPassword: "old",
	}, nil)
	customerRepo.On("Update", mock.AnythingOfType("*entities.Customer")).Return(nil)

	err := uc.ChangePassword(1, "newpassword123")

	assert.NoError(t, err)
}

func TestPortalChangePassword_TooShort(t *testing.T) {
	uc, _, _, _ := newTestPortalUsecase(t)

	err := uc.ChangePassword(1, "abc")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be at least 6 characters")
}

func TestPortalChangePassword_NotFound(t *testing.T) {
	uc, customerRepo, _, _ := newTestPortalUsecase(t)
	customerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := uc.ChangePassword(999, "newpassword123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")
}

func TestPortalGetInvoices_Success(t *testing.T) {
	uc, _, invoiceRepo, _ := newTestPortalUsecase(t)

	invoiceRepo.On("FindByCustomerID", uint(1), 1, 10).Return(
		[]*entities.Invoice{{ID: 1, Number: "INV-001"}}, int64(1), nil,
	)

	invoices, total, err := uc.GetInvoices(1, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, invoices, 1)
}

func TestPortalGetTickets_Success(t *testing.T) {
	uc, _, _, ticketRepo := newTestPortalUsecase(t)

	ticketRepo.On("FindByCustomerID", uint(1)).Return(
		[]*entities.TroubleTicket{{ID: 1, Subject: "Internet down"}}, nil,
	)

	tickets, err := uc.GetTickets(1)

	assert.NoError(t, err)
	assert.Len(t, tickets, 1)
}

func TestPortalCreateTicket_Success(t *testing.T) {
	uc, _, _, ticketRepo := newTestPortalUsecase(t)

	ticketRepo.On("Create", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.CreateTicket(1, "No internet connection", "high", "Internet down")

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "open", ticket.Status)
	assert.Equal(t, "high", ticket.Priority)
	assert.Equal(t, "Internet down", ticket.Subject)
}

func TestPortalCreateTicket_EmptyDescription(t *testing.T) {
	uc, _, _, _ := newTestPortalUsecase(t)

	ticket, err := uc.CreateTicket(1, "", "high", "Issue")

	assert.Error(t, err)
	assert.Nil(t, ticket)
	assert.Contains(t, err.Error(), "description is required")
}
