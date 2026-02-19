package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) repositories.InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *entities.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) FindByID(id uint) (*entities.Invoice, error) {
	var invoice entities.Invoice
	err := r.db.Preload("Customer").First(&invoice, id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) FindByNumber(number string) (*entities.Invoice, error) {
	var invoice entities.Invoice
	err := r.db.Preload("Customer").Where("number = ?", number).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) FindByCustomerID(customerID uint, page, perPage int) ([]*entities.Invoice, int64, error) {
	var invoices []*entities.Invoice
	var total int64

	query := r.db.Model(&entities.Invoice{}).Preload("Customer").Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&invoices).Error
	return invoices, total, err
}

func (r *invoiceRepository) Update(invoice *entities.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *invoiceRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Invoice{}, id).Error
}

func (r *invoiceRepository) FindAll(page, perPage int) ([]*entities.Invoice, int64, error) {
	var invoices []*entities.Invoice
	var total int64

	query := r.db.Model(&entities.Invoice{}).Preload("Customer")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&invoices).Error
	return invoices, total, err
}

func (r *invoiceRepository) FindByStatus(status string, page, perPage int) ([]*entities.Invoice, int64, error) {
	var invoices []*entities.Invoice
	var total int64

	query := r.db.Model(&entities.Invoice{}).Preload("Customer").Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&invoices).Error
	return invoices, total, err
}

func (r *invoiceRepository) FindLastInvoiceNumber() (*entities.Invoice, error) {
	var invoice entities.Invoice
	err := r.db.Order("id DESC").First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
