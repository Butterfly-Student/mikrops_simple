package usecase

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
)

type PackageUsecase struct {
	packageRepo repositories.PackageRepository
}

func NewPackageUsecase(packageRepo repositories.PackageRepository) *PackageUsecase {
	return &PackageUsecase{packageRepo: packageRepo}
}

// GetAll returns all packages.
func (u *PackageUsecase) GetAll() ([]*entities.Package, error) {
	return u.packageRepo.FindAll()
}

// GetByID returns a single package.
func (u *PackageUsecase) GetByID(id uint) (*entities.Package, error) {
	pkg, err := u.packageRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found")
	}
	return pkg, nil
}

// CreatePackageRequest holds the fields for creating/updating a package.
type CreatePackageRequest struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	Speed         string  `json:"speed"`
	Description   string  `json:"description"`
	ProfileNormal string  `json:"profile_normal"`
	ProfileIsolir string  `json:"profile_isolir"`
	Status        string  `json:"status"`
}

// Create creates a new internet package.
func (u *PackageUsecase) Create(req CreatePackageRequest) (*entities.Package, error) {
	if req.Status == "" {
		req.Status = "active"
	}

	pkg := &entities.Package{
		Name:          req.Name,
		Price:         req.Price,
		Speed:         req.Speed,
		Description:   req.Description,
		ProfileNormal: req.ProfileNormal,
		ProfileIsolir: req.ProfileIsolir,
		Status:        req.Status,
	}

	if err := u.packageRepo.Create(pkg); err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}
	return pkg, nil
}

// Update updates an existing package.
func (u *PackageUsecase) Update(id uint, req CreatePackageRequest) (*entities.Package, error) {
	pkg, err := u.packageRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("package not found")
	}

	if req.Name != "" {
		pkg.Name = req.Name
	}
	if req.Price > 0 {
		pkg.Price = req.Price
	}
	if req.Speed != "" {
		pkg.Speed = req.Speed
	}
	if req.Description != "" {
		pkg.Description = req.Description
	}
	if req.ProfileNormal != "" {
		pkg.ProfileNormal = req.ProfileNormal
	}
	if req.ProfileIsolir != "" {
		pkg.ProfileIsolir = req.ProfileIsolir
	}
	if req.Status != "" {
		pkg.Status = req.Status
	}

	if err := u.packageRepo.Update(pkg); err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}
	return pkg, nil
}

// Delete removes a package.
func (u *PackageUsecase) Delete(id uint) error {
	if _, err := u.packageRepo.FindByID(id); err != nil {
		return fmt.Errorf("package not found")
	}
	return u.packageRepo.Delete(id)
}
