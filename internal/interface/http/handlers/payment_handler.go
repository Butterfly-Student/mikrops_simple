package handlers

import (
	"io"
	"net/http"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/tripay"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

// POST /api/payment/create
func (h *PaymentHandler) CreateTransaction(c *gin.Context) {
	var req usecase.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.paymentUsecase.CreateTransaction(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment transaction created",
		"data":    resp,
	})
}

// GET /api/payment/gateways
func (h *PaymentHandler) GetGateways(c *gin.Context) {
	gateways, err := h.paymentUsecase.GetPaymentGateways()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gateways,
	})
}

// POST /api/payment/callback  (Tripay webhook - no auth required)
func (h *PaymentHandler) TripayCallback(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	signature := c.GetHeader("X-Callback-Signature")

	var payload tripay.TripayCallbackPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	if err := h.paymentUsecase.HandleCallback(payload, string(rawBody), signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
