//go:build integration

package impl

import (
	"testing"
	"time"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	// Create customer for FK
	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	invoice := &entities.Invoice{
		CustomerID: customer.ID,
		Number:     "INV-000001",
		Amount:     100000,
		Period:     "2024-01",
		DueDate:    time.Now().Add(7 * 24 * time.Hour),
		Status:     "unpaid",
	}

	err := repo.Create(invoice)
	assert.NoError(t, err)
	assert.NotZero(t, invoice.ID)

	found, err := repo.FindByID(invoice.ID)
	assert.NoError(t, err)
	assert.Equal(t, "INV-000001", found.Number)
	assert.NotNil(t, found.Customer) // Preloaded
}

func TestInvoiceRepository_FindByNumber(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	_ = repo.Create(&entities.Invoice{
		CustomerID: customer.ID, Number: "INV-TEST", Amount: 50000,
		Period: "2024-01", DueDate: time.Now(), Status: "unpaid",
	})

	found, err := repo.FindByNumber("INV-TEST")
	assert.NoError(t, err)
	assert.Equal(t, "INV-TEST", found.Number)

	_, err = repo.FindByNumber("NONEXIST")
	assert.Error(t, err)
}

func TestInvoiceRepository_FindByCustomerID(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-001", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid"})
	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-002", Amount: 100000, Period: "2024-02", DueDate: time.Now(), Status: "paid"})

	invoices, total, err := repo.FindByCustomerID(customer.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, invoices, 2)
}

func TestInvoiceRepository_FindAll_Paginated(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	for i := 0; i < 5; i++ {
		_ = repo.Create(&entities.Invoice{
			CustomerID: customer.ID, Number: "INV-P" + string(rune('0'+i)),
			Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid",
		})
	}

	invoices, total, err := repo.FindAll(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, invoices, 2)
}

func TestInvoiceRepository_FindByStatus(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-U1", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid"})
	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-P1", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "paid"})

	invoices, total, err := repo.FindByStatus("unpaid", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, invoices, 1)
}

func TestInvoiceRepository_FindLastInvoiceNumber(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-001", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid"})
	_ = repo.Create(&entities.Invoice{CustomerID: customer.ID, Number: "INV-002", Amount: 200000, Period: "2024-02", DueDate: time.Now(), Status: "unpaid"})

	last, err := repo.FindLastInvoiceNumber()
	assert.NoError(t, err)
	assert.Equal(t, "INV-002", last.Number)
}

func TestInvoiceRepository_Update(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	invoice := &entities.Invoice{CustomerID: customer.ID, Number: "INV-UPD", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid"}
	_ = repo.Create(invoice)

	invoice.Status = "paid"
	now := time.Now()
	invoice.PaidAt = &now
	err := repo.Update(invoice)
	assert.NoError(t, err)

	found, _ := repo.FindByID(invoice.ID)
	assert.Equal(t, "paid", found.Status)
	assert.NotNil(t, found.PaidAt)
}

func TestInvoiceRepository_Delete(t *testing.T) {
	cleanTable(t, "invoices")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewInvoiceRepository(testDB)

	invoice := &entities.Invoice{CustomerID: customer.ID, Number: "INV-DEL", Amount: 100000, Period: "2024-01", DueDate: time.Now(), Status: "unpaid"}
	_ = repo.Create(invoice)

	err := repo.Delete(invoice.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(invoice.ID)
	assert.Error(t, err)
}
