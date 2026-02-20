package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories/mocks"
)

func TestPackageGetAll_Success(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	expected := []*entities.Package{
		{ID: 1, Name: "Basic", Price: 100000, Speed: "10Mbps", Status: "active"},
		{ID: 2, Name: "Premium", Price: 250000, Speed: "50Mbps", Status: "active"},
	}
	mockRepo.On("FindAll").Return(expected, nil)

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Basic", result[0].Name)
}

func TestPackageGetByID_Success(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	expected := &entities.Package{ID: 1, Name: "Basic", Price: 100000, Status: "active"}
	mockRepo.On("FindByID", uint(1)).Return(expected, nil)

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Basic", result.Name)
}

func TestPackageGetByID_NotFound(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("record not found"))

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "package not found", err.Error())
}

func TestPackageCreate_Success(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("Create", mock.AnythingOfType("*entities.Package")).Return(nil)

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.Create(CreatePackageRequest{
		Name: "Basic", Price: 100000, Speed: "10Mbps", Status: "active",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Basic", result.Name)
	assert.Equal(t, "active", result.Status)
}

func TestPackageCreate_DefaultStatus(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("Create", mock.MatchedBy(func(pkg *entities.Package) bool {
		return pkg.Status == "active"
	})).Return(nil)

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.Create(CreatePackageRequest{Name: "Basic", Price: 100000})

	assert.NoError(t, err)
	assert.Equal(t, "active", result.Status)
}

func TestPackageUpdate_Success(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	existing := &entities.Package{
		ID: 1, Name: "Basic", Price: 100000, Speed: "10Mbps",
		Description: "Basic package", ProfileNormal: "pn", ProfileIsolir: "pi", Status: "active",
	}
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*entities.Package")).Return(nil)

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.Update(1, CreatePackageRequest{Name: "Premium", Price: 250000})

	assert.NoError(t, err)
	assert.Equal(t, "Premium", result.Name)
	assert.Equal(t, float64(250000), result.Price)
	assert.Equal(t, "Basic package", result.Description) // unchanged
}

func TestPackageUpdate_NotFound(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	uc := NewPackageUsecase(mockRepo)
	result, err := uc.Update(999, CreatePackageRequest{Name: "X"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "package not found", err.Error())
}

func TestPackageDelete_Success(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("FindByID", uint(1)).Return(&entities.Package{ID: 1}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	uc := NewPackageUsecase(mockRepo)
	err := uc.Delete(1)

	assert.NoError(t, err)
}

func TestPackageDelete_NotFound(t *testing.T) {
	mockRepo := mocks.NewPackageRepository(t)
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	uc := NewPackageUsecase(mockRepo)
	err := uc.Delete(999)

	assert.Error(t, err)
	assert.Equal(t, "package not found", err.Error())
}
