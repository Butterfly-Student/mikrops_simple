package usecase

import (
	"fmt"

	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"github.com/alijayanet/gembok-backend/internal/domain/repositories"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/genieacs"
)

type ONUUsecase struct {
	onuRepo        repositories.ONULocationRepository
	genieacsClient *genieacs.GenieACSClient
}

func NewONUUsecase(onuRepo repositories.ONULocationRepository, genieacsClient *genieacs.GenieACSClient) *ONUUsecase {
	return &ONUUsecase{
		onuRepo:        onuRepo,
		genieacsClient: genieacsClient,
	}
}

type ONULocationWithStatus struct {
	entities.ONULocation
	Status     string      `json:"status"`
	DeviceInfo interface{} `json:"device_info,omitempty"`
	SSID       string      `json:"ssid,omitempty"`
}

type UpsertONURequest struct {
	ID           uint    `json:"id"`
	CustomerID   uint    `json:"customer_id"`
	SerialNumber string  `json:"serial_number"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	Notes        string  `json:"notes"`
}

func (u *ONUUsecase) GetAll(page, perPage int) ([]*ONULocationWithStatus, int64, error) {
	locs, total, err := u.onuRepo.FindAll(page, perPage)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*ONULocationWithStatus, 0, len(locs))
	for _, loc := range locs {
		item := &ONULocationWithStatus{ONULocation: *loc}

		if loc.SerialNumber != "" && u.genieacsClient != nil {
			deviceInfo, err := u.genieacsClient.GetDeviceWithStatus(loc.SerialNumber)
			if err == nil && deviceInfo != nil {
				item.Status = deviceInfo.Status
				item.DeviceInfo = deviceInfo
			} else {
				item.Status = "unknown"
			}
		} else {
			item.Status = "unknown"
		}
		result = append(result, item)
	}
	return result, total, nil
}

func (u *ONUUsecase) Upsert(req UpsertONURequest) error {
	if req.SerialNumber == "" {
		return fmt.Errorf("serial_number is required")
	}

	existing, _ := u.onuRepo.FindBySerialNumber(req.SerialNumber)
	if existing != nil {
		existing.Latitude = req.Latitude
		existing.Longitude = req.Longitude
		existing.Address = req.Address
		existing.Notes = req.Notes
		if req.Name != "" {
			existing.ONUID = req.Name
		}
		return u.onuRepo.Update(existing)
	}

	loc := &entities.ONULocation{
		CustomerID:   req.CustomerID,
		SerialNumber: req.SerialNumber,
		ONUID:        req.Name,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		Notes:        req.Notes,
	}
	return u.onuRepo.Create(loc)
}

func (u *ONUUsecase) SetWiFi(pppoeUsername, serial, ssid, password string) error {
	// If pppoe username provided, find device first via GenieACS
	if pppoeUsername != "" && serial == "" {
		device, err := u.genieacsClient.FindDeviceByPPPoEUsername(pppoeUsername)
		if err != nil || device == nil {
			return fmt.Errorf("device not found for pppoe username: %s", pppoeUsername)
		}
		// GenieACSDevice has field ID (the device's _id)
		serial = device.ID
	}

	if serial == "" {
		return fmt.Errorf("serial or pppoe_username is required")
	}

	if ssid != "" {
		if len(ssid) < 3 {
			return fmt.Errorf("ssid must be at least 3 characters")
		}
		if err := u.genieacsClient.SetParameter(serial, "InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", ssid); err != nil {
			return fmt.Errorf("failed to set SSID: %w", err)
		}
	}

	if password != "" {
		if len(password) < 8 {
			return fmt.Errorf("password must be at least 8 characters")
		}
		if err := u.genieacsClient.SetParameter(serial, "InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase", password); err != nil {
			return fmt.Errorf("failed to set password: %w", err)
		}
	}

	return nil
}
