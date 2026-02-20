package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/tripay"
)

func newTestPaymentUsecase(t *testing.T, tripayClient *tripay.TripayClient) (*PaymentUsecase, *mocks.InvoiceRepository, *mocks.CustomerRepository) {
	invoiceRepo := mocks.NewInvoiceRepository(t)
	customerRepo := mocks.NewCustomerRepository(t)
	// mikrotikSvc is nil - auto-activate after payment will panic if reached
	uc := NewPaymentUsecase(invoiceRepo, customerRepo, tripayClient, nil, "http://localhost:8080")
	return uc, invoiceRepo, customerRepo
}

func TestPaymentCreateTransaction_InvoiceNotFound(t *testing.T) {
	uc, invoiceRepo, _ := newTestPaymentUsecase(t, nil)

	invoiceRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	resp, err := uc.CreateTransaction(CreatePaymentRequest{InvoiceID: 999, PaymentMethod: "QRIS"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invoice not found")
}

func TestPaymentCreateTransaction_AlreadyPaid(t *testing.T) {
	uc, invoiceRepo, _ := newTestPaymentUsecase(t, nil)

	now := time.Now()
	invoiceRepo.On("FindByID", uint(1)).Return(&entities.Invoice{
		ID: 1, Status: "paid", DueDate: now, CreatedAt: now, UpdatedAt: now,
	}, nil)

	resp, err := uc.CreateTransaction(CreatePaymentRequest{InvoiceID: 1, PaymentMethod: "QRIS"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invoice is already paid")
}

func TestPaymentCreateTransaction_GatewayNotConfigured(t *testing.T) {
	tripayClient := tripay.NewTripayClient("", "", "", "sandbox")
	uc, invoiceRepo, _ := newTestPaymentUsecase(t, tripayClient)

	now := time.Now()
	invoiceRepo.On("FindByID", uint(1)).Return(&entities.Invoice{
		ID: 1, Status: "unpaid", CustomerID: 1, DueDate: now, CreatedAt: now, UpdatedAt: now,
	}, nil)

	resp, err := uc.CreateTransaction(CreatePaymentRequest{InvoiceID: 1, PaymentMethod: "QRIS"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "payment gateway not configured")
}

func TestPaymentGetGateways_Unconfigured(t *testing.T) {
	tripayClient := tripay.NewTripayClient("", "", "", "sandbox")
	uc, _, _ := newTestPaymentUsecase(t, tripayClient)

	channels, err := uc.GetPaymentGateways()

	assert.NoError(t, err)
	assert.Len(t, channels, 3) // Default channels
	assert.Equal(t, "QRIS", channels[0].Code)
}

func TestPaymentHandleCallback_InvalidSignature(t *testing.T) {
	tripayClient := tripay.NewTripayClient("key", "secret", "merchant", "sandbox")
	uc, _, _ := newTestPaymentUsecase(t, tripayClient)

	payload := tripay.TripayCallbackPayload{
		MerchantRef: "INV-000001",
		Status:      "PAID",
	}

	err := uc.HandleCallback(payload, "tampered-body", "invalid-signature")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid callback signature")
}

// computeCallbackSignature generates a valid HMAC-SHA256 signature for tripay callback testing
func computeCallbackSignature(privateKey, body string) string {
	mac := hmac.New(sha256.New, []byte(privateKey))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestPaymentHandleCallback_InvoiceNotFound(t *testing.T) {
	privateKey := "test-private-key"
	tripayClient := tripay.NewTripayClient("key", privateKey, "merchant", "sandbox")
	uc, invoiceRepo, _ := newTestPaymentUsecase(t, tripayClient)

	payload := tripay.TripayCallbackPayload{
		MerchantRef: "INV-NONEXIST",
		Status:      "PAID",
	}
	rawBody := `{"merchant_ref":"INV-NONEXIST","status":"PAID"}`
	signature := computeCallbackSignature(privateKey, rawBody)

	invoiceRepo.On("FindByNumber", "INV-NONEXIST").Return(nil, errors.New("not found"))

	err := uc.HandleCallback(payload, rawBody, signature)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invoice not found")
}

func TestPaymentHandleCallback_PaidSuccess(t *testing.T) {
	privateKey := "test-private-key"
	tripayClient := tripay.NewTripayClient("key", privateKey, "merchant", "sandbox")
	uc, invoiceRepo, customerRepo := newTestPaymentUsecase(t, tripayClient)

	payload := tripay.TripayCallbackPayload{
		MerchantRef: "INV-000001",
		Status:      "PAID",
		Reference:   "REF-123",
	}
	rawBody := `{"merchant_ref":"INV-000001","status":"PAID","reference":"REF-123"}`
	signature := computeCallbackSignature(privateKey, rawBody)

	now := time.Now()
	invoice := &entities.Invoice{
		ID: 1, Number: "INV-000001", CustomerID: 1, Status: "unpaid",
		Amount: 100000, DueDate: now, CreatedAt: now, UpdatedAt: now,
	}
	invoiceRepo.On("FindByNumber", "INV-000001").Return(invoice, nil)
	invoiceRepo.On("Update", mock.AnythingOfType("*entities.Invoice")).Return(nil)
	// Customer is active, so MikroTik auto-activate won't be called
	customerRepo.On("FindByID", uint(1)).Return(&entities.Customer{
		ID: 1, Status: "active",
	}, nil)

	err := uc.HandleCallback(payload, rawBody, signature)

	assert.NoError(t, err)
	assert.Equal(t, "paid", invoice.Status)
	assert.NotNil(t, invoice.PaidAt)
	assert.Equal(t, "REF-123", invoice.PaymentReference)
}

func TestPaymentHandleCallback_NonPaidStatus(t *testing.T) {
	privateKey := "test-private-key"
	tripayClient := tripay.NewTripayClient("key", privateKey, "merchant", "sandbox")
	uc, invoiceRepo, _ := newTestPaymentUsecase(t, tripayClient)

	payload := tripay.TripayCallbackPayload{
		MerchantRef: "INV-000001",
		Status:      "EXPIRED",
	}
	rawBody := `{"merchant_ref":"INV-000001","status":"EXPIRED"}`
	signature := computeCallbackSignature(privateKey, rawBody)

	now := time.Now()
	invoiceRepo.On("FindByNumber", "INV-000001").Return(&entities.Invoice{
		ID: 1, Number: "INV-000001", Status: "unpaid",
		DueDate: now, CreatedAt: now, UpdatedAt: now,
	}, nil)

	err := uc.HandleCallback(payload, rawBody, signature)

	// Non-PAID status doesn't update invoice, just returns nil
	assert.NoError(t, err)
}
