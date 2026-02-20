package handlers

import (
	"net/http"
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/interface/dto"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type HotspotHandler struct {
	profileUC usecase.HotspotProfileUsecase
	userUC    usecase.HotspotUserUsecase
	voucherUC usecase.VoucherUsecase
	saleUC    usecase.HotspotSaleUsecase
	sessionUC usecase.HotspotSessionUsecase
}

func NewHotspotHandler(
	profileUC usecase.HotspotProfileUsecase,
	userUC usecase.HotspotUserUsecase,
	voucherUC usecase.VoucherUsecase,
	saleUC usecase.HotspotSaleUsecase,
	sessionUC usecase.HotspotSessionUsecase,
) *HotspotHandler {
	return &HotspotHandler{
		profileUC: profileUC,
		userUC:    userUC,
		voucherUC: voucherUC,
		saleUC:    saleUC,
		sessionUC: sessionUC,
	}
}

type VoucherTemplateData struct {
	HotspotName string
	BatchName   string
	GeneratedAt string
	TotalCount  int
	DnsName     string
	Currency    string
	LogoURL     string
	Vouchers    []VoucherInfo
}

type VoucherInfo struct {
	Username  string
	Password  string
	Profile   string
	Validity  string
	Price     string
	TimeLimit string
	DataLimit string
	QRCodeURL string
}

func getRouterID(c *gin.Context) (uint, error) {
	routerIDStr := c.DefaultQuery("router_id", "0")
	routerID64, err := strconv.ParseUint(routerIDStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(routerID64), nil
}

func (h *HotspotHandler) CreateProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var req dto.CreateHotspotProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.profileUC.CreateProfile(routerID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Profile created successfully", nil)
}

func (h *HotspotHandler) UpdateProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profileName := c.Param("name")
	var req dto.UpdateHotspotProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.profileUC.UpdateProfile(routerID, profileName, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Profile updated successfully", nil)
}

func (h *HotspotHandler) DeleteProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profileName := c.Param("name")

	if err := h.profileUC.DeleteProfile(routerID, profileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Profile deleted successfully", nil)
}

func (h *HotspotHandler) GetProfiles(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profiles, err := h.profileUC.GetAllProfiles(routerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"profiles": profiles})
}

func (h *HotspotHandler) GetProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profileName := c.Param("name")

	profile, err := h.profileUC.GetProfile(routerID, profileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, profile)
}

func (h *HotspotHandler) CreateUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var req dto.CreateHotspotUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.userUC.CreateUser(routerID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User created successfully", nil)
}

func (h *HotspotHandler) UpdateUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")
	var req dto.UpdateHotspotUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.userUC.UpdateUser(routerID, username, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User updated successfully", nil)
}

func (h *HotspotHandler) DeleteUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	if err := h.userUC.DeleteUser(routerID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User deleted successfully", nil)
}

func (h *HotspotHandler) GetUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	user, err := h.userUC.GetUser(routerID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, user)
}

func (h *HotspotHandler) GetUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var filter dto.HotspotUserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	users, total, err := h.userUC.GetAllUsers(routerID, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"users": users, "total": total})
}

func (h *HotspotHandler) DisableUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	if err := h.userUC.DisableUser(routerID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User disabled successfully", nil)
}

func (h *HotspotHandler) EnableUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	if err := h.userUC.EnableUser(routerID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User enabled successfully", nil)
}

func (h *HotspotHandler) RemoveExpiredUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profile := c.Query("profile")
	if profile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile parameter required"})
		return
	}

	removed, err := h.userUC.RemoveExpiredUsers(routerID, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Expired users removed", gin.H{"removed_count": removed})
}

func (h *HotspotHandler) RemoveUnusedVouchers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	profile := c.Query("profile")
	if profile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile parameter required"})
		return
	}

	removed, err := h.userUC.RemoveUnusedVouchers(routerID, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Unused vouchers removed", gin.H{"removed_count": removed})
}

func (h *HotspotHandler) GenerateVouchers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var req dto.GenerateVouchers
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := h.voucherUC.GenerateVouchers(routerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Vouchers generated successfully", result)
}

func (h *HotspotHandler) GetVouchersByPrefix(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	prefix := c.Param("prefix")

	vouchers, err := h.voucherUC.GetVouchersByPrefix(routerID, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"vouchers": vouchers})
}

func (h *HotspotHandler) PrintVoucher(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	user, err := h.userUC.GetUser(routerID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qrURL, _ := utils.GenerateQRCodeURL(user.Name+":"+user.Password, 150)

	c.HTML(http.StatusOK, "voucher.html", gin.H{
		"HotspotName": "Hotspot WiFi",
		"Username":    user.Name,
		"Password":    user.Password,
		"Profile":     user.Profile,
		"Validity":    "Unlimited",
		"TimeLimit":   "",
		"DataLimit":   "",
		"Price":       "N/A",
		"QRCodeURL":   qrURL,
		"DnsName":     "10.10.10.1",
		"Currency":    "IDR",
		"GeneratedAt": "Just now",
		"ExpiresAt":   "N/A",
		"LogoURL":     "",
	})
}

func (h *HotspotHandler) RecordSale(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var req dto.RecordHotspotSale
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.saleUC.RecordSale(routerID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "Sale recorded successfully", nil)
}

func (h *HotspotHandler) GetSales(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	var filter dto.HotspotSaleFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sales, total, err := h.saleUC.GetAllSales(routerID, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"sales": sales, "total": total})
}

func (h *HotspotHandler) GetTotalRevenue(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	revenue, err := h.saleUC.GetTotalRevenue(routerID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"revenue": revenue})
}

func (h *HotspotHandler) GetActiveSessions(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	sessions, err := h.sessionUC.GetActiveSessions(routerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, gin.H{"sessions": sessions, "total": len(sessions)})
}

func (h *HotspotHandler) GetSession(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	session, err := h.sessionUC.GetSessionByUsername(routerID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, session)
}

func (h *HotspotHandler) DisconnectUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	username := c.Param("username")

	if err := h.sessionUC.DisconnectUser(routerID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessWithMessage(c, "User disconnected successfully", nil)
}

func (h *HotspotHandler) GetSessionStats(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid router_id"})
		return
	}

	stats, err := h.sessionUC.GetSessionStats(routerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccess(c, stats)
}
