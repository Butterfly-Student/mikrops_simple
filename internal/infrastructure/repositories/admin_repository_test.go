//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "admin_users")

	repo := NewAdminRepository(testDB)

	admin := &entities.AdminUser{
		Username: "testadmin",
		Password: "hashed-password",
		Email:    "admin@test.com",
		Role:     "admin",
		Status:   "active",
	}

	err := repo.Create(admin)
	assert.NoError(t, err)
	assert.NotZero(t, admin.ID)

	found, err := repo.FindByID(admin.ID)
	assert.NoError(t, err)
	assert.Equal(t, "testadmin", found.Username)
	assert.Equal(t, "admin@test.com", found.Email)
}

func TestAdminRepository_FindByUsername(t *testing.T) {
	cleanTable(t, "admin_users")

	repo := NewAdminRepository(testDB)

	admin := &entities.AdminUser{
		Username: "findme",
		Password: "pass",
		Role:     "admin",
		Status:   "active",
	}
	_ = repo.Create(admin)

	found, err := repo.FindByUsername("findme")
	assert.NoError(t, err)
	assert.Equal(t, "findme", found.Username)

	_, err = repo.FindByUsername("nonexistent")
	assert.Error(t, err)
}

func TestAdminRepository_Update(t *testing.T) {
	cleanTable(t, "admin_users")

	repo := NewAdminRepository(testDB)

	admin := &entities.AdminUser{
		Username: "updateme",
		Password: "pass",
		Role:     "admin",
		Status:   "active",
	}
	_ = repo.Create(admin)

	admin.Email = "updated@test.com"
	err := repo.Update(admin)
	assert.NoError(t, err)

	found, _ := repo.FindByID(admin.ID)
	assert.Equal(t, "updated@test.com", found.Email)
}

func TestAdminRepository_Delete(t *testing.T) {
	cleanTable(t, "admin_users")

	repo := NewAdminRepository(testDB)

	admin := &entities.AdminUser{
		Username: "deleteme",
		Password: "pass",
		Role:     "admin",
		Status:   "active",
	}
	_ = repo.Create(admin)

	err := repo.Delete(admin.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(admin.ID)
	assert.Error(t, err)
}

func TestAdminRepository_FindAll(t *testing.T) {
	cleanTable(t, "admin_users")

	repo := NewAdminRepository(testDB)

	_ = repo.Create(&entities.AdminUser{Username: "admin1", Password: "p", Role: "admin", Status: "active"})
	_ = repo.Create(&entities.AdminUser{Username: "admin2", Password: "p", Role: "operator", Status: "active"})

	admins, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, admins, 2)
}
