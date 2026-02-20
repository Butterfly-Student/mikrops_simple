package whatsapp

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/gowa"
)

type IWhatsAppService interface {
	SendInvoiceNotification(invoice *entities.Invoice) error
	SendPaymentConfirmation(invoice *entities.Invoice) error
	SendIsolationNotification(customer *entities.Customer) error
	SendActivationNotification(customer *entities.Customer) error
	SendWelcomeMessage(customer *entities.Customer) error
	SendBulkNotification(message string, phones []string) error
	SetDeviceID(deviceID string)
	GetClient() gowa.GOWAClient
	GetCustomerByPhone(phone string) (*entities.Customer, error)
	IsAdmin(phone string) bool
}
