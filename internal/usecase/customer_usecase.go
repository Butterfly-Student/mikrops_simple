package usecase

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/whatsapp"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/logger"

	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"go.uber.org/zap"
)

type CustomerUsecase interface {
	GetCustomers(page, perPage int, search string) (*dto.CustomerListResponse, error)
	GetCustomerByID(id uint) (*dto.CustomerDetail, error)
	CreateCustomer(customer *dto.CustomerDetail) error
	UpdateCustomer(id uint, customer *dto.CustomerDetail) error
	DeleteCustomer(id uint) error
	IsolateCustomer(id uint) error
	ActivateCustomer(id uint) error
	SyncCustomer(id uint) error
	BulkIsolate(ids []uint) error
	BulkActivate(ids []uint) error
}

type customerUsecase struct {
	customerRepo    repositories.CustomerRepository
	mikrotikService *mikrotik.MikroTikService
	whatsappService *whatsapp.WhatsAppService
}

func NewCustomerUsecase(customerRepo repositories.CustomerRepository, mikrotikService *mikrotik.MikroTikService, whatsappService *whatsapp.WhatsAppService) CustomerUsecase {
	return &customerUsecase{
		customerRepo:    customerRepo,
		mikrotikService: mikrotikService,
		whatsappService: whatsappService,
	}
}

func (u *customerUsecase) GetCustomers(page, perPage int, search string) (*dto.CustomerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	customers, total, err := u.customerRepo.FindAll(page, perPage, search)
	if err != nil {
		return nil, err
	}

	customerDTOs := make([]dto.CustomerDetail, 0, len(customers))
	for _, customer := range customers {
		customerDTOs = append(customerDTOs, *u.entityToDTO(customer))
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.CustomerListResponse{
		Customers:  customerDTOs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (u *customerUsecase) GetCustomerByID(id uint) (*dto.CustomerDetail, error) {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return u.entityToDTO(customer), nil
}

func (u *customerUsecase) CreateCustomer(customerDTO *dto.CustomerDetail) error {
	customer := &entities.Customer{
		Name:          customerDTO.Name,
		Phone:         customerDTO.Phone,
		Email:         customerDTO.Email,
		Address:       customerDTO.Address,
		PackageID:     customerDTO.PackageID,
		PPPoEUsername: customerDTO.PPPoEUsername,
		PPPoEPassword: customerDTO.PPPoEPassword,
		Status:        "active",
		RouterID:      customerDTO.RouterID,
		ONUID:         customerDTO.ONUID,
		ONUSerial:     customerDTO.ONUSerial,
		ONUMacAddress: customerDTO.ONUMacAddress,
		ONUIPAddress:  customerDTO.ONUIPAddress,
		Latitude:      customerDTO.Latitude,
		Longitude:     customerDTO.Longitude,
	}

	// Step 1: Jika ada PPPoE credentials, daftarkan ke MikroTik TERLEBIH DAHULU.
	// Jika MikroTik gagal, customer tidak disimpan ke database.
	if customer.PPPoEUsername != "" && customer.PPPoEPassword != "" {
		if u.mikrotikService == nil {
			return fmt.Errorf("mikrotik service tidak tersedia")
		}
		logger.Info("Mendaftarkan customer ke MikroTik",
			zap.String("username", customer.PPPoEUsername),
			zap.Uint("router_id", customer.RouterID),
		)
		if err := u.mikrotikService.CreateCustomerOnMikroTik(customer); err != nil {
			logger.Error("Gagal mendaftarkan customer ke MikroTik — customer TIDAK disimpan ke database",
				zap.String("username", customer.PPPoEUsername),
				zap.Error(err),
			)
			return fmt.Errorf("gagal membuat PPPoE user di MikroTik: %w", err)
		}
		logger.Info("Berhasil mendaftarkan customer ke MikroTik",
			zap.String("username", customer.PPPoEUsername),
		)
	}

	// Step 2: Simpan ke database setelah MikroTik berhasil (atau tidak ada PPPoE).
	if err := u.customerRepo.Create(customer); err != nil {
		logger.Error("Gagal menyimpan customer ke database",
			zap.String("name", customer.Name),
			zap.Error(err),
		)
		return err
	}

	if u.whatsappService != nil {
		go u.whatsappService.SendWelcomeMessage(customer)
	}

	return nil
}

func (u *customerUsecase) UpdateCustomer(id uint, customerDTO *dto.CustomerDetail) error {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return err
	}

	customer.Name = customerDTO.Name
	customer.Phone = customerDTO.Phone
	customer.Email = customerDTO.Email
	customer.Address = customerDTO.Address
	customer.PackageID = customerDTO.PackageID
	customer.RouterID = customerDTO.RouterID
	customer.ONUID = customerDTO.ONUID
	customer.ONUSerial = customerDTO.ONUSerial
	customer.ONUMacAddress = customerDTO.ONUMacAddress
	customer.ONUIPAddress = customerDTO.ONUIPAddress
	customer.Latitude = customerDTO.Latitude
	customer.Longitude = customerDTO.Longitude

	if customerDTO.PPPoEPassword != "" {
		customer.PPPoEPassword = customerDTO.PPPoEPassword
	}

	if err := u.customerRepo.Update(customer); err != nil {
		return err
	}

	return nil
}

func (u *customerUsecase) DeleteCustomer(id uint) error {
	return u.customerRepo.Delete(id)
}

func (u *customerUsecase) IsolateCustomer(id uint) error {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := u.mikrotikService.IsolateCustomer(customer); err != nil {
		return err
	}

	if u.whatsappService != nil {
		go u.whatsappService.SendIsolationNotification(customer)
	}

	return nil
}

func (u *customerUsecase) ActivateCustomer(id uint) error {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := u.mikrotikService.ActivateCustomer(customer); err != nil {
		return err
	}

	if u.whatsappService != nil {
		go u.whatsappService.SendActivationNotification(customer)
	}

	return nil
}

func (u *customerUsecase) SyncCustomer(id uint) error {
	customer, err := u.customerRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := u.mikrotikService.SyncCustomerToMikroTik(customer); err != nil {
		return err
	}

	return nil
}

func (u *customerUsecase) BulkIsolate(ids []uint) error {
	return u.mikrotikService.BulkSyncCustomers(ids)
}

func (u *customerUsecase) BulkActivate(ids []uint) error {
	return u.mikrotikService.BulkSyncCustomers(ids)
}

func (u *customerUsecase) entityToDTO(customer *entities.Customer) *dto.CustomerDetail {
	packageName := ""
	price := 0.0
	if customer.Package != nil {
		packageName = customer.Package.Name
		price = customer.Package.Price
	}

	var isolationDate, activationDate *string
	if customer.IsolationDate != nil {
		date := customer.IsolationDate.Format("2006-01-02 15:04:05")
		isolationDate = &date
	}
	if customer.ActivationDate != nil {
		date := customer.ActivationDate.Format("2006-01-02 15:04:05")
		activationDate = &date
	}

	return &dto.CustomerDetail{
		ID:             customer.ID,
		Name:           customer.Name,
		Phone:          customer.Phone,
		Email:          customer.Email,
		Address:        customer.Address,
		PackageID:      customer.PackageID,
		PackageName:    packageName,
		PackagePrice:   price,
		PPPoEUsername:  customer.PPPoEUsername,
		PPPoEPassword:  "",
		Status:         customer.Status,
		RouterID:       customer.RouterID,
		ONUID:          customer.ONUID,
		ONUSerial:      customer.ONUSerial,
		ONUMacAddress:  customer.ONUMacAddress,
		ONUIPAddress:   customer.ONUIPAddress,
		Latitude:       customer.Latitude,
		Longitude:      customer.Longitude,
		IsolationDate:  isolationDate,
		ActivationDate: activationDate,
		CreatedAt:      customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
