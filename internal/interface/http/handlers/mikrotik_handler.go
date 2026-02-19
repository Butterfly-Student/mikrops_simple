package handlers

import (
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MikroTikHandler struct {
	mikrotikUC usecase.MikroTikUsecase
}

func NewMikroTikHandler(mikrotikUC usecase.MikroTikUsecase) *MikroTikHandler {
	return &MikroTikHandler{
		mikrotikUC: mikrotikUC,
	}
}

func (h *MikroTikHandler) GetPPPUsers(c *gin.Context) {
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	result, err := h.mikrotikUC.GetPPPUsers(uint(routerID))
	if err != nil {
		utils.SendError(c, 500, "Failed to get PPP users")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *MikroTikHandler) GetActiveSessions(c *gin.Context) {
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	result, err := h.mikrotikUC.GetActiveSessions(uint(routerID))
	if err != nil {
		utils.SendError(c, 500, "Failed to get active sessions")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *MikroTikHandler) GetPPPProfiles(c *gin.Context) {
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	result, err := h.mikrotikUC.GetPPPProfiles(uint(routerID))
	if err != nil {
		utils.SendError(c, 500, "Failed to get PPP profiles")
		return
	}

	utils.SendSuccess(c, result)
}

func (h *MikroTikHandler) AddPPPUser(c *gin.Context) {
	var req dto.AddPPPUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.mikrotikUC.AddPPPUser(&req); err != nil {
		utils.SendError(c, 500, "Failed to add PPP user")
		return
	}

	utils.SendSuccessWithMessage(c, "PPP user added successfully", nil)
}

func (h *MikroTikHandler) UpdatePPPUser(c *gin.Context) {
	username := c.Param("username")
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.mikrotikUC.UpdatePPPUser(username, uint(routerID), params); err != nil {
		utils.SendError(c, 500, "Failed to update PPP user")
		return
	}

	utils.SendSuccessWithMessage(c, "PPP user updated successfully", nil)
}

func (h *MikroTikHandler) RemovePPPUser(c *gin.Context) {
	username := c.Param("username")
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	if err := h.mikrotikUC.RemovePPPUser(username, uint(routerID)); err != nil {
		utils.SendError(c, 500, "Failed to remove PPP user")
		return
	}

	utils.SendSuccessWithMessage(c, "PPP user removed successfully", nil)
}

func (h *MikroTikHandler) DisconnectUser(c *gin.Context) {
	username := c.Param("username")
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID, _ := strconv.ParseUint(routerIDStr, 10, 32)

	if err := h.mikrotikUC.DisconnectUser(username, uint(routerID)); err != nil {
		utils.SendError(c, 500, "Failed to disconnect user")
		return
	}

	utils.SendSuccessWithMessage(c, "User disconnected successfully", nil)
}

func (h *MikroTikHandler) IsolateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.mikrotikUC.IsolateCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to isolate customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer isolated successfully", nil)
}

func (h *MikroTikHandler) ActivateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.mikrotikUC.ActivateCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to activate customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer activated successfully", nil)
}

func (h *MikroTikHandler) BulkIsolate(c *gin.Context) {
	var req dto.BulkOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.mikrotikUC.BulkIsolate(req.CustomerIDs); err != nil {
		utils.SendError(c, 500, "Failed to bulk isolate customers")
		return
	}

	utils.SendSuccessWithMessage(c, "Bulk isolation completed", nil)
}

func (h *MikroTikHandler) BulkActivate(c *gin.Context) {
	var req dto.BulkOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request")
		return
	}

	if err := h.mikrotikUC.BulkActivate(req.CustomerIDs); err != nil {
		utils.SendError(c, 500, "Failed to bulk activate customers")
		return
	}

	utils.SendSuccessWithMessage(c, "Bulk activation completed", nil)
}

func (h *MikroTikHandler) SyncCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendError(c, 400, "Invalid customer ID")
		return
	}

	if err := h.mikrotikUC.SyncCustomer(uint(id)); err != nil {
		utils.SendError(c, 500, "Failed to sync customer")
		return
	}

	utils.SendSuccessWithMessage(c, "Customer synced successfully", nil)
}

func (h *MikroTikHandler) SyncAllCustomers(c *gin.Context) {
	if err := h.mikrotikUC.SyncAllCustomers(); err != nil {
		utils.SendError(c, 500, "Failed to sync all customers")
		return
	}

	utils.SendSuccessWithMessage(c, "All customers synced successfully", nil)
}

// GET /api/mikrotik/hotspot/logs - stub (MikroTik hotspot log)
func (h *MikroTikHandler) GetHotspotLog(c *gin.Context) {
	limit := 20
	// Return stub data — real implementation requires routeros v2 integration
	c.JSON(200, gin.H{
		"success": true,
		"data":    []interface{}{},
		"message": "Hotspot log endpoint ready. MikroTik routeros v2 integration pending.",
		"limit":   limit,
	})
}

// GET /api/mikrotik/traffic - stub (Real-time traffic monitor)
func (h *MikroTikHandler) GetTraffic(c *gin.Context) {
	iface := c.DefaultQuery("interface", "ether1")
	// Return stub data — real implementation requires routeros v2 integration
	c.JSON(200, []gin.H{
		{"data": 0, "interface": iface, "direction": "tx"},
		{"data": 0, "interface": iface, "direction": "rx"},
	})
}
