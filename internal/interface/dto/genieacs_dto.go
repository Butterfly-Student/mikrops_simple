package dto

type GenieACSDeviceResponse struct {
	ID              string `json:"id"`
	SerialNumber    string `json:"serial_number"`
	LastInform      string `json:"last_inform"`
	Status          string `json:"status"`
	Uptime          string `json:"uptime"`
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	SoftwareVersion string `json:"software_version"`
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	SSID            string `json:"ssid"`
	WiFiPassword    string `json:"wifi_password"`
	RXPower         string `json:"rx_power"`
	TXPower         string `json:"tx_power"`
}

type GenieACSDeviceListResponse struct {
	Devices []GenieACSDeviceResponse `json:"devices"`
	Total   int                      `json:"total"`
}

type GenieACSDeviceRebootRequest struct {
	Serial string `json:"serial" binding:"required"`
}

type GenieACSParameterSetRequest struct {
	Serial    string `json:"serial" binding:"required"`
	Parameter string `json:"parameter" binding:"required"`
	Value     string `json:"value" binding:"required"`
}
