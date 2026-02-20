//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestONURepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "John", Phone: "081", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	loc := &entities.ONULocation{
		CustomerID:   customer.ID,
		ONUID:        "ONU-001",
		SerialNumber: "SN-001",
		Latitude:     -6.2,
		Longitude:    106.8,
		Address:      "Jakarta",
	}

	err := repo.Create(loc)
	assert.NoError(t, err)
	assert.NotZero(t, loc.ID)

	found, err := repo.FindByID(loc.ID)
	assert.NoError(t, err)
	assert.Equal(t, "ONU-001", found.ONUID)
	assert.NotNil(t, found.Customer) // Preloaded
}

func TestONURepository_FindByCustomerID(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Jane", Phone: "082", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	_ = repo.Create(&entities.ONULocation{CustomerID: customer.ID, ONUID: "ONU-C1", SerialNumber: "SN-C1"})

	found, err := repo.FindByCustomerID(customer.ID)
	assert.NoError(t, err)
	assert.Equal(t, "ONU-C1", found.ONUID)
}

func TestONURepository_FindBySerialNumber(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Bob", Phone: "083", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	_ = repo.Create(&entities.ONULocation{CustomerID: customer.ID, ONUID: "ONU-SN", SerialNumber: "FIND-ME-SN"})

	found, err := repo.FindBySerialNumber("FIND-ME-SN")
	assert.NoError(t, err)
	assert.Equal(t, "ONU-SN", found.ONUID)

	_, err = repo.FindBySerialNumber("NONEXIST")
	assert.Error(t, err)
}

func TestONURepository_FindAll_Paginated(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "084", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	for i := 0; i < 5; i++ {
		_ = repo.Create(&entities.ONULocation{
			CustomerID: customer.ID, ONUID: "ONU-" + string(rune('A'+i)),
			SerialNumber: "SN-" + string(rune('A'+i)),
		})
	}

	locs, total, err := repo.FindAll(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, locs, 2)
}

func TestONURepository_Update(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "085", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	loc := &entities.ONULocation{CustomerID: customer.ID, ONUID: "ONU-UPD", SerialNumber: "SN-UPD"}
	_ = repo.Create(loc)

	loc.Address = "Updated Address"
	err := repo.Update(loc)
	assert.NoError(t, err)

	found, _ := repo.FindByID(loc.ID)
	assert.Equal(t, "Updated Address", found.Address)
}

func TestONURepository_Delete(t *testing.T) {
	cleanTable(t, "onu_locations")
	cleanTable(t, "customers")

	customer := &entities.Customer{Name: "Test", Phone: "086", PPPoEPassword: "p", Status: "active"}
	testDB.Create(customer)

	repo := NewONULocationRepository(testDB)

	loc := &entities.ONULocation{CustomerID: customer.ID, ONUID: "ONU-DEL", SerialNumber: "SN-DEL"}
	_ = repo.Create(loc)

	err := repo.Delete(loc.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(loc.ID)
	assert.Error(t, err)
}
