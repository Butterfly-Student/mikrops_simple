package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) repositories.PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) Create(pkg *entities.Package) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) FindByID(id uint) (*entities.Package, error) {
	var pkg entities.Package
	err := r.db.First(&pkg, id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) FindByName(name string) (*entities.Package, error) {
	var pkg entities.Package
	err := r.db.Where("name = ?", name).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) Update(pkg *entities.Package) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Package{}, id).Error
}

func (r *packageRepository) FindAll() ([]*entities.Package, error) {
	var packages []*entities.Package
	err := r.db.Order("created_at DESC").Find(&packages).Error
	return packages, err
}
