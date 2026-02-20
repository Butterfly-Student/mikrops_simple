package middleware

import (
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/casbin"
	"github.com/alijayanet/gembok-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RBACMiddleware struct {
	casbinService *casbin.CasbinService
}

func NewRBACMiddleware(casbinService *casbin.CasbinService) *RBACMiddleware {
	return &RBACMiddleware{casbinService: casbinService}
}

func (m *RBACMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.SendError(c, 401, "Role not found in context")
			c.AbortWithStatus(401)
			return
		}

		userRole, ok := role.(string)
		if !ok {
			utils.SendError(c, 401, "Invalid role format")
			c.AbortWithStatus(401)
			return
		}

		roleAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			utils.SendError(c, 403, "Insufficient permissions")
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}

func (m *RBACMiddleware) EnforceWithOwner(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.SendError(c, 401, "Role not found in context")
			c.AbortWithStatus(401)
			return
		}

		userRole, ok := role.(string)
		if !ok {
			utils.SendError(c, 401, "Invalid role format")
			c.AbortWithStatus(401)
			return
		}

		owner := "*"
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint); ok {
				owner = strconv.FormatUint(uint64(uid), 10)
			}
		}

		allowed, err := m.casbinService.EnforceWithOwner(userRole, resource, action, owner)
		if err != nil {
			utils.SendError(c, 500, "Authorization check failed")
			c.AbortWithStatus(500)
			return
		}

		if !allowed {
			utils.SendError(c, 403, "Insufficient permissions")
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}

func (m *RBACMiddleware) RequireOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.SendError(c, 401, "Role not found in context")
			c.AbortWithStatus(401)
			return
		}

		_, ok := role.(string)
		if !ok {
			utils.SendError(c, 401, "Invalid role format")
			c.AbortWithStatus(401)
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			utils.SendError(c, 401, "User ID not found in context")
			c.AbortWithStatus(401)
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			utils.SendError(c, 401, "Invalid user ID format")
			c.AbortWithStatus(401)
			return
		}

		resourceID := c.Param("id")

		resourceOwner := strconv.FormatUint(uint64(uid), 10)
		if resourceID != resourceOwner {
			utils.SendError(c, 403, "Access denied: You can only access your own resources")
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}
