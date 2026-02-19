package handlers

import (
	"net/http"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type GenieACSHandler struct {
	usecase GenieACSUsecase
}

func NewGenieACSHandler(usecase GenieACSUsecase) *GenieACSHandler {
	return &GenieACSHandler{
		usecase: usecase,
	}
}

type GenieACSUsecase interface {
	GetDevices() (*dto.GenieACSDeviceListResponse, error)
	GetDevice(serial string) (*dto.GenieACSDeviceResponse, error)
	RebootDevice(serial string) error
	SetParameter(serial, parameter, value string) error
	FindDeviceByPPPoE(username string) (*dto.GenieACSDeviceResponse, error)
}

func (h *GenieACSHandler) GetDevices(c *gin.Context) {
	devices, err := h.usecase.GetDevices()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Devices retrieved successfully", devices)
}

func (h *GenieACSHandler) GetDevice(c *gin.Context) {
	serial := c.Param("serial")

	device, err := h.usecase.GetDevice(serial)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Device retrieved successfully", device)
}

func (h *GenieACSHandler) RebootDevice(c *gin.Context) {
	var req dto.GenieACSDeviceRebootRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.RebootDevice(req.Serial); err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Device reboot initiated successfully", nil)
}

func (h *GenieACSHandler) SetParameter(c *gin.Context) {
	var req dto.GenieACSParameterSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.SetParameter(req.Serial, req.Parameter, req.Value); err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Parameter set successfully", nil)
}

func (h *GenieACSHandler) FindDeviceByPPPoE(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		utils.SendError(c, http.StatusBadRequest, "Username query parameter is required")
		return
	}

	device, err := h.usecase.FindDeviceByPPPoE(username)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessWithMessage(c, "Device found successfully", device)
}
