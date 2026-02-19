package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"gorm.io/gorm"
)

type routerRepository struct {
	db *gorm.DB
}

func NewRouterRepository(db *gorm.DB) repositories.RouterRepository {
	return &routerRepository{db: db}
}

func (r *routerRepository) Create(router *entities.Router) error {
	return r.db.Create(router).Error
}

func (r *routerRepository) FindByID(id uint) (*entities.Router, error) {
	var router entities.Router
	err := r.db.First(&router, id).Error
	if err != nil {
		return nil, err
	}
	return &router, nil
}

func (r *routerRepository) FindActive() (*entities.Router, error) {
	var router entities.Router
	err := r.db.Where("is_active = ?", true).First(&router).Error
	if err != nil {
		return nil, err
	}
	return &router, nil
}

func (r *routerRepository) FindAll() ([]*entities.Router, error) {
	var routers []*entities.Router
	err := r.db.Order("created_at DESC").Find(&routers).Error
	return routers, err
}

func (r *routerRepository) Update(router *entities.Router) error {
	return r.db.Save(router).Error
}

func (r *routerRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Router{}, id).Error
}
