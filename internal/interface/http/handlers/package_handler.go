package handlers

import (
	"net/http"
	"strconv"

	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PackageHandler struct {
	packageUC *usecase.PackageUsecase
}

func NewPackageHandler(packageUC *usecase.PackageUsecase) *PackageHandler {
	return &PackageHandler{packageUC: packageUC}
}

// GET /api/packages
func (h *PackageHandler) GetAll(c *gin.Context) {
	packages, err := h.packageUC.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get packages")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": packages})
}

// GET /api/packages/:id
func (h *PackageHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid package ID")
		return
	}
	pkg, err := h.packageUC.GetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pkg})
}

// POST /api/packages
func (h *PackageHandler) Create(c *gin.Context) {
	var req usecase.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	pkg, err := h.packageUC.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": pkg, "message": "Package created successfully"})
}

// PUT /api/packages/:id
func (h *PackageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid package ID")
		return
	}
	var req usecase.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	pkg, err := h.packageUC.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pkg, "message": "Package updated successfully"})
}

// DELETE /api/packages/:id
func (h *PackageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid package ID")
		return
	}
	if err := h.packageUC.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Package deleted successfully"})
}
