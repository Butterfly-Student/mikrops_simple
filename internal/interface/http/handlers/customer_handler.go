package handlers

import (
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")

	result, err := h.customerUsecase.GetCustomers(page, perPage, search)
	if err != nil {
		utils.SendError(c, 500, "Failed to get customers")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	customer, err := h.customerUsecase.GetCustomerByID(uint(id))
	if err != nil {
		utils.SendError(c, 404, "Customer not found")
		return
	}

	utils.SendSuccess(c, customer)
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req dto.CustomerDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.customerUsecase.CreateCustomer(&req); err != nil {
		utils.SendError(c, 500, "Failed to create customer: "+err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Customer created successfully", nil)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	var req dto.CustomerDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.customerUsecase.UpdateCustomer(uint(id), &req); err != nil {
		utils.SendError(c, 500, "Failed to update customer: "+err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Customer updated successfully", nil)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.customerUsecase.DeleteCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to delete customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer deleted successfully", nil)
}

func (h *CustomerHandler) IsolateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.customerUsecase.IsolateCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to isolate customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer isolated successfully", nil)
}

func (h *CustomerHandler) ActivateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.customerUsecase.ActivateCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to activate customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer activated successfully", nil)
}

func (h *CustomerHandler) SyncCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.customerUsecase.SyncCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to sync customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer synced successfully", nil)
}

func (h *CustomerHandler) BulkIsolate(c *gin.Context) {
	var req dto.BulkOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.customerUsecase.BulkIsolate(req.CustomerIDs); err != nil {
		utils.SendError(c, 500, "Failed to bulk isolate customers")
		return
	}

	utils.SendSuccessWithMessage(c, "Bulk isolation completed", nil)
}

func (h *CustomerHandler) BulkActivate(c *gin.Context) {
	var req dto.BulkOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.customerUsecase.BulkActivate(req.CustomerIDs); err != nil {
		utils.SendError(c, 500, "Failed to bulk activate customers")
		return
	}

	utils.SendSuccessWithMessage(c, "Bulk activation completed", nil)
}
