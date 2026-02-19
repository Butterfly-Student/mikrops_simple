package genieacs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

type GenieACSClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	timeout    time.Duration
}

func NewGenieACSClient(baseURL, username, password string) *GenieACSClient {
	return &GenieACSClient{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 10 * time.Second,
	}
}

func (c *GenieACSClient) makeRequest(method, endpoint string, body io.Reader, queryParams map[string]string) ([]byte, error) {
	u := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	if len(queryParams) > 0 {
		values := url.Values{}
		for k, v := range queryParams {
			values.Set(k, v)
		}
		u = fmt.Sprintf("%s?%s", u, values.Encode())
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *GenieACSClient) GetDevices() ([]GenieACSDevice, error) {
	projection := []string{
		"_id",
		"_lastInform",
		"_deviceId",
		"DeviceID",
		"VirtualParameters.pppoeUsername",
		"VirtualParameters.pppoeUsername2",
		"VirtualParameters.gettemp",
		"VirtualParameters.RXPower",
		"VirtualParameters.pppoeIP",
		"VirtualParameters.IPTR069",
		"VirtualParameters.pppoeMac",
		"VirtualParameters.getponmode",
		"VirtualParameters.PonMac",
		"VirtualParameters.getSerialNumber",
		"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID",
		"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase",
		"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.TotalAssociations",
		"VirtualParameters.activedevices",
		"VirtualParameters.getdeviceuptime",
	}

	query := map[string]interface{}{"_id": map[string]interface{}{"$regex": ""}}
	queryJSON, _ := json.Marshal(query)
	projectionStr := strings.Join(projection, ",")

	queryParams := map[string]string{
		"query":      string(queryJSON),
		"projection": projectionStr,
	}

	respBody, err := c.makeRequest("GET", "/devices/", nil, queryParams)
	if err != nil {
		return nil, err
	}

	var devices []GenieACSDevice
	if err := json.Unmarshal(respBody, &devices); err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	return devices, nil
}

func (c *GenieACSClient) GetDevice(serial string) (*GenieACSDevice, error) {
	attempts := []struct {
		query map[string]interface{}
	}{
		{map[string]interface{}{"_deviceId._SerialNumber": serial}},
		{map[string]interface{}{"_id": serial}},
		{map[string]interface{}{"_id": url.QueryEscape(serial)}},
	}

	for _, attempt := range attempts {
		queryJSON, _ := json.Marshal(attempt.query)
		queryParams := map[string]string{"query": string(queryJSON)}

		respBody, err := c.makeRequest("GET", "/devices/", nil, queryParams)
		if err != nil {
			logger.Warn("GenieACS query attempt failed", zap.String("serial", serial), zap.Error(err))
			continue
		}

		var devices []GenieACSDevice
		if err := json.Unmarshal(respBody, &devices); err != nil {
			continue
		}

		if len(devices) > 0 {
			return &devices[0], nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", serial)
}

func (c *GenieACSClient) GetDeviceInfo(serial string) (*GenieACSDeviceInfo, error) {
	device, err := c.GetDevice(serial)
	if err != nil {
		return nil, err
	}

	lastInform := ""
	if device.LastInform != nil {
		if str, ok := device.LastInform.(string); ok {
			lastInform = str
		}
	}

	info := &GenieACSDeviceInfo{
		ID:           device.ID,
		SerialNumber: serial,
		LastInform:   lastInform,
		Status:       "unknown",
	}

	if lastInform != "" {
		t, _ := time.Parse(time.RFC3339, lastInform)
		if time.Since(t) < 5*time.Minute {
			info.Status = "online"
		} else {
			info.Status = "offline"
		}
	}

	info.Manufacturer = c.getValue(device, "InternetGatewayDevice.DeviceInfo.Manufacturer", "Device.DeviceInfo.Manufacturer", "DeviceID.Manufacturer")
	info.Model = c.getValue(device, "InternetGatewayDevice.DeviceInfo.ModelName", "Device.DeviceInfo.ModelName", "DeviceID.ProductClass")
	info.SoftwareVersion = c.getValue(device, "InternetGatewayDevice.DeviceInfo.SoftwareVersion", "Device.DeviceInfo.SoftwareVersion")
	info.Uptime = c.getValue(device, "InternetGatewayDevice.DeviceInfo.UpTime", "Device.DeviceInfo.UpTime")
	info.IPAddress = c.getValue(device, "InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress", "InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANPPPConnection.1.ExternalIPAddress")
	info.MACAddress = c.getValue(device, "InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.MACAddress", "Device.Ethernet.Interface.1.MACAddress")
	info.SSID = c.getValue(device, "InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", "InternetGatewayDevice.LANDevice.1.WiFi.Radio.1.SSID", "Device.WiFi.SSID.1.SSID")
	info.WiFiPassword = c.getValue(device, "InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase", "InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.KeyPassphrase", "Device.WiFi.AccessPoint.1.Security.KeyPassphrase")
	info.RXPower = c.getValue(device, "InternetGatewayDevice.WANDevice.1.X_GponInterafceConfig.RxPower", "Device.Optical.Interface.1.RXPower")
	info.TXPower = c.getValue(device, "InternetGatewayDevice.WANDevice.1.X_GponInterafceConfig.TxPower", "Device.Optical.Interface.1.TXPower")

	return info, nil
}

func (c *GenieACSClient) getValue(device *GenieACSDevice, paths ...string) string {
	for _, path := range paths {
		keys := strings.Split(path, ".")
		value := c.extractValue(device, keys)
		if value != "" {
			return value
		}
	}
	return ""
}

func (c *GenieACSClient) extractValue(device interface{}, keys []string) string {
	current := device

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, ok := v[key]; ok {
				current = val
			} else {
				return ""
			}
		case *GenieACSDevice:
			if key == "_id" {
				return v.ID
			}
			if key == "_lastInform" {
				return v.LastInform.(string)
			}
			if val, ok := v.VirtualParameters[key]; ok {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		default:
			return ""
		}
	}

	if str, ok := current.(string); ok {
		return str
	}
	return ""
}

func (c *GenieACSClient) SetParameter(serial, parameter, value string) error {
	device, err := c.GetDevice(serial)
	if err != nil {
		return err
	}

	deviceID := url.QueryEscape(device.ID)
	taskURL := fmt.Sprintf("/devices/%s/tasks?timeout=3000&connection_request", deviceID)

	taskData := map[string]interface{}{
		"name": "setParameterValues",
		"parameterValues": []interface{}{
			[]interface{}{parameter, value, "xsd:string"},
		},
	}

	body, _ := json.Marshal(taskData)
	respBody, err := c.makeRequest("POST", taskURL, strings.NewReader(string(body)), nil)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err == nil {
		logger.Info("GenieACS parameter set", zap.String("serial", serial), zap.String("parameter", parameter))
	}

	return nil
}

func (c *GenieACSClient) FindDeviceByPPPoE(pppoeUsername string) (*GenieACSDevice, error) {
	attempts := []struct {
		query map[string]interface{}
	}{
		{map[string]interface{}{"VirtualParameters.pppoeUsername": pppoeUsername}},
		{map[string]interface{}{"VirtualParameters.pppoeUsername2": pppoeUsername}},
		{map[string]interface{}{"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Username": pppoeUsername}},
		{map[string]interface{}{"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANPPPConnection.1.Username": pppoeUsername}},
		{map[string]interface{}{"Device.PPP.Interface.1.Credentials.Username": pppoeUsername}},
		{map[string]interface{}{"InternetGatewayDevice.PPPPEngine.PPPoE.UnicastDiscovery.Username": pppoeUsername}},
		{map[string]interface{}{"Device.DeviceInfo.Description": pppoeUsername}},
		{map[string]interface{}{"Device.DeviceInfo.FriendlyName": pppoeUsername}},
	}

	for _, attempt := range attempts {
		queryJSON, _ := json.Marshal(attempt.query)
		queryParams := map[string]string{"query": string(queryJSON)}

		respBody, err := c.makeRequest("GET", "/devices/", nil, queryParams)
		if err != nil {
			continue
		}

		var devices []GenieACSDevice
		if err := json.Unmarshal(respBody, &devices); err != nil {
			continue
		}

		if len(devices) > 0 {
			return &devices[0], nil
		}
	}

	queryParams := map[string]string{"query": fmt.Sprintf(`"%s"`, pppoeUsername)}
	respBody, err := c.makeRequest("GET", "/devices/", nil, queryParams)
	if err == nil {
		var devices []GenieACSDevice
		if err := json.Unmarshal(respBody, &devices); err == nil && len(devices) > 0 {
			return &devices[0], nil
		}
	}

	return nil, fmt.Errorf("device not found with PPPoE username: %s", pppoeUsername)
}

func (c *GenieACSClient) Reboot(serial string) error {
	device, err := c.GetDevice(serial)
	if err != nil {
		return err
	}

	deviceID := url.QueryEscape(device.ID)
	taskURL := fmt.Sprintf("/devices/%s/tasks?connection_request", deviceID)

	taskData := map[string]interface{}{
		"name": "reboot",
	}

	body, _ := json.Marshal(taskData)
	_, err = c.makeRequest("POST", taskURL, strings.NewReader(string(body)), nil)
	if err != nil {
		return err
	}

	logger.Info("GenieACS device reboot initiated", zap.String("serial", serial))
	return nil
}

// GetDeviceWithStatus is a convenience alias for GetDeviceInfo
func (c *GenieACSClient) GetDeviceWithStatus(serial string) (*GenieACSDeviceInfo, error) {
	return c.GetDeviceInfo(serial)
}

// FindDeviceByPPPoEUsername is an alias for FindDeviceByPPPoE
func (c *GenieACSClient) FindDeviceByPPPoEUsername(pppoeUsername string) (*GenieACSDevice, error) {
	return c.FindDeviceByPPPoE(pppoeUsername)
}
