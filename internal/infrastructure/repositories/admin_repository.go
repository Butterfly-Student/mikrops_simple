package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepositoryImpl {
	return &AdminRepositoryImpl{db: db}
}

func (r *AdminRepositoryImpl) Create(admin *entities.AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepositoryImpl) FindByID(id uint) (*entities.AdminUser, error) {
	var admin entities.AdminUser
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepositoryImpl) FindByUsername(username string) (*entities.AdminUser, error) {
	var admin entities.AdminUser
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepositoryImpl) Update(admin *entities.AdminUser) error {
	return r.db.Save(admin).Error
}

func (r *AdminRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.AdminUser{}, id).Error
}

func (r *AdminRepositoryImpl) FindAll() ([]*entities.AdminUser, error) {
	var admins []*entities.AdminUser
	if err := r.db.Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}
