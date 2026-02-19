package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) repositories.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *entities.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) FindByID(id uint) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.db.Preload("Package").First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByPhone(phone string) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.db.Preload("Package").Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByPPPoEUsername(username string) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.db.Preload("Package").Where("pppoe_username = ?", username).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(customer *entities.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Customer{}, id).Error
}

func (r *customerRepository) FindAll(page, perPage int, search string) ([]*entities.Customer, int64, error) {
	var customers []*entities.Customer
	var total int64

	query := r.db.Model(&entities.Customer{}).Preload("Package")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR pppoe_username LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&customers).Error
	return customers, total, err
}

func (r *customerRepository) FindByStatus(status string, page, perPage int) ([]*entities.Customer, int64, error) {
	var customers []*entities.Customer
	var total int64

	query := r.db.Model(&entities.Customer{}).Preload("Package").Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&customers).Error
	return customers, total, err
}

func (r *customerRepository) FindByPackageID(packageID uint, page, perPage int) ([]*entities.Customer, int64, error) {
	var customers []*entities.Customer
	var total int64

	query := r.db.Model(&entities.Customer{}).Preload("Package").Where("package_id = ?", packageID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&customers).Error
	return customers, total, err
}
