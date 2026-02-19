package tripay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

type TripayClient struct {
	apiKey       string
	privateKey   string
	merchantCode string
	mode         string
	httpClient   *http.Client
}

type TripayTransactionRequest struct {
	Method        string            `json:"method"`
	MerchantRef   string            `json:"merchant_ref"`
	Amount        int64             `json:"amount"`
	CustomerName  string            `json:"customer_name"`
	CustomerEmail string            `json:"customer_email"`
	CustomerPhone string            `json:"customer_phone"`
	OrderItems    []TripayOrderItem `json:"order_items"`
	Signature     string            `json:"signature"`
	ReturnURL     string            `json:"return_url,omitempty"`
	ExpiredTime   int64             `json:"expired_time,omitempty"`
}

type TripayOrderItem struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

type TripayTransactionData struct {
	Reference   string `json:"reference"`
	MerchantRef string `json:"merchant_ref"`
	PaymentURL  string `json:"payment_url"`
	QRString    string `json:"qr_string,omitempty"`
	Status      string `json:"status"`
	Amount      int64  `json:"amount"`
	ExpiredTime int64  `json:"expired_time"`
}

type TripayTransactionResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    *TripayTransactionData `json:"data"`
}

type TripayCallbackPayload struct {
	Reference         string `json:"reference"`
	MerchantRef       string `json:"merchant_ref"`
	PaymentMethod     string `json:"payment_method"`
	PaymentMethodCode string `json:"payment_method_code"`
	TotalAmount       int64  `json:"total_amount"`
	FeeMerchant       int64  `json:"fee_merchant"`
	FeeCustomer       int64  `json:"fee_customer"`
	AmountReceived    int64  `json:"amount_received"`
	IsClosedPayment   int    `json:"is_closed_payment"`
	Status            string `json:"status"`
	PaidAt            int64  `json:"paid_at"`
	Note              string `json:"note"`
}

type TripayPaymentChannel struct {
	Group   string `json:"group"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	IconURL string `json:"icon_url"`
	Active  bool   `json:"active"`
}

func NewTripayClient(apiKey, privateKey, merchantCode, mode string) *TripayClient {
	return &TripayClient{
		apiKey:       apiKey,
		privateKey:   privateKey,
		merchantCode: merchantCode,
		mode:         mode,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *TripayClient) baseURL() string {
	if c.mode == "sandbox" {
		return "https://tripay.co.id/api-sandbox"
	}
	return "https://tripay.co.id/api"
}

func (c *TripayClient) IsConfigured() bool {
	return c.apiKey != "" && c.privateKey != "" && c.merchantCode != ""
}

func (c *TripayClient) GenerateSignature(merchantRef string, amount int64) string {
	data := fmt.Sprintf("%s%s%d", c.merchantCode, merchantRef, amount)
	mac := hmac.New(sha256.New, []byte(c.privateKey))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *TripayClient) CreateTransaction(req TripayTransactionRequest) (*TripayTransactionResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("tripay not configured")
	}

	req.Signature = c.GenerateSignature(req.MerchantRef, req.Amount)

	url := c.baseURL() + "/transaction/create"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.Info("Tripay creating transaction",
		zap.String("merchant_ref", req.MerchantRef),
		zap.Int64("amount", req.Amount),
	)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response TripayTransactionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *TripayClient) GetPaymentChannels() ([]TripayPaymentChannel, error) {
	url := c.baseURL() + "/merchant/payment-channel"

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    []TripayPaymentChannel `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("tripay error: %s", result.Message)
	}

	return result.Data, nil
}

func (c *TripayClient) ValidateCallback(signature, body string) bool {
	mac := hmac.New(sha256.New, []byte(c.privateKey))
	mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
