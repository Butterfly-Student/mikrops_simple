package whatsapp

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/gowa"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

type WhatsAppService struct {
	client       gowa.GOWAClient
	customerRepo repositories.CustomerRepository
	invoiceRepo  repositories.InvoiceRepository
	adminPhones  []string
}

func NewWhatsAppService(
	client gowa.GOWAClient,
	customerRepo repositories.CustomerRepository,
	invoiceRepo repositories.InvoiceRepository,
	adminPhones []string,
) *WhatsAppService {
	return &WhatsAppService{
		client:       client,
		customerRepo: customerRepo,
		invoiceRepo:  invoiceRepo,
		adminPhones:  adminPhones,
	}
}

func (s *WhatsAppService) SendInvoiceNotification(invoice *entities.Invoice) error {
	customer, err := s.customerRepo.FindByID(invoice.CustomerID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf(`*Invoice Baru Dibuat* 🧾

No: %s
Pelanggan: %s
Jumlah: Rp %.2f
Periode: %s
Jatuh Tempo: %s

Silakan lakukan pembayaran sebelum jatuh tempo. Terima kasih!`,
		invoice.Number,
		customer.Name,
		invoice.Amount,
		invoice.Period,
		invoice.DueDate.Format("2006-01-02"),
	)

	return s.client.SendText(customer.Phone, message)
}

func (s *WhatsAppService) SendPaymentConfirmation(invoice *entities.Invoice) error {
	customer, err := s.customerRepo.FindByID(invoice.CustomerID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf(`*Pembayaran Berhasil* ✅

No Invoice: %s
Pelanggan: %s
Jumlah: Rp %.2f
Metode: %s

Terima kasih atas pembayaran Anda! Koneksi Anda kini aktif.`,
		invoice.Number,
		customer.Name,
		invoice.Amount,
		invoice.PaymentMethod,
	)

	return s.client.SendText(customer.Phone, message)
}

func (s *WhatsAppService) SendIsolationNotification(customer *entities.Customer) error {
	message := fmt.Sprintf(`*Akun Diisolir* ⚠️

Pelanggan: %s
Status: %s

Maaf, koneksi internet Anda telah diisolir karena belum melakukan pembayaran.
Silakan hubungi admin atau lakukan pembayaran untuk mengaktifkan kembali.`,
		customer.Name,
		customer.Status,
	)

	return s.client.SendText(customer.Phone, message)
}

func (s *WhatsAppService) SendActivationNotification(customer *entities.Customer) error {
	message := fmt.Sprintf(`*Akun Diaktifkan* ✅

Pelanggan: %s
Status: %s

Koneksi internet Anda telah diaktifkan kembali.
Terima kasih telah melakukan pembayaran!`,
		customer.Name,
		customer.Status,
	)

	return s.client.SendText(customer.Phone, message)
}

func (s *WhatsAppService) SendWelcomeMessage(customer *entities.Customer) error {
	pkgName := ""
	if customer.Package != nil {
		pkgName = customer.Package.Name
	}

	message := fmt.Sprintf(`*Selamat Datang!* 👋

Nama: %s
Paket: %s
Username: %s

Terima kasih telah berlangganan! 
Hubungi kami jika ada kendala.`,
		customer.Name,
		pkgName,
		customer.PPPoEUsername,
	)

	return s.client.SendText(customer.Phone, message)
}

func (s *WhatsAppService) SendBulkNotification(message string, phones []string) error {
	for _, phone := range phones {
		if err := s.client.SendText(phone, message); err != nil {
			logger.Error("Failed to send bulk notification",
				zap.String("phone", phone),
			)
		}
	}
	return nil
}

func (s *WhatsAppService) SetDeviceID(deviceID string) {
	s.client.SetDeviceID(deviceID)
}

func (s *WhatsAppService) GetClient() gowa.GOWAClient {
	return s.client
}

func (s *WhatsAppService) GetCustomerByPhone(phone string) (*entities.Customer, error) {
	formattedPhone := formatPhone(phone)
	return s.customerRepo.FindByPhone(formattedPhone)
}

func (s *WhatsAppService) IsAdmin(phone string) bool {
	formattedPhone := formatPhone(phone)
	for _, adminPhone := range s.adminPhones {
		if formatPhone(adminPhone) == formattedPhone {
			return true
		}
	}
	return false
}

func formatPhone(phone string) string {
	if len(phone) >= 10 && phone[0:2] == "08" {
		return "62" + phone[1:]
	}
	if len(phone) >= 10 && phone[0:1] == "0" {
		return "62" + phone[1:]
	}
	return phone
}
