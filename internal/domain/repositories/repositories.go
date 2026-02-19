package repositories

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
)

type AdminRepository interface {
	Create(admin *entities.AdminUser) error
	FindByID(id uint) (*entities.AdminUser, error)
	FindByUsername(username string) (*entities.AdminUser, error)
	Update(admin *entities.AdminUser) error
	Delete(id uint) error
	FindAll() ([]*entities.AdminUser, error)
}

type CustomerRepository interface {
	Create(customer *entities.Customer) error
	FindByID(id uint) (*entities.Customer, error)
	FindByPhone(phone string) (*entities.Customer, error)
	FindByPPPoEUsername(username string) (*entities.Customer, error)
	Update(customer *entities.Customer) error
	Delete(id uint) error
	FindAll(page, perPage int, search string) ([]*entities.Customer, int64, error)
	FindByStatus(status string, page, perPage int) ([]*entities.Customer, int64, error)
	FindByPackageID(packageID uint, page, perPage int) ([]*entities.Customer, int64, error)
}

type PackageRepository interface {
	Create(pkg *entities.Package) error
	FindByID(id uint) (*entities.Package, error)
	FindByName(name string) (*entities.Package, error)
	Update(pkg *entities.Package) error
	Delete(id uint) error
	FindAll() ([]*entities.Package, error)
}

type InvoiceRepository interface {
	Create(invoice *entities.Invoice) error
	FindByID(id uint) (*entities.Invoice, error)
	FindByNumber(number string) (*entities.Invoice, error)
	FindByCustomerID(customerID uint, page, perPage int) ([]*entities.Invoice, int64, error)
	FindLastInvoiceNumber() (*entities.Invoice, error)
	Update(invoice *entities.Invoice) error
	Delete(id uint) error
	FindAll(page, perPage int) ([]*entities.Invoice, int64, error)
	FindByStatus(status string, page, perPage int) ([]*entities.Invoice, int64, error)
}

type RouterRepository interface {
	Create(router *entities.Router) error
	FindByID(id uint) (*entities.Router, error)
	FindActive() (*entities.Router, error)
	FindAll() ([]*entities.Router, error)
	Update(router *entities.Router) error
	Delete(id uint) error
}

type ONULocationRepository interface {
	Create(location *entities.ONULocation) error
	FindByID(id uint) (*entities.ONULocation, error)
	FindByCustomerID(customerID uint) (*entities.ONULocation, error)
	FindByONUID(onuID string) (*entities.ONULocation, error)
	FindBySerialNumber(serial string) (*entities.ONULocation, error)
	Update(location *entities.ONULocation) error
	Delete(id uint) error
	FindAll(page, perPage int) ([]*entities.ONULocation, int64, error)
}

type TroubleTicketRepository interface {
	Create(ticket *entities.TroubleTicket) error
	FindByID(id uint) (*entities.TroubleTicket, error)
	FindByCustomerID(customerID uint) ([]*entities.TroubleTicket, error)
	Update(ticket *entities.TroubleTicket) error
	Delete(id uint) error
	FindAll(page, perPage int) ([]*entities.TroubleTicket, int64, error)
	FindByStatus(status string, page, perPage int) ([]*entities.TroubleTicket, int64, error)
}

type SettingRepository interface {
	Get(key string) (string, error)
	Set(key, value string) error
	GetAll() (map[string]string, error)
}

type CronScheduleRepository interface {
	Create(schedule *entities.CronSchedule) error
	FindByID(id uint) (*entities.CronSchedule, error)
	FindByTaskType(taskType string) (*entities.CronSchedule, error)
	FindActive() ([]*entities.CronSchedule, error)
	Update(schedule *entities.CronSchedule) error
	Delete(id uint) error
}
