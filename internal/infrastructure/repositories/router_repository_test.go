//go:build integration

package impl

import (
	"testing"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestRouterRepository_CreateAndFindByID(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	router := &entities.Router{
		Name: "Main Router", Host: "10.0.0.1", Username: "admin",
		Password: "secret", Port: 8728, IsActive: true,
	}

	err := repo.Create(router)
	assert.NoError(t, err)
	assert.NotZero(t, router.ID)

	found, err := repo.FindByID(router.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Main Router", found.Name)
	assert.Equal(t, "10.0.0.1", found.Host)
}

func TestRouterRepository_FindActive(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	_ = repo.Create(&entities.Router{Name: "R1", Host: "10.0.0.1", Username: "a", Password: "p", Port: 8728, IsActive: false})
	_ = repo.Create(&entities.Router{Name: "R2", Host: "10.0.0.2", Username: "a", Password: "p", Port: 8728, IsActive: true})

	found, err := repo.FindActive()
	assert.NoError(t, err)
	assert.Equal(t, "R2", found.Name)
	assert.True(t, found.IsActive)
}

func TestRouterRepository_FindActive_None(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	_ = repo.Create(&entities.Router{Name: "R1", Host: "10.0.0.1", Username: "a", Password: "p", Port: 8728, IsActive: false})

	_, err := repo.FindActive()
	assert.Error(t, err)
}

func TestRouterRepository_FindAll(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	_ = repo.Create(&entities.Router{Name: "R1", Host: "10.0.0.1", Username: "a", Password: "p", Port: 8728})
	_ = repo.Create(&entities.Router{Name: "R2", Host: "10.0.0.2", Username: "a", Password: "p", Port: 8728})

	routers, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, routers, 2)
}

func TestRouterRepository_Update(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	router := &entities.Router{Name: "Old", Host: "10.0.0.1", Username: "a", Password: "p", Port: 8728}
	_ = repo.Create(router)

	router.Name = "Updated"
	router.Host = "10.0.0.99"
	err := repo.Update(router)
	assert.NoError(t, err)

	found, _ := repo.FindByID(router.ID)
	assert.Equal(t, "Updated", found.Name)
	assert.Equal(t, "10.0.0.99", found.Host)
}

func TestRouterRepository_Delete(t *testing.T) {
	cleanTable(t, "routers")

	repo := NewRouterRepository(testDB)

	router := &entities.Router{Name: "Del", Host: "10.0.0.1", Username: "a", Password: "p", Port: 8728}
	_ = repo.Create(router)

	err := repo.Delete(router.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(router.ID)
	assert.Error(t, err)
}
