package usecase

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/genieacs"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
)

type GenieACSUsecase interface {
	GetDevices() (*dto.GenieACSDeviceListResponse, error)
	GetDevice(serial string) (*dto.GenieACSDeviceResponse, error)
	RebootDevice(serial string) error
	SetParameter(serial, parameter, value string) error
	FindDeviceByPPPoE(username string) (*dto.GenieACSDeviceResponse, error)
}

type genieacsUsecase struct {
	client *genieacs.GenieACSClient
}

func NewGenieACSUsecase(client *genieacs.GenieACSClient) GenieACSUsecase {
	return &genieacsUsecase{
		client: client,
	}
}

func (u *genieacsUsecase) GetDevices() (*dto.GenieACSDeviceListResponse, error) {
	devices, err := u.client.GetDevices()
	if err != nil {
		return nil, err
	}

	deviceDTOs := make([]dto.GenieACSDeviceResponse, 0, len(devices))
	for _, device := range devices {
		info, err := u.client.GetDeviceInfo(device.ID)
		if err != nil {
			continue
		}
		deviceDTOs = append(deviceDTOs, dto.GenieACSDeviceResponse{
			ID:              info.ID,
			SerialNumber:    info.SerialNumber,
			LastInform:      info.LastInform,
			Status:          info.Status,
			Uptime:          info.Uptime,
			Manufacturer:    info.Manufacturer,
			Model:           info.Model,
			SoftwareVersion: info.SoftwareVersion,
			IPAddress:       info.IPAddress,
			MACAddress:      info.MACAddress,
			SSID:            info.SSID,
			WiFiPassword:    info.WiFiPassword,
			RXPower:         info.RXPower,
			TXPower:         info.TXPower,
		})
	}

	return &dto.GenieACSDeviceListResponse{
		Devices: deviceDTOs,
		Total:   len(deviceDTOs),
	}, nil
}

func (u *genieacsUsecase) GetDevice(serial string) (*dto.GenieACSDeviceResponse, error) {
	info, err := u.client.GetDeviceInfo(serial)
	if err != nil {
		return nil, err
	}

	return &dto.GenieACSDeviceResponse{
		ID:              info.ID,
		SerialNumber:    info.SerialNumber,
		LastInform:      info.LastInform,
		Status:          info.Status,
		Uptime:          info.Uptime,
		Manufacturer:    info.Manufacturer,
		Model:           info.Model,
		SoftwareVersion: info.SoftwareVersion,
		IPAddress:       info.IPAddress,
		MACAddress:      info.MACAddress,
		SSID:            info.SSID,
		WiFiPassword:    info.WiFiPassword,
		RXPower:         info.RXPower,
		TXPower:         info.TXPower,
	}, nil
}

func (u *genieacsUsecase) RebootDevice(serial string) error {
	return u.client.Reboot(serial)
}

func (u *genieacsUsecase) SetParameter(serial, parameter, value string) error {
	return u.client.SetParameter(serial, parameter, value)
}

func (u *genieacsUsecase) FindDeviceByPPPoE(username string) (*dto.GenieACSDeviceResponse, error) {
	device, err := u.client.FindDeviceByPPPoE(username)
	if err != nil {
		return nil, err
	}

	info, err := u.client.GetDeviceInfo(device.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GenieACSDeviceResponse{
		ID:              info.ID,
		SerialNumber:    info.SerialNumber,
		LastInform:      info.LastInform,
		Status:          info.Status,
		Uptime:          info.Uptime,
		Manufacturer:    info.Manufacturer,
		Model:           info.Model,
		SoftwareVersion: info.SoftwareVersion,
		IPAddress:       info.IPAddress,
		MACAddress:      info.MACAddress,
		SSID:            info.SSID,
		WiFiPassword:    info.WiFiPassword,
		RXPower:         info.RXPower,
		TXPower:         info.TXPower,
	}, nil
}
