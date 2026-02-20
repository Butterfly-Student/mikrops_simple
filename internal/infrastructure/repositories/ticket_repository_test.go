//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestTicketRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	ticket := &entities.TroubleTicket{
		CustomerID:  customer.ID,
		Subject:     "Internet down",
		Description: "No connection since morning",
		Priority:    "high",
		Status:      "open",
	}

	err := repo.Create(ticket)
	assert.NoError(t, err)
	assert.NotZero(t, ticket.ID)

	found, err := repo.FindByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Internet down", found.Subject)
	assert.NotNil(t, found.Customer) // Preloaded
}

func TestTicketRepository_FindByCustomerID(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Jane", Phone: "082", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	_ = repo.Create(&entities.TroubleTicket{CustomerID: customer.ID, Subject: "T1", Description: "D1", Status: "open"})
	_ = repo.Create(&entities.TroubleTicket{CustomerID: customer.ID, Subject: "T2", Description: "D2", Status: "closed"})

	tickets, err := repo.FindByCustomerID(customer.ID)
	assert.NoError(t, err)
	assert.Len(t, tickets, 2)
}

func TestTicketRepository_FindAll_Paginated(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "083", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	for i := 0; i < 5; i++ {
		_ = repo.Create(&entities.TroubleTicket{
			CustomerID: customer.ID, Subject: "Ticket", Description: "Desc", Status: "open",
		})
	}

	tickets, total, err := repo.FindAll(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, tickets, 2)
}

func TestTicketRepository_FindByStatus(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "084", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	_ = repo.Create(&entities.TroubleTicket{CustomerID: customer.ID, Subject: "T1", Description: "D1", Status: "open"})
	_ = repo.Create(&entities.TroubleTicket{CustomerID: customer.ID, Subject: "T2", Description: "D2", Status: "closed"})
	_ = repo.Create(&entities.TroubleTicket{CustomerID: customer.ID, Subject: "T3", Description: "D3", Status: "open"})

	tickets, total, err := repo.FindByStatus("open", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, tickets, 2)
}

func TestTicketRepository_Update(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "085", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	ticket := &entities.TroubleTicket{CustomerID: customer.ID, Subject: "Old", Description: "Old desc", Status: "open"}
	_ = repo.Create(ticket)

	ticket.Status = "in_progress"
	ticket.AssignedTo = "tech-A"
	err := repo.Update(ticket)
	assert.NoError(t, err)

	found, _ := repo.FindByID(ticket.ID)
	assert.Equal(t, "in_progress", found.Status)
	assert.Equal(t, "tech-A", found.AssignedTo)
}

func TestTicketRepository_Delete(t *testing.T) {
	cleanTable(t, "trouble_tickets")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "086", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewTroubleTicketRepository(testDB)

	ticket := &entities.TroubleTicket{CustomerID: customer.ID, Subject: "Del", Description: "Del desc", Status: "open"}
	_ = repo.Create(ticket)

	err := repo.Delete(ticket.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ticket.ID)
	assert.Error(t, err)
}
