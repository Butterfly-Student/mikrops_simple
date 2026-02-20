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

func newTestInvoiceUsecase(t *testing.T) (InvoiceUsecase, *mocks.InvoiceRepository, *mocks.SettingRepository) {
	invoiceRepo := mocks.NewInvoiceRepository(t)
	settingRepo := mocks.NewSettingRepository(t)
	// whatsappService is nil - the code nil-checks before use
	uc := NewInvoiceUsecase(invoiceRepo, settingRepo, nil)
	return uc, invoiceRepo, settingRepo
}

func TestInvoiceGetAll_Success(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)

	now := time.Now()
	invoiceRepo.On("FindAll", 1, 20).Return([]*entities.Invoice{
		{ID: 1, Number: "INV-000001", Amount: 100000, Status: "unpaid", DueDate: now, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Number: "INV-000002", Amount: 250000, Status: "paid", DueDate: now, CreatedAt: now, UpdatedAt: now},
	}, int64(2), nil)

	result, err := uc.GetInvoices(1, 20)

	assert.NoError(t, err)
	assert.Len(t, result.Invoices, 2)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 1, result.TotalPages)
}

func TestInvoiceGetAll_DefaultPagination(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)

	now := time.Now()
	invoiceRepo.On("FindAll", 1, 20).Return([]*entities.Invoice{
		{ID: 1, Number: "INV-000001", DueDate: now, CreatedAt: now, UpdatedAt: now},
	}, int64(1), nil)

	// page=0 and perPage=0 should be corrected to defaults
	result, err := uc.GetInvoices(0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestInvoiceGetByID_Success(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)

	now := time.Now()
	invoiceRepo.On("FindByID", uint(1)).Return(&entities.Invoice{
		ID: 1, Number: "INV-000001", Amount: 100000, CustomerID: 1,
		DueDate: now, CreatedAt: now, UpdatedAt: now,
	}, nil)

	result, err := uc.GetInvoiceByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "INV-000001", result.Number)
}

func TestInvoiceGetByID_NotFound(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)
	invoiceRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := uc.GetInvoiceByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestInvoiceCreate_Success(t *testing.T) {
	uc, invoiceRepo, settingRepo := newTestInvoiceUsecase(t)

	settingRepo.On("Get", "INVOICE_PREFIX").Return("INV-", nil)
	invoiceRepo.On("FindLastInvoiceNumber").Return(nil, errors.New("no invoices"))
	settingRepo.On("Get", "INVOICE_START").Return("1", nil)
	invoiceRepo.On("Create", mock.AnythingOfType("*entities.Invoice")).Return(nil)

	err := uc.CreateInvoice(&dto.InvoiceDetail{
		CustomerID: 1, Amount: 100000, Period: "2024-01",
	})

	assert.NoError(t, err)
}

func TestInvoiceCreate_AutoNumbering(t *testing.T) {
	uc, invoiceRepo, settingRepo := newTestInvoiceUsecase(t)

	settingRepo.On("Get", "INVOICE_PREFIX").Return("INV-", nil)
	invoiceRepo.On("FindLastInvoiceNumber").Return(&entities.Invoice{Number: "INV-000005"}, nil)
	invoiceRepo.On("Create", mock.MatchedBy(func(inv *entities.Invoice) bool {
		return inv.Number == "INV-000006"
	})).Return(nil)

	err := uc.CreateInvoice(&dto.InvoiceDetail{
		CustomerID: 1, Amount: 200000, Period: "2024-02",
	})

	assert.NoError(t, err)
}

func TestInvoiceCreate_CustomDueDate(t *testing.T) {
	uc, invoiceRepo, settingRepo := newTestInvoiceUsecase(t)

	settingRepo.On("Get", "INVOICE_PREFIX").Return("INV-", nil)
	invoiceRepo.On("FindLastInvoiceNumber").Return(nil, errors.New("none"))
	settingRepo.On("Get", "INVOICE_START").Return("1", nil)
	invoiceRepo.On("Create", mock.MatchedBy(func(inv *entities.Invoice) bool {
		return inv.DueDate.Format("2006-01-02") == "2024-06-15"
	})).Return(nil)

	err := uc.CreateInvoice(&dto.InvoiceDetail{
		CustomerID: 1, Amount: 100000, Period: "2024-06", DueDate: "2024-06-15",
	})

	assert.NoError(t, err)
}

func TestInvoiceUpdate_Success(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)

	now := time.Now()
	existing := &entities.Invoice{
		ID: 1, Number: "INV-000001", Amount: 100000, Status: "unpaid",
		DueDate: now, CreatedAt: now, UpdatedAt: now,
	}
	invoiceRepo.On("FindByID", uint(1)).Return(existing, nil)
	invoiceRepo.On("Update", mock.AnythingOfType("*entities.Invoice")).Return(nil)

	err := uc.UpdateInvoice(1, &dto.InvoiceDetail{Amount: 150000, Period: "2024-02"})

	assert.NoError(t, err)
	assert.Equal(t, float64(150000), existing.Amount)
	assert.Equal(t, "2024-02", existing.Period)
}

func TestInvoiceUpdate_MarkPaid(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)

	now := time.Now()
	existing := &entities.Invoice{
		ID: 1, Number: "INV-000001", Amount: 100000, Status: "unpaid",
		DueDate: now, CreatedAt: now, UpdatedAt: now,
	}
	invoiceRepo.On("FindByID", uint(1)).Return(existing, nil)
	invoiceRepo.On("Update", mock.AnythingOfType("*entities.Invoice")).Return(nil)

	err := uc.UpdateInvoice(1, &dto.InvoiceDetail{Status: "paid"})

	assert.NoError(t, err)
	assert.Equal(t, "paid", existing.Status)
	assert.NotNil(t, existing.PaidAt)
}

func TestInvoiceUpdate_NotFound(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)
	invoiceRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := uc.UpdateInvoice(999, &dto.InvoiceDetail{Amount: 100})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invoice not found")
}

func TestInvoiceDelete_Success(t *testing.T) {
	uc, invoiceRepo, _ := newTestInvoiceUsecase(t)
	invoiceRepo.On("Delete", uint(1)).Return(nil)

	err := uc.DeleteInvoice(1)

	assert.NoError(t, err)
}
