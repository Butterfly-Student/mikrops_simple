package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MikroTikUsecase depends entirely on *mikrotik.MikroTikService (concrete type).
// We can only test the stub methods (GetRouter, CreateRouter, UpdateRouter, DeleteRouter, ActivateRouter, GetRouterStatus)
// which return nil without calling the service.
// Full integration tests for MikroTik operations require a running MikroTik router.

func newTestMikroTikUsecase() MikroTikUsecase {
	// All methods that delegate to mikrotikService will panic with nil pointer.
	// Only stub methods that return nil can be tested.
	return NewMikroTikUsecase(nil)
}

func TestMikroTikGetRouter_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	result, err := uc.GetRouter(1)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMikroTikCreateRouter_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	err := uc.CreateRouter(nil)

	assert.NoError(t, err)
}

func TestMikroTikUpdateRouter_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	err := uc.UpdateRouter(1, nil)

	assert.NoError(t, err)
}

func TestMikroTikDeleteRouter_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	err := uc.DeleteRouter(1)

	assert.NoError(t, err)
}

func TestMikroTikActivateRouter_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	err := uc.ActivateRouter(1)

	assert.NoError(t, err)
}

func TestMikroTikGetRouterStatus_Stub(t *testing.T) {
	uc := newTestMikroTikUsecase()

	result, err := uc.GetRouterStatus(1)

	assert.NoError(t, err)
	assert.Nil(t, result)
}
