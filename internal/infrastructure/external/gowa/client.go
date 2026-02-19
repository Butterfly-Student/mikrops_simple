package gowa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type GOWAClient interface {
	SendText(phone, message string) error
	SendImage(phone, imageURL, caption string) error
	SendFile(phone, fileURL, caption string) error
	SendVideo(phone, videoURL, caption string) error
	SendContact(phone, contactName, contactPhone string) error
	SendLocation(phone, lat, long string) error
	CheckConnection() (bool, error)
	GetDevices() ([]map[string]interface{}, error)
	SetDeviceID(deviceID string)
}

type RealGOWAClient struct {
	baseURL    string
	apiKey     string
	deviceID   string
	httpClient *http.Client
}

type MockGOWAClient struct {
	messagesSent []map[string]interface{}
}

type GOWAResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewGOWAClient(baseURL, apiKey, deviceID string, useMock bool) GOWAClient {
	if useMock {
		return &MockGOWAClient{
			messagesSent: []map[string]interface{}{},
		}
	}
	return &RealGOWAClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		deviceID:   deviceID,
		httpClient: &http.Client{},
	}
}

func (c *RealGOWAClient) SendText(phone, message string) error {
	payload := map[string]interface{}{
		"phone":   formatPhone(phone),
		"message": message,
	}
	_, err := c.doRequest("POST", "/send/message", payload)
	return err
}

func (c *RealGOWAClient) SendImage(phone, imageURL, caption string) error {
	payload := map[string]interface{}{
		"phone":     formatPhone(phone),
		"image_url": imageURL,
		"caption":   caption,
	}
	_, err := c.doRequest("POST", "/send/image", payload)
	return err
}

func (c *RealGOWAClient) SendFile(phone, fileURL, caption string) error {
	payload := map[string]interface{}{
		"phone":    formatPhone(phone),
		"file_url": fileURL,
		"caption":  caption,
	}
	_, err := c.doRequest("POST", "/send/file", payload)
	return err
}

func (c *RealGOWAClient) SendVideo(phone, videoURL, caption string) error {
	payload := map[string]interface{}{
		"phone":     formatPhone(phone),
		"video_url": videoURL,
		"caption":   caption,
	}
	_, err := c.doRequest("POST", "/send/video", payload)
	return err
}

func (c *RealGOWAClient) SendContact(phone, contactName, contactPhone string) error {
	payload := map[string]interface{}{
		"phone":         formatPhone(phone),
		"contact_name":  contactName,
		"contact_phone": contactPhone,
	}
	_, err := c.doRequest("POST", "/send/contact", payload)
	return err
}

func (c *RealGOWAClient) SendLocation(phone, lat, long string) error {
	payload := map[string]interface{}{
		"phone":     formatPhone(phone),
		"latitude":  lat,
		"longitude": long,
	}
	_, err := c.doRequest("POST", "/send/location", payload)
	return err
}

func (c *RealGOWAClient) CheckConnection() (bool, error) {
	_, err := c.doRequest("GET", "/app/status", nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *RealGOWAClient) GetDevices() ([]map[string]interface{}, error) {
	resp, err := c.doRequest("GET", "/devices", nil)
	if err != nil {
		return nil, err
	}
	devices := resp.([]map[string]interface{})
	return devices, nil
}

func (c *RealGOWAClient) SetDeviceID(deviceID string) {
	c.deviceID = deviceID
}

func (c *RealGOWAClient) doRequest(method, path string, body interface{}) (interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.apiKey != "" {
		parts := strings.Split(c.apiKey, ":")
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		}
	}

	if c.deviceID != "" {
		req.Header.Set("X-Device-Id", c.deviceID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GOWAResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if result.Status != 200 && result.Status != 201 {
		return nil, fmt.Errorf("GOWA API error: %s", result.Message)
	}

	return result.Data, nil
}

func formatPhone(phone string) string {
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}
	if strings.HasPrefix(phone, "0") {
		return "62" + phone[1:]
	}
	if !strings.Contains(phone, "@") {
		return phone + "@s.whatsapp.net"
	}
	return phone
}

func createMultipartRequest(url, apiKey, deviceID string, fields map[string]io.Reader, fieldName string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, reader := range fields {
		part, err := writer.CreateFormFile(fieldName, key)
		if err != nil {
			return nil, err
		}
		io.Copy(part, reader)
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	if apiKey != "" {
		parts := strings.Split(apiKey, ":")
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		}
	}

	if deviceID != "" {
		req.Header.Set("X-Device-Id", deviceID)
	}

	return req, nil
}

func (m *MockGOWAClient) SendText(phone, message string) error {
	log.Printf("[MOCK GOWA] SendText to %s: %s", phone, message)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":    "text",
		"phone":   phone,
		"message": message,
	})
	return nil
}

func (m *MockGOWAClient) SendImage(phone, imageURL, caption string) error {
	log.Printf("[MOCK GOWA] SendImage to %s: %s", phone, imageURL)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":      "image",
		"phone":     phone,
		"image_url": imageURL,
		"caption":   caption,
	})
	return nil
}

func (m *MockGOWAClient) SendFile(phone, fileURL, caption string) error {
	log.Printf("[MOCK GOWA] SendFile to %s: %s", phone, fileURL)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":     "file",
		"phone":    phone,
		"file_url": fileURL,
		"caption":  caption,
	})
	return nil
}

func (m *MockGOWAClient) SendVideo(phone, videoURL, caption string) error {
	log.Printf("[MOCK GOWA] SendVideo to %s: %s", phone, videoURL)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":      "video",
		"phone":     phone,
		"video_url": videoURL,
		"caption":   caption,
	})
	return nil
}

func (m *MockGOWAClient) SendContact(phone, contactName, contactPhone string) error {
	log.Printf("[MOCK GOWA] SendContact to %s: %s", phone, contactName)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":          "contact",
		"phone":         phone,
		"contact_name":  contactName,
		"contact_phone": contactPhone,
	})
	return nil
}

func (m *MockGOWAClient) SendLocation(phone, lat, long string) error {
	log.Printf("[MOCK GOWA] SendLocation to %s: %s, %s", phone, lat, long)
	m.messagesSent = append(m.messagesSent, map[string]interface{}{
		"type":      "location",
		"phone":     phone,
		"latitude":  lat,
		"longitude": long,
	})
	return nil
}

func (m *MockGOWAClient) CheckConnection() (bool, error) {
	log.Printf("[MOCK GOWA] CheckConnection: OK")
	return true, nil
}

func (m *MockGOWAClient) GetDevices() ([]map[string]interface{}, error) {
	log.Printf("[MOCK GOWA] GetDevices")
	return []map[string]interface{}{
		{
			"device_id": "mock-device-1@s.whatsapp.net",
			"status":    "connected",
		},
	}, nil
}

func (m *MockGOWAClient) SetDeviceID(deviceID string) {
	log.Printf("[MOCK GOWA] SetDeviceID: %s", deviceID)
}
