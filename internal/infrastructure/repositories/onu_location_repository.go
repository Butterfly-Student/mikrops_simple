package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"gorm.io/gorm"
)

type ONULocationRepositoryImpl struct {
	db *gorm.DB
}

func NewONULocationRepository(db *gorm.DB) *ONULocationRepositoryImpl {
	return &ONULocationRepositoryImpl{db: db}
}

func (r *ONULocationRepositoryImpl) Create(loc *entities.ONULocation) error {
	return r.db.Create(loc).Error
}

func (r *ONULocationRepositoryImpl) FindByID(id uint) (*entities.ONULocation, error) {
	var loc entities.ONULocation
	if err := r.db.Preload("Customer").First(&loc, id).Error; err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *ONULocationRepositoryImpl) FindByCustomerID(customerID uint) (*entities.ONULocation, error) {
	var loc entities.ONULocation
	if err := r.db.Where("customer_id = ?", customerID).First(&loc).Error; err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *ONULocationRepositoryImpl) FindByONUID(onuID string) (*entities.ONULocation, error) {
	var loc entities.ONULocation
	if err := r.db.Where("onu_id = ?", onuID).First(&loc).Error; err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *ONULocationRepositoryImpl) FindBySerialNumber(serial string) (*entities.ONULocation, error) {
	var loc entities.ONULocation
	if err := r.db.Where("serial_number = ?", serial).First(&loc).Error; err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *ONULocationRepositoryImpl) Update(loc *entities.ONULocation) error {
	return r.db.Save(loc).Error
}

func (r *ONULocationRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.ONULocation{}, id).Error
}

func (r *ONULocationRepositoryImpl) FindAll(page, perPage int) ([]*entities.ONULocation, int64, error) {
	var locs []*entities.ONULocation
	var total int64

	r.db.Model(&entities.ONULocation{}).Count(&total)

	offset := (page - 1) * perPage
	if err := r.db.Preload("Customer").Offset(offset).Limit(perPage).Find(&locs).Error; err != nil {
		return nil, 0, err
	}
	return locs, total, nil
}
