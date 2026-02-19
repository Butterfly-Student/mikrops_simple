package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/whatsapp"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type InvoiceUsecase interface {
	GetInvoices(page, perPage int) (*dto.InvoiceListResponse, error)
	GetInvoiceByID(id uint) (*dto.InvoiceDetail, error)
	CreateInvoice(invoice *dto.InvoiceDetail) error
	UpdateInvoice(id uint, invoice *dto.InvoiceDetail) error
	DeleteInvoice(id uint) error
}

type invoiceUsecase struct {
	invoiceRepo     repositories.InvoiceRepository
	settingRepo     repositories.SettingRepository
	whatsappService *whatsapp.WhatsAppService
}

func NewInvoiceUsecase(invoiceRepo repositories.InvoiceRepository, settingRepo repositories.SettingRepository, whatsappService *whatsapp.WhatsAppService) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceRepo:     invoiceRepo,
		settingRepo:     settingRepo,
		whatsappService: whatsappService,
	}
}

func (u *invoiceUsecase) GetInvoices(page, perPage int) (*dto.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	invoices, total, err := u.invoiceRepo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	invoiceDTOs := make([]dto.InvoiceDetail, 0, len(invoices))
	for _, invoice := range invoices {
		invoiceDTOs = append(invoiceDTOs, *u.entityToDTO(invoice))
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.InvoiceListResponse{
		Invoices:   invoiceDTOs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (u *invoiceUsecase) GetInvoiceByID(id uint) (*dto.InvoiceDetail, error) {
	invoice, err := u.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return u.entityToDTO(invoice), nil
}

func (u *invoiceUsecase) CreateInvoice(invoiceDTO *dto.InvoiceDetail) error {
	prefix, _ := u.settingRepo.Get("INVOICE_PREFIX")
	if prefix == "" {
		prefix = "INV-"
	}

	lastInvoice, _ := u.invoiceRepo.FindLastInvoiceNumber()
	var nextNum int
	if lastInvoice != nil {
		numStr := strings.TrimPrefix(lastInvoice.Number, prefix)
		lastNum, _ := strconv.Atoi(numStr)
		nextNum = lastNum + 1
	} else {
		start, _ := u.settingRepo.Get("INVOICE_START")
		nextNum, _ = strconv.Atoi(start)
		if nextNum == 0 {
			nextNum = 1
		}
	}

	invoiceNumber := fmt.Sprintf("%s%06d", prefix, nextNum)

	var dueDate time.Time
	if invoiceDTO.DueDate != "" {
		dueDate, _ = time.Parse("2006-01-02", invoiceDTO.DueDate)
	} else {
		dueDate = time.Now().AddDate(0, 0, 7)
	}

	invoice := &entities.Invoice{
		CustomerID:       invoiceDTO.CustomerID,
		Number:           invoiceNumber,
		Amount:           invoiceDTO.Amount,
		Period:           invoiceDTO.Period,
		DueDate:          dueDate,
		Status:           "unpaid",
		PaymentMethod:    invoiceDTO.PaymentMethod,
		PaymentReference: invoiceDTO.PaymentReference,
	}

	if err := u.invoiceRepo.Create(invoice); err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	if u.whatsappService != nil {
		go u.whatsappService.SendInvoiceNotification(invoice)
	}

	return nil
}

func (u *invoiceUsecase) UpdateInvoice(id uint, invoiceDTO *dto.InvoiceDetail) error {
	invoice, err := u.invoiceRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("invoice not found")
	}

	if invoiceDTO.Amount > 0 {
		invoice.Amount = invoiceDTO.Amount
	}
	if invoiceDTO.Period != "" {
		invoice.Period = invoiceDTO.Period
	}
	if invoiceDTO.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", invoiceDTO.DueDate)
		if err == nil {
			invoice.DueDate = dueDate
		}
	}
	if invoiceDTO.Status != "" {
		invoice.Status = invoiceDTO.Status
	}
	if invoiceDTO.PaymentMethod != "" {
		invoice.PaymentMethod = invoiceDTO.PaymentMethod
	}
	if invoiceDTO.PaymentReference != "" {
		invoice.PaymentReference = invoiceDTO.PaymentReference
	}

	if invoice.Status == "paid" && invoice.PaidAt == nil {
		now := time.Now()
		invoice.PaidAt = &now
	}

	if err := u.invoiceRepo.Update(invoice); err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	if u.whatsappService != nil && invoice.Status == "paid" {
		go u.whatsappService.SendPaymentConfirmation(invoice)
	}

	return nil
}

func (u *invoiceUsecase) DeleteInvoice(id uint) error {
	return u.invoiceRepo.Delete(id)
}

func (u *invoiceUsecase) entityToDTO(invoice *entities.Invoice) *dto.InvoiceDetail {
	customerName := ""
	customerPhone := ""
	if invoice.Customer != nil {
		customerName = invoice.Customer.Name
		customerPhone = invoice.Customer.Phone
	}

	var paidAt *string
	if invoice.PaidAt != nil {
		date := invoice.PaidAt.Format("2006-01-02 15:04:05")
		paidAt = &date
	}

	return &dto.InvoiceDetail{
		ID:               invoice.ID,
		CustomerID:       invoice.CustomerID,
		CustomerName:     customerName,
		CustomerPhone:    customerPhone,
		Number:           invoice.Number,
		Amount:           invoice.Amount,
		Period:           invoice.Period,
		DueDate:          invoice.DueDate.Format("2006-01-02"),
		Status:           invoice.Status,
		PaidAt:           paidAt,
		PaymentMethod:    invoice.PaymentMethod,
		PaymentReference: invoice.PaymentReference,
		CreatedAt:        invoice.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        invoice.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
