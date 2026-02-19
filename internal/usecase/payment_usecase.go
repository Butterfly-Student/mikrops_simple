package usecase

import (
	"fmt"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/tripay"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

type PaymentUsecase struct {
	invoiceRepo  repositories.InvoiceRepository
	customerRepo repositories.CustomerRepository
	tripay       *tripay.TripayClient
	mikrotikSvc  *mikrotik.MikroTikService
	appURL       string
}

func NewPaymentUsecase(
	invoiceRepo repositories.InvoiceRepository,
	customerRepo repositories.CustomerRepository,
	tripay *tripay.TripayClient,
	mikrotikSvc *mikrotik.MikroTikService,
	appURL string,
) *PaymentUsecase {
	return &PaymentUsecase{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		tripay:       tripay,
		mikrotikSvc:  mikrotikSvc,
		appURL:       appURL,
	}
}

type CreatePaymentRequest struct {
	InvoiceID     uint   `json:"invoice_id"`
	PaymentMethod string `json:"payment_method"` // e.g. "QRIS", "BRIVA", "MANDIRIVA"
}

type CreatePaymentResponse struct {
	PaymentURL  string `json:"payment_url"`
	Reference   string `json:"reference"`
	MerchantRef string `json:"merchant_ref"`
	Amount      int64  `json:"amount"`
	ExpiredTime int64  `json:"expired_time"`
}

func (u *PaymentUsecase) CreateTransaction(req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	invoice, err := u.invoiceRepo.FindByID(req.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found")
	}

	if invoice.Status == "paid" {
		return nil, fmt.Errorf("invoice is already paid")
	}

	if !u.tripay.IsConfigured() {
		return nil, fmt.Errorf("payment gateway not configured, please set tripay credentials in config")
	}

	customer := invoice.Customer
	if customer == nil {
		cust, err := u.customerRepo.FindByID(invoice.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("customer not found")
		}
		customer = cust
	}

	expiredTime := time.Now().Add(24 * time.Hour).Unix()

	tripayReq := tripay.TripayTransactionRequest{
		Method:        req.PaymentMethod,
		MerchantRef:   invoice.Number,
		Amount:        int64(invoice.Amount),
		CustomerName:  customer.Name,
		CustomerEmail: customer.Email,
		CustomerPhone: customer.Phone,
		OrderItems: []tripay.TripayOrderItem{
			{
				SKU:      invoice.Number,
				Name:     fmt.Sprintf("Tagihan Internet - %s", invoice.Period),
				Price:    int64(invoice.Amount),
				Quantity: 1,
			},
		},
		ReturnURL:   u.appURL,
		ExpiredTime: expiredTime,
	}

	resp, err := u.tripay.CreateTransaction(tripayReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	if !resp.Success || resp.Data == nil {
		return nil, fmt.Errorf("payment gateway error: %s", resp.Message)
	}

	// Save reference to invoice
	invoice.PaymentReference = resp.Data.Reference
	invoice.PaymentMethod = req.PaymentMethod
	_ = u.invoiceRepo.Update(invoice)

	return &CreatePaymentResponse{
		PaymentURL:  resp.Data.PaymentURL,
		Reference:   resp.Data.Reference,
		MerchantRef: resp.Data.MerchantRef,
		Amount:      resp.Data.Amount,
		ExpiredTime: resp.Data.ExpiredTime,
	}, nil
}

func (u *PaymentUsecase) GetPaymentGateways() ([]tripay.TripayPaymentChannel, error) {
	// Return defaults if Tripay not configured
	if !u.tripay.IsConfigured() {
		return []tripay.TripayPaymentChannel{
			{Code: "QRIS", Name: "QRIS", Group: "E-Money", Active: true},
			{Code: "BRIVA", Name: "BRI Virtual Account", Group: "Virtual Account", Active: true},
			{Code: "MANDIRIVA", Name: "Mandiri Virtual Account", Group: "Virtual Account", Active: true},
		}, nil
	}
	return u.tripay.GetPaymentChannels()
}

func (u *PaymentUsecase) HandleCallback(payload tripay.TripayCallbackPayload, rawBody, signature string) error {
	if !u.tripay.ValidateCallback(signature, rawBody) {
		return fmt.Errorf("invalid callback signature")
	}

	invoice, err := u.invoiceRepo.FindByNumber(payload.MerchantRef)
	if err != nil {
		return fmt.Errorf("invoice not found: %s", payload.MerchantRef)
	}

	if payload.Status == "PAID" {
		now := time.Now()
		invoice.Status = "paid"
		invoice.PaidAt = &now
		invoice.PaymentReference = payload.Reference

		if err := u.invoiceRepo.Update(invoice); err != nil {
			return fmt.Errorf("failed to update invoice: %w", err)
		}

		// Auto-activate customer after payment
		customer, err := u.customerRepo.FindByID(invoice.CustomerID)
		if err == nil && customer.Status == "isolated" {
			if err := u.mikrotikSvc.ActivateCustomer(customer); err != nil {
				logger.Warn("Failed to auto-activate customer after payment",
					zap.Uint("customer_id", customer.ID),
					zap.Error(err),
				)
			} else {
				logger.Info("Customer auto-activated after payment",
					zap.Uint("customer_id", customer.ID),
				)
			}
		}
	}

	return nil
}
