package handlers

import (
	"fmt"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/casbin"
	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CasbinHandler struct {
	casbinService *casbin.CasbinService
}

func NewCasbinHandler(casbinService *casbin.CasbinService) *CasbinHandler {
	return &CasbinHandler{
		casbinService: casbinService,
	}
}

func (h *CasbinHandler) GetAllPolicies(c *gin.Context) {
	policies, err := h.casbinService.GetAllPolicies()
	if err != nil {
		utils.SendError(c, 403, "Failed to get policies")
		return
	}

	dtos := make([]dto.PolicyResponse, len(policies))
	for i, p := range policies {
		if len(p) >= 4 {
			dtos[i] = dto.PolicyResponse{
				ID:      p[0],
				Subject: p[1],
				Object:  p[2],
				Action:  p[3],
				Owner:   p[4],
			}
		}
	}

	utils.SendSuccess(c, gin.H{"policies": dtos})
}

func (h *CasbinHandler) CreatePolicy(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 403, "Invalid request")
		return
	}

	performedBy := fmt.Sprintf("%v", userID)

	if err := h.casbinService.AddPolicy(c, req.Subject, req.Object, req.Action, req.Owner, performedBy); err != nil {
		utils.SendError(c, 403, "Failed to create policy")
		return
	}

	utils.SendSuccessWithMessage(c, "Policy created successfully", nil)
}

func (h *CasbinHandler) DeletePolicy(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req dto.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 403, "Invalid request")
		return
	}

	performedBy := fmt.Sprintf("%v", userID)

	if err := h.casbinService.RemovePolicy(c, req.Subject, req.Object, req.Action, req.Owner, performedBy); err != nil {
		utils.SendError(c, 403, "Failed to delete policy")
		return
	}

	utils.SendSuccessWithMessage(c, "Policy deleted successfully", nil)
}

func (h *CasbinHandler) ReloadPolicies(c *gin.Context) {
	if err := h.casbinService.ReloadPolicies(); err != nil {
		utils.SendError(c, 403, "Failed to reload policies")
		return
	}

	policies, _ := h.casbinService.GetAllPolicies()

	utils.SendSuccess(c, dto.ReloadPoliciesResponse{
		Message:    "Policies reloaded successfully",
		Policies:   len(policies),
		ReloadedAt: utils.GetCurrentTimestamp(),
	})
}

func (h *CasbinHandler) CheckPermission(c *gin.Context) {
	var req dto.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 403, "Invalid request")
		return
	}

	allowed, err := h.casbinService.EnforceWithOwner(req.Role, req.Resource, req.Action, req.Owner)
	if err != nil {
		utils.SendError(c, 403, "Authorization check failed")
		return
	}

	explanation := fmt.Sprintf("Role '%s' %s %s access to %s",
		req.Role,
		map[bool]string{true: "HAS", false: "DOES NOT HAVE"},
		req.Action,
		req.Resource,
	)

	utils.SendSuccess(c, dto.CheckPermissionResponse{
		Allowed:     allowed,
		Role:        req.Role,
		Resource:    req.Resource,
		Action:      req.Action,
		Explanation: explanation,
	})
}
