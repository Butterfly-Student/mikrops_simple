//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "customers")
	cleanTable(t, "packages")

	// Create a package first for FK
	pkg := &entities.Package{Name: "Basic", Price: 100000, Status: "active"}
	testDB.Create(pkg)

	repo := NewCustomerRepository(testDB)

	customer := &entities.Customer{
		Name:          "John Doe",
		Phone:         "08123456789",
		Email:         "john@test.com",
		PackageID:     pkg.ID,
		PPPoEUsername: "john-pppoe",
		PPPoEPassword: "secret",
		Status:        "active",
	}

	err := repo.Create(customer)
	assert.NoError(t, err)
	assert.NotZero(t, customer.ID)

	found, err := repo.FindByID(customer.ID)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", found.Name)
	assert.NotNil(t, found.Package) // Preloaded
	assert.Equal(t, "Basic", found.Package.Name)
}

func TestCustomerRepository_FindByPhone(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	customer := &entities.Customer{
		Name: "Jane", Phone: "08111222333", PPPoEPassword: "p", Status: "active",
	}
	_ = repo.Create(customer)

	found, err := repo.FindByPhone("08111222333")
	assert.NoError(t, err)
	assert.Equal(t, "Jane", found.Name)

	_, err = repo.FindByPhone("00000000000")
	assert.Error(t, err)
}

func TestCustomerRepository_FindByPPPoEUsername(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	customer := &entities.Customer{
		Name: "Bob", Phone: "08999", PPPoEUsername: "bob-pppoe", PPPoEPassword: "p", Status: "active",
	}
	_ = repo.Create(customer)

	found, err := repo.FindByPPPoEUsername("bob-pppoe")
	assert.NoError(t, err)
	assert.Equal(t, "Bob", found.Name)
}

func TestCustomerRepository_FindAll_Paginated(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	for i := 0; i < 5; i++ {
		_ = repo.Create(&entities.Customer{
			Name: "Customer", Phone: "0800000000" + string(rune('0'+i)),
			PPPoEUsername: "pppoe-" + string(rune('0'+i)),
			PPPoEPassword: "p", Status: "active",
		})
	}

	customers, total, err := repo.FindAll(1, 2, "")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, customers, 2)
}

func TestCustomerRepository_FindAll_WithSearch(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	_ = repo.Create(&entities.Customer{Name: "Alpha", Phone: "081", PPPoEUsername: "alpha", PPPoEPassword: "p", Status: "active"})
	_ = repo.Create(&entities.Customer{Name: "Beta", Phone: "082", PPPoEUsername: "beta", PPPoEPassword: "p", Status: "active"})

	customers, total, err := repo.FindAll(1, 10, "Alpha")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, customers, 1)
	assert.Equal(t, "Alpha", customers[0].Name)
}

func TestCustomerRepository_FindByStatus(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	_ = repo.Create(&entities.Customer{Name: "Active1", Phone: "081", PPPoEUsername: "active1", PPPoEPassword: "p", Status: "active"})
	_ = repo.Create(&entities.Customer{Name: "Isolated1", Phone: "082", PPPoEUsername: "isolated1", PPPoEPassword: "p", Status: "isolated"})
	_ = repo.Create(&entities.Customer{Name: "Active2", Phone: "083", PPPoEUsername: "active2", PPPoEPassword: "p", Status: "active"})

	customers, total, err := repo.FindByStatus("active", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, customers, 2)
}

func TestCustomerRepository_Update(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	customer := &entities.Customer{Name: "Old", Phone: "081", PPPoEPassword: "p", Status: "active"}
	_ = repo.Create(customer)

	customer.Name = "New"
	err := repo.Update(customer)
	assert.NoError(t, err)

	found, _ := repo.FindByID(customer.ID)
	assert.Equal(t, "New", found.Name)
}

func TestCustomerRepository_Delete(t *testing.T) {
	cleanTable(t, "customers")

	repo := NewCustomerRepository(testDB)

	customer := &entities.Customer{Name: "Delete", Phone: "081", PPPoEPassword: "p", Status: "active"}
	_ = repo.Create(customer)

	err := repo.Delete(customer.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(customer.ID)
	assert.Error(t, err)
}
