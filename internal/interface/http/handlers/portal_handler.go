package handlers

import (
	"net/http"
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PortalHandler struct {
	portalUsecase *usecase.PortalUsecase
}

func NewPortalHandler(portalUsecase *usecase.PortalUsecase) *PortalHandler {
	return &PortalHandler{portalUsecase: portalUsecase}
}

// POST /api/portal/login  (public)
func (h *PortalHandler) Login(c *gin.Context) {
	var req usecase.PortalLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.portalUsecase.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    resp,
	})
}

// GET /api/portal/profile  (customer auth required)
func (h *PortalHandler) GetProfile(c *gin.Context) {
	customerID := getCustomerID(c)
	if customerID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	customer, err := h.portalUsecase.GetProfile(customerID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccess(c, customer)
}

// PUT /api/portal/password  (customer auth required)
func (h *PortalHandler) ChangePassword(c *gin.Context) {
	customerID := getCustomerID(c)
	if customerID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Password is required")
		return
	}

	if err := h.portalUsecase.ChangePassword(customerID, req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password updated successfully"})
}

// GET /api/portal/invoices  (customer auth required)
func (h *PortalHandler) GetInvoices(c *gin.Context) {
	customerID := getCustomerID(c)
	if customerID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}

	invoices, total, err := h.portalUsecase.GetInvoices(customerID, page, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendPaginatedSuccess(c, invoices, total, page, perPage)
}

// GET /api/portal/tickets  (customer auth required)
func (h *PortalHandler) GetTickets(c *gin.Context) {
	customerID := getCustomerID(c)
	if customerID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tickets, err := h.portalUsecase.GetTickets(customerID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, tickets)
}

// POST /api/portal/tickets  (customer auth required)
func (h *PortalHandler) CreateTicket(c *gin.Context) {
	customerID := getCustomerID(c)
	if customerID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Subject     string `json:"subject"`
		Description string `json:"description" binding:"required"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	ticket, err := h.portalUsecase.CreateTicket(customerID, req.Description, priority, req.Subject)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Ticket submitted",
		"data":    ticket,
	})
}

// Helper to get authenticated customer ID from context
func getCustomerID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}
