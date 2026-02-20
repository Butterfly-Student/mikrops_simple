package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/genieacs"
)

// GenieACSUsecase depends on *genieacs.GenieACSClient (concrete type).
// We create a client pointing to a non-existent server to test error handling.
// The client will return errors when trying to make HTTP requests.

func newTestGenieACSUsecase() GenieACSUsecase {
	// Point to a non-routable address to ensure requests fail fast
	client := genieacs.NewGenieACSClient("http://127.0.0.1:1", "", "")
	return NewGenieACSUsecase(client)
}

func TestGenieACSGetDevices_ClientError(t *testing.T) {
	uc := newTestGenieACSUsecase()

	result, err := uc.GetDevices()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGenieACSGetDevice_ClientError(t *testing.T) {
	uc := newTestGenieACSUsecase()

	result, err := uc.GetDevice("TEST-SERIAL")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGenieACSRebootDevice_ClientError(t *testing.T) {
	uc := newTestGenieACSUsecase()

	err := uc.RebootDevice("TEST-SERIAL")

	assert.Error(t, err)
}

func TestGenieACSSetParameter_ClientError(t *testing.T) {
	uc := newTestGenieACSUsecase()

	err := uc.SetParameter("TEST-SERIAL", "some.param", "value")

	assert.Error(t, err)
}

func TestGenieACSFindDeviceByPPPoE_ClientError(t *testing.T) {
	uc := newTestGenieACSUsecase()

	result, err := uc.FindDeviceByPPPoE("user@pppoe")

	assert.Error(t, err)
	assert.Nil(t, result)
}
