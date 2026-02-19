package usecase

import (
	"fmt"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type TroubleTicketUsecase struct {
	ticketRepo   repositories.TroubleTicketRepository
	customerRepo repositories.CustomerRepository
}

func NewTroubleTicketUsecase(ticketRepo repositories.TroubleTicketRepository, customerRepo repositories.CustomerRepository) *TroubleTicketUsecase {
	return &TroubleTicketUsecase{
		ticketRepo:   ticketRepo,
		customerRepo: customerRepo,
	}
}

type CreateTicketRequest struct {
	CustomerID  uint   `json:"customer_id"`
	Subject     string `json:"subject"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority"`
}

func (u *TroubleTicketUsecase) Create(req CreateTicketRequest) (*entities.TroubleTicket, error) {
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
	if !validPriorities[priority] {
		priority = "medium"
	}

	ticket := &entities.TroubleTicket{
		CustomerID:  req.CustomerID,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    priority,
		Status:      "open",
	}

	if err := u.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

func (u *TroubleTicketUsecase) GetAll(page, perPage int, status string) ([]*entities.TroubleTicket, int64, error) {
	if status != "" {
		return u.ticketRepo.FindByStatus(status, page, perPage)
	}
	return u.ticketRepo.FindAll(page, perPage)
}

func (u *TroubleTicketUsecase) GetByID(id uint) (*entities.TroubleTicket, error) {
	return u.ticketRepo.FindByID(id)
}

func (u *TroubleTicketUsecase) GetByCustomerID(customerID uint) ([]*entities.TroubleTicket, error) {
	return u.ticketRepo.FindByCustomerID(customerID)
}

type UpdateTicketRequest struct {
	Status     string `json:"status"`
	AssignedTo string `json:"assigned_to"`
}

func (u *TroubleTicketUsecase) Update(id uint, req UpdateTicketRequest) (*entities.TroubleTicket, error) {
	ticket, err := u.ticketRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("ticket not found")
	}

	if req.Status != "" {
		ticket.Status = req.Status
		if req.Status == "resolved" || req.Status == "closed" {
			now := time.Now()
			ticket.ResolvedAt = &now
		}
	}

	if req.AssignedTo != "" {
		ticket.AssignedTo = req.AssignedTo
	}

	if err := u.ticketRepo.Update(ticket); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return ticket, nil
}
