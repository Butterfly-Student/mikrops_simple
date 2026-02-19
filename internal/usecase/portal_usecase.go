package usecase

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/pkg/utils"
)

// PortalUsecase handles customer portal operations
type PortalUsecase struct {
	customerRepo repositories.CustomerRepository
	invoiceRepo  repositories.InvoiceRepository
	ticketRepo   repositories.TroubleTicketRepository
	jwtSecret    string
	jwtExpiry    interface{} // time.Duration
}

func NewPortalUsecase(
	customerRepo repositories.CustomerRepository,
	invoiceRepo repositories.InvoiceRepository,
	ticketRepo repositories.TroubleTicketRepository,
	jwtSecret string,
) *PortalUsecase {
	return &PortalUsecase{
		customerRepo: customerRepo,
		invoiceRepo:  invoiceRepo,
		ticketRepo:   ticketRepo,
		jwtSecret:    jwtSecret,
	}
}

type PortalLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PortalLoginResponse struct {
	Token    string             `json:"token"`
	Customer *entities.Customer `json:"customer"`
}

func (u *PortalUsecase) Login(req PortalLoginRequest) (*PortalLoginResponse, error) {
	customer, err := u.customerRepo.FindByPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("phone or password incorrect")
	}

	if customer.Status == "inactive" {
		return nil, fmt.Errorf("account is inactive")
	}

	// Portal uses PPPoE password as portal password
	if !utils.CheckPassword(req.Password, customer.PPPoEPassword) {
		// Try plain text comparison for backward compat
		if req.Password != customer.PPPoEPassword {
			return nil, fmt.Errorf("phone or password incorrect")
		}
	}

	token, err := utils.GenerateToken(
		customer.ID,
		customer.Phone,
		"customer",
		u.jwtSecret,
		24*60*60*1000000000, // 24h in nanoseconds (time.Duration)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Hide sensitive data
	customer.PPPoEPassword = ""

	return &PortalLoginResponse{
		Token:    token,
		Customer: customer,
	}, nil
}

func (u *PortalUsecase) GetProfile(customerID uint) (*entities.Customer, error) {
	customer, err := u.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}
	customer.PPPoEPassword = ""
	return customer, nil
}

func (u *PortalUsecase) ChangePassword(customerID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	customer, err := u.customerRepo.FindByID(customerID)
	if err != nil {
		return fmt.Errorf("customer not found")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	customer.PPPoEPassword = hashed
	return u.customerRepo.Update(customer)
}

func (u *PortalUsecase) GetInvoices(customerID uint, page, perPage int) ([]*entities.Invoice, int64, error) {
	return u.invoiceRepo.FindByCustomerID(customerID, page, perPage)
}

func (u *PortalUsecase) GetTickets(customerID uint) ([]*entities.TroubleTicket, error) {
	return u.ticketRepo.FindByCustomerID(customerID)
}

func (u *PortalUsecase) CreateTicket(customerID uint, description, priority, subject string) (*entities.TroubleTicket, error) {
	if description == "" {
		return nil, fmt.Errorf("description is required")
	}

	ticket := &entities.TroubleTicket{
		CustomerID:  customerID,
		Subject:     subject,
		Description: description,
		Priority:    priority,
		Status:      "open",
	}

	if err := u.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}
