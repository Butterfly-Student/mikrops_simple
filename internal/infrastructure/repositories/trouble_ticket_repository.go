package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"gorm.io/gorm"
)

type TroubleTicketRepositoryImpl struct {
	db *gorm.DB
}

func NewTroubleTicketRepository(db *gorm.DB) *TroubleTicketRepositoryImpl {
	return &TroubleTicketRepositoryImpl{db: db}
}

func (r *TroubleTicketRepositoryImpl) Create(ticket *entities.TroubleTicket) error {
	return r.db.Create(ticket).Error
}

func (r *TroubleTicketRepositoryImpl) FindByID(id uint) (*entities.TroubleTicket, error) {
	var ticket entities.TroubleTicket
	if err := r.db.Preload("Customer").First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TroubleTicketRepositoryImpl) FindByCustomerID(customerID uint) ([]*entities.TroubleTicket, error) {
	var tickets []*entities.TroubleTicket
	if err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TroubleTicketRepositoryImpl) Update(ticket *entities.TroubleTicket) error {
	return r.db.Save(ticket).Error
}

func (r *TroubleTicketRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.TroubleTicket{}, id).Error
}

func (r *TroubleTicketRepositoryImpl) FindAll(page, perPage int) ([]*entities.TroubleTicket, int64, error) {
	var tickets []*entities.TroubleTicket
	var total int64

	r.db.Model(&entities.TroubleTicket{}).Count(&total)

	offset := (page - 1) * perPage
	if err := r.db.Preload("Customer").Order("created_at DESC").Offset(offset).Limit(perPage).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *TroubleTicketRepositoryImpl) FindByStatus(status string, page, perPage int) ([]*entities.TroubleTicket, int64, error) {
	var tickets []*entities.TroubleTicket
	var total int64

	r.db.Model(&entities.TroubleTicket{}).Where("status = ?", status).Count(&total)

	offset := (page - 1) * perPage
	if err := r.db.Preload("Customer").Where("status = ?", status).Order("created_at DESC").Offset(offset).Limit(perPage).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}
