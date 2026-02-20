package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
)

func newTestTicketUsecase(t *testing.T) (*TroubleTicketUsecase, *mocks.TroubleTicketRepository, *mocks.CustomerRepository) {
	ticketRepo := mocks.NewTroubleTicketRepository(t)
	customerRepo := mocks.NewCustomerRepository(t)
	uc := NewTroubleTicketUsecase(ticketRepo, customerRepo)
	return uc, ticketRepo, customerRepo
}

func TestTicketCreate_Success(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	ticketRepo.On("Create", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Create(CreateTicketRequest{
		CustomerID: 1, Subject: "Internet down", Description: "No connection", Priority: "high",
	})

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "high", ticket.Priority)
	assert.Equal(t, "open", ticket.Status)
}

func TestTicketCreate_EmptyDescription(t *testing.T) {
	uc, _, _ := newTestTicketUsecase(t)

	ticket, err := uc.Create(CreateTicketRequest{CustomerID: 1, Description: ""})

	assert.Error(t, err)
	assert.Nil(t, ticket)
	assert.Contains(t, err.Error(), "description is required")
}

func TestTicketCreate_DefaultPriority(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	ticketRepo.On("Create", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Create(CreateTicketRequest{
		CustomerID: 1, Subject: "Slow", Description: "Slow speed", Priority: "",
	})

	assert.NoError(t, err)
	assert.Equal(t, "medium", ticket.Priority)
}

func TestTicketCreate_InvalidPriority(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	ticketRepo.On("Create", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Create(CreateTicketRequest{
		CustomerID: 1, Subject: "Issue", Description: "Details", Priority: "critical",
	})

	assert.NoError(t, err)
	assert.Equal(t, "medium", ticket.Priority)
}

func TestTicketGetAll_WithStatus(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	expected := []*entities.TroubleTicket{{ID: 1, Status: "open"}, {ID: 2, Status: "open"}}
	ticketRepo.On("FindByStatus", "open", 1, 10).Return(expected, int64(2), nil)

	tickets, total, err := uc.GetAll(1, 10, "open")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, tickets, 2)
}

func TestTicketGetAll_NoStatus(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	expected := []*entities.TroubleTicket{{ID: 1}, {ID: 2}}
	ticketRepo.On("FindAll", 1, 10).Return(expected, int64(2), nil)

	tickets, total, err := uc.GetAll(1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, tickets, 2)
}

func TestTicketGetByID_Success(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	expected := &entities.TroubleTicket{ID: 1, Subject: "Internet down"}
	ticketRepo.On("FindByID", uint(1)).Return(expected, nil)

	ticket, err := uc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Internet down", ticket.Subject)
}

func TestTicketUpdate_Success(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	existing := &entities.TroubleTicket{ID: 1, Status: "open", AssignedTo: ""}
	ticketRepo.On("FindByID", uint(1)).Return(existing, nil)
	ticketRepo.On("Update", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Update(1, UpdateTicketRequest{Status: "in_progress", AssignedTo: "tech-A"})

	assert.NoError(t, err)
	assert.Equal(t, "in_progress", ticket.Status)
	assert.Equal(t, "tech-A", ticket.AssignedTo)
	assert.Nil(t, ticket.ResolvedAt)
}

func TestTicketUpdate_NotFound(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	ticketRepo.On("FindByID", uint(999)).Return(nil, fmt.Errorf("not found"))

	ticket, err := uc.Update(999, UpdateTicketRequest{Status: "closed"})

	assert.Error(t, err)
	assert.Nil(t, ticket)
	assert.Contains(t, err.Error(), "ticket not found")
}

func TestTicketUpdate_ResolveSetsTimestamp(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	existing := &entities.TroubleTicket{ID: 2, Status: "in_progress"}
	ticketRepo.On("FindByID", uint(2)).Return(existing, nil)
	ticketRepo.On("Update", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Update(2, UpdateTicketRequest{Status: "resolved"})

	assert.NoError(t, err)
	assert.Equal(t, "resolved", ticket.Status)
	assert.NotNil(t, ticket.ResolvedAt)
}

func TestTicketUpdate_ClosedSetsTimestamp(t *testing.T) {
	uc, ticketRepo, _ := newTestTicketUsecase(t)
	existing := &entities.TroubleTicket{ID: 3, Status: "in_progress"}
	ticketRepo.On("FindByID", uint(3)).Return(existing, nil)
	ticketRepo.On("Update", mock.AnythingOfType("*entities.TroubleTicket")).Return(nil)

	ticket, err := uc.Update(3, UpdateTicketRequest{Status: "closed"})

	assert.NoError(t, err)
	assert.Equal(t, "closed", ticket.Status)
	assert.NotNil(t, ticket.ResolvedAt)
}
