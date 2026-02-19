package handlers

import (
	"net/http"
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TroubleTicketHandler struct {
	ticketUsecase *usecase.TroubleTicketUsecase
}

func NewTroubleTicketHandler(ticketUsecase *usecase.TroubleTicketUsecase) *TroubleTicketHandler {
	return &TroubleTicketHandler{ticketUsecase: ticketUsecase}
}

// GET /api/tickets
func (h *TroubleTicketHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}

	tickets, total, err := h.ticketUsecase.GetAll(page, perPage, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendPaginatedSuccess(c, tickets, total, page, perPage)
}

// GET /api/tickets/:id
func (h *TroubleTicketHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ticket, err := h.ticketUsecase.GetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Ticket not found")
		return
	}

	utils.SendSuccess(c, ticket)
}

// PUT /api/tickets/:id
func (h *TroubleTicketHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req usecase.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ticket, err := h.ticketUsecase.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Ticket updated", ticket)
}

// POST /api/tickets  (admin can create on behalf of customer)
func (h *TroubleTicketHandler) Create(c *gin.Context) {
	var req usecase.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ticket, err := h.ticketUsecase.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Ticket created",
		"data":    ticket,
	})
}
