//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestPackageRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "packages")

	repo := NewPackageRepository(testDB)

	pkg := &entities.Package{
		Name: "Basic", Price: 100000, Speed: "10Mbps",
		ProfileNormal: "10M", ProfileIsolir: "256K", Status: "active",
	}

	err := repo.Create(pkg)
	assert.NoError(t, err)
	assert.NotZero(t, pkg.ID)

	found, err := repo.FindByID(pkg.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Basic", found.Name)
	assert.Equal(t, float64(100000), found.Price)
}

func TestPackageRepository_FindByName(t *testing.T) {
	cleanTable(t, "packages")

	repo := NewPackageRepository(testDB)

	_ = repo.Create(&entities.Package{Name: "Premium", Price: 250000, Status: "active"})

	found, err := repo.FindByName("Premium")
	assert.NoError(t, err)
	assert.Equal(t, "Premium", found.Name)

	_, err = repo.FindByName("Nonexistent")
	assert.Error(t, err)
}

func TestPackageRepository_FindAll(t *testing.T) {
	cleanTable(t, "packages")

	repo := NewPackageRepository(testDB)

	_ = repo.Create(&entities.Package{Name: "Basic", Price: 100000, Status: "active"})
	_ = repo.Create(&entities.Package{Name: "Premium", Price: 250000, Status: "active"})

	packages, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, packages, 2)
}

func TestPackageRepository_Update(t *testing.T) {
	cleanTable(t, "packages")

	repo := NewPackageRepository(testDB)

	pkg := &entities.Package{Name: "Old", Price: 100000, Status: "active"}
	_ = repo.Create(pkg)

	pkg.Name = "Updated"
	pkg.Price = 150000
	err := repo.Update(pkg)
	assert.NoError(t, err)

	found, _ := repo.FindByID(pkg.ID)
	assert.Equal(t, "Updated", found.Name)
	assert.Equal(t, float64(150000), found.Price)
}

func TestPackageRepository_Delete(t *testing.T) {
	cleanTable(t, "packages")

	repo := NewPackageRepository(testDB)

	pkg := &entities.Package{Name: "ToDelete", Price: 100000, Status: "active"}
	_ = repo.Create(pkg)

	err := repo.Delete(pkg.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(pkg.ID)
	assert.Error(t, err)
}
