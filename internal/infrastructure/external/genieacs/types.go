package genieacs

type GenieACSDevice struct {
	ID                string                 `json:"_id"`
	LastInform        interface{}            `json:"_lastInform"`
	DeviceID          map[string]interface{} `json:"_deviceId"`
	VirtualParameters map[string]interface{} `json:"VirtualParameters"`
	LANDevice         map[string]interface{} `json:"InternetGatewayDevice.LANDevice.1"`
	WANDevice         map[string]interface{} `json:"InternetGatewayDevice.WANDevice.1"`
	DeviceInfo        map[string]interface{} `json:"InternetGatewayDevice.DeviceInfo"`
}

type GenieACSDeviceInfo struct {
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

type GenieACSParameter struct {
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}

type GenieACSTask struct {
	ID         string              `json:"_id"`
	Name       string              `json:"name"`
	Parameters []GenieACSParameter `json:"parameterValues"`
	Status     string              `json:"status"`
	Timestamp  string              `json:"timestamp"`
}
