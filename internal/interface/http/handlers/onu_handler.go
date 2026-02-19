package handlers

import (
	"net/http"
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ONUHandler struct {
	onuUsecase *usecase.ONUUsecase
}

func NewONUHandler(onuUsecase *usecase.ONUUsecase) *ONUHandler {
	return &ONUHandler{onuUsecase: onuUsecase}
}

// GET /api/onu-locations
func (h *ONUHandler) GetLocations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	locs, total, err := h.onuUsecase.GetAll(page, perPage)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendPaginatedSuccess(c, locs, total, page, perPage)
}

// POST /api/onu-locations
func (h *ONUHandler) UpsertLocation(c *gin.Context) {
	var req usecase.UpsertONURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if err := h.onuUsecase.Upsert(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ONU location saved successfully",
	})
}

// POST /api/onu-wifi
func (h *ONUHandler) SetWiFi(c *gin.Context) {
	var req struct {
		PPPoEUsername string `json:"pppoe_username"`
		Serial        string `json:"serial"`
		SSID          string `json:"ssid"`
		Password      string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if err := h.onuUsecase.SetWiFi(req.PPPoEUsername, req.Serial, req.SSID, req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "WiFi settings updated successfully",
	})
}
