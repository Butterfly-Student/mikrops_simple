package handlers

import (
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RouterHandler struct {
	routerUC usecase.RouterUsecase
}

func NewRouterHandler(routerUC usecase.RouterUsecase) *RouterHandler {
	return &RouterHandler{
		routerUC: routerUC,
	}
}

func (h *RouterHandler) GetRouters(c *gin.Context) {
	routers, err := h.routerUC.GetAll()
	if err != nil {
		utils.SendError(c, 500, "Failed to get routers")
		return
	}

	utils.SendSuccess(c, routers)
}

func (h *RouterHandler) GetRouter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	router, err := h.routerUC.GetByID(uint(id))
	if err != nil {
		utils.SendError(c, 404, "Router not found")
		return
	}

	utils.SendSuccess(c, router)
}

func (h *RouterHandler) CreateRouter(c *gin.Context) {
	var req dto.RouterCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.routerUC.Create(&req); err != nil {
		utils.SendError(c, 500, "Failed to create router")
		return
	}

	utils.SendSuccessWithMessage(c, "Router created successfully", nil)
}

func (h *RouterHandler) UpdateRouter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	var req dto.RouterUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.routerUC.Update(uint(id), &req); err != nil {
		utils.SendError(c, 500, "Failed to update router")
		return
	}

	utils.SendSuccessWithMessage(c, "Router updated successfully", nil)
}

func (h *RouterHandler) DeleteRouter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	if err := h.routerUC.Delete(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to delete router")
		return
	}

	utils.SendSuccessWithMessage(c, "Router deleted successfully", nil)
}

func (h *RouterHandler) TestConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	result, err := h.routerUC.TestConnection(uint(id))
	if err != nil {
		utils.SendError(c, 500, "Failed to test connection")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *RouterHandler) SetActive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	if err := h.routerUC.SetActive(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to set active router")
		return
	}

	utils.SendSuccessWithMessage(c, "Router set as active", nil)
}

func (h *RouterHandler) GetActive(c *gin.Context) {
	router, err := h.routerUC.GetActive()
	if err != nil {
		utils.SendError(c, 500, "Failed to get active router")
		return
	}

	utils.SendSuccess(c, router)
}

func (h *RouterHandler) GetStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid router ID")
		return
	}

	status, err := h.routerUC.GetStatus(uint(id))
	if err != nil {
		utils.SendError(c, 500, "Failed to get router status")
		return
	}

	utils.SendSuccess(c, status)
}

func (h *RouterHandler) GetAllStatus(c *gin.Context) {
	statuses, err := h.routerUC.GetAllStatus()
	if err != nil {
		utils.SendError(c, 500, "Failed to get all router statuses")
		return
	}

	utils.SendSuccess(c, statuses)
}
