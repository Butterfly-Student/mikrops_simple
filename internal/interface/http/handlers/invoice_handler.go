package handlers

import (
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceUsecase usecase.InvoiceUsecase
}

func NewInvoiceHandler(invoiceUsecase usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceUsecase: invoiceUsecase,
	}
}

func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := h.invoiceUsecase.GetInvoices(page, perPage)
	if err != nil {
		utils.SendError(c, 500, "Failed to get invoices")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *InvoiceHandler) GetInvoiceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid invoice ID")
		return
	}

	invoice, err := h.invoiceUsecase.GetInvoiceByID(uint(id))
	if err != nil {
		utils.SendError(c, 404, "Invoice not found")
		return
	}

	utils.SendSuccess(c, invoice)
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req dto.InvoiceDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.invoiceUsecase.CreateInvoice(&req); err != nil {
		utils.SendError(c, 500, "Failed to create invoice: "+err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Invoice created successfully", nil)
}

func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid invoice ID")
		return
	}

	var req dto.InvoiceDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.invoiceUsecase.UpdateInvoice(uint(id), &req); err != nil {
		utils.SendError(c, 500, "Failed to update invoice: "+err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Invoice updated successfully", nil)
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid invoice ID")
		return
	}

	if err := h.invoiceUsecase.DeleteInvoice(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to delete invoice")
		return
	}

	utils.SendSuccessWithMessage(c, "Invoice deleted successfully", nil)
}
