package handlers

import (
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUsecase usecase.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.dashboardUsecase.GetDashboardStats()
	if err != nil {
		utils.SendError(c, 500, "Failed to get dashboard stats")
		return
	}

	utils.SendSuccess(c, stats)
}
