//go:build integration

package impl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingRepository_GetAndSet(t *testing.T) {
	cleanTable(t, "settings")

	repo := NewSettingRepository(testDB)

	// Get non-existent key
	_, err := repo.Get("NON_EXISTENT")
	assert.Error(t, err)

	// Set a new key
	err = repo.Set("INVOICE_PREFIX", "INV-")
	assert.NoError(t, err)

	// Get the key
	val, err := repo.Get("INVOICE_PREFIX")
	assert.NoError(t, err)
	assert.Equal(t, "INV-", val)

	// Update existing key (upsert)
	err = repo.Set("INVOICE_PREFIX", "BILL-")
	assert.NoError(t, err)

	val, err = repo.Get("INVOICE_PREFIX")
	assert.NoError(t, err)
	assert.Equal(t, "BILL-", val)
}

func TestSettingRepository_GetAll(t *testing.T) {
	cleanTable(t, "settings")

	repo := NewSettingRepository(testDB)

	_ = repo.Set("KEY1", "val1")
	_ = repo.Set("KEY2", "val2")
	_ = repo.Set("KEY3", "val3")

	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 3)
	assert.Equal(t, "val1", all["KEY1"])
	assert.Equal(t, "val2", all["KEY2"])
	assert.Equal(t, "val3", all["KEY3"])
}
