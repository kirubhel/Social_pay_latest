package gin

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	auth_entity "github.com/socialpay/socialpay/src/pkg/authv2/core/entity"
	auth_service "github.com/socialpay/socialpay/src/pkg/authv2/core/service"
	rbac_entity "github.com/socialpay/socialpay/src/pkg/rbac/core/entity"
	ginn "github.com/socialpay/socialpay/src/pkg/shared/middleware/gin"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/entity"
	"github.com/socialpay/socialpay/src/pkg/team_member/core/service"
)

// TeamMemberHandler handles team member management HTTP requests
type TeamMemberHandler struct {
	teamMemberService service.TeamMemberService
	rbac              *ginn.RBACV2
	authService       auth_service.AuthService
}

// NewTeamMemberHandler creates a new team member handler
func NewTeamMemberHandler(authService auth_service.AuthService, teamMemberService service.TeamMemberService, rbac *ginn.RBACV2) *TeamMemberHandler {
	return &TeamMemberHandler{
		authService:       authService,
		teamMemberService: teamMemberService,
		rbac:              rbac,
	}
}

// CreateGroupRequest represents the create group route request payload
type CreateGroupRequest struct {
	Title       string                                `json:"title" binding:"required"`
	Description string                                `json:"description" binding:"required"`
	Permissions []rbac_entity.CreatePermissionRequest `json:"permissions" binding:"required"`
}

// UpdateGroupRequest represents the update group route request payload
type UpdateGroupRequest struct {
	Title       *string                                `json:"title,omitempty"`
	Description *string                                `json:"description,omitempty"`
	Permissions *[]rbac_entity.CreatePermissionRequest `json:"permissions,omitempty"`
}

// CreateTeamMemberRequest represents the create team member route request payload
type CreateTeamMemberRequest struct {
	GroupId  uuid.UUID                     `json:"group_id" binding:"required"`
	UserData auth_entity.CreateUserRequest `json:"user_data" binding:"required"`
}

// UpdateTeamMemberRequest represents the update team member route request payload
type UpdateTeamMemberRequest struct {
	GroupId  uuid.UUID                     `json:"group_id" binding:"required"`
	UserData auth_entity.UpdateUserRequest `json:"user_data" binding:"required"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error"`
}

// ApiError represents API error details
type ApiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// CreateGroup handler group creation
// @Summary Create new group
// @Description Create new group with title, description, and set of permissions
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "CreateGroup details"
// @Success 201 {object} map[string]interface{} "Group creation successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Group already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /group/create [post]
func (h *TeamMemberHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	entityReq := &entity.CreateGroupRequest{
		Title:       req.Title,
		Description: req.Description,
		MerchantID:  &merchantId,
		Permissions: req.Permissions,
	}

	// Create group
	err := h.teamMemberService.CreateGroup(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Group creation failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Group creation successfully",
	})
}

// CreateAdminGroup handler group creation
// @Summary Create new admin group
// @Description Create new admin group with title, description, and set of permissions
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "CreateGroup details"
// @Success 201 {object} map[string]interface{} "Group creation successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Group already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/group/create [post]
func (h *TeamMemberHandler) CreateAdminGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	entityReq := &entity.CreateGroupRequest{
		Title:       req.Title,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	// Create group
	err := h.teamMemberService.CreateGroup(c.Request.Context(), entityReq)
	if err != nil {
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Admin Group creation failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Admin Group creation successfully",
	})
}

// GetGroups godoc
// @Summary Get groups
// @Description Get list of groups
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /groups [get]
func (h *TeamMemberHandler) GetGroups(c *gin.Context) {
	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	groups, err := h.teamMemberService.GetGroups(c, &merchantId)
	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Fetching groups failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Groups fetched successfully",
		"data":    groups,
	})
}

// GetAdminGroups godoc
// @Summary Get admin groups
// @Description Get list of admin groups
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/groups [get]
func (h *TeamMemberHandler) GetAdminGroups(c *gin.Context) {
	groups, err := h.teamMemberService.GetGroups(c, nil)
	if err != nil {
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Fetching groups failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Groups fetched successfully",
		"data":    groups,
	})
}

// GetGroup godoc
// @Summary Get group
// @Description Get group by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /group/:groupId [get]
func (h *TeamMemberHandler) GetGroup(c *gin.Context) {
	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	group, err := h.teamMemberService.GetGroupById(c, &merchantId, groupId)

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Fetching group failed",
				},
			})
		}
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrGroupNotFound,
				"message": entity.MsgGroupNotFound,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group fetched successfully",
		"data":    group,
	})
}

// GetAdminGroup godoc
// @Summary Get admin group
// @Description Get admin group by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/group/:groupId [get]
func (h *TeamMemberHandler) GetAdminGroup(c *gin.Context) {
	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	group, err := h.teamMemberService.GetGroupById(c, nil, groupId)

	if err != nil {
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Fetching admin group failed",
				},
			})
		}
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrGroupNotFound,
				"message": entity.MsgGroupNotFound,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group fetched successfully",
		"data":    group,
	})
}

// UpdateGroup godoc
// @Summary Update group
// @Description Update group by id with optional title, description, and permissions
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Param request body UpdateGroupRequest true "Update group details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /group/:groupId [put]
func (h *TeamMemberHandler) UpdateGroup(c *gin.Context) {
	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	merchantId, exists := ginn.GetMerchantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	// Create entity request
	entityReq := &entity.UpdateGroupRequest{
		ID:          groupId,
		Title:       req.Title,
		Description: req.Description,
		MerchantID:  &merchantId,
		Permissions: req.Permissions,
	}

	// Update group
	err = h.teamMemberService.UpdateGroup(c.Request.Context(), entityReq)
	if err != nil {
		log.Printf("Update group error: %v", err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			} else if authErr.Type == entity.ErrGroupNotFound {
				statusCode = http.StatusNotFound
			} else if authErr.Type == entity.ErrGroupAlreadyExists {
				statusCode = http.StatusConflict
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Group update failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group updated successfully",
	})
}

// UpdateAdminGroup godoc
// @Summary Update admin group
// @Description Update admin group by id with optional title, description, and permissions
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Param request body UpdateGroupRequest true "Update group details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/group/:groupId [put]
func (h *TeamMemberHandler) UpdateAdminGroup(c *gin.Context) {
	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	// Create entity request
	entityReq := &entity.UpdateGroupRequest{
		ID:          groupId,
		Title:       req.Title,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	// Update group
	err = h.teamMemberService.UpdateGroup(c.Request.Context(), entityReq)
	if err != nil {
		log.Printf("Update group error: %v", err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			} else if authErr.Type == entity.ErrGroupNotFound {
				statusCode = http.StatusNotFound
			} else if authErr.Type == entity.ErrGroupAlreadyExists {
				statusCode = http.StatusConflict
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Admin Group update failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin Group updated successfully",
	})
}

// DeleteGroup godoc
// @Summary Get group
// @Description Delete group by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /group/:groupId [delete]
func (h *TeamMemberHandler) DeleteGroup(c *gin.Context) {
	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	err = h.teamMemberService.DeleteGroupById(c, &merchantId, groupId)

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Deleting group failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Group deleted successfully",
	})
}

// DeleteAdminGroup godoc
// @Summary Delete admin group
// @Description Delete admin group by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/group/:groupId [delete]
func (h *TeamMemberHandler) DeleteAdminGroup(c *gin.Context) {
	groupId, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INVALID_REQUEST",
				"message": "Invalid group ID",
			},
		})
		return
	}

	err = h.teamMemberService.DeleteGroupById(c, nil, groupId)

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Deleting admin group failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Admin Group deleted successfully",
	})
}

// CreateTeamMember godoc
// @Summary Create team member
// @Description Create team member
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /create [post]
func (h *TeamMemberHandler) CreateTeamMember(c *gin.Context) {
	var req CreateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	req.UserData.UserType = "member"
	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	err := h.teamMemberService.CreateTeamMember(c, &merchantId, &entity.CreateTeamMemberRequest{
		GroupId:  req.GroupId,
		UserData: req.UserData,
	})

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Creating team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Team member created successfully",
	})
}

// CreateAdminTeamMember godoc
// @Summary Create Admin team member
// @Description Create admin team member
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/create [post]
func (h *TeamMemberHandler) CreateAdminTeamMember(c *gin.Context) {
	var req CreateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	req.UserData.UserType = "admin"

	err := h.teamMemberService.CreateTeamMember(c, nil, &entity.CreateTeamMemberRequest{
		GroupId:  req.GroupId,
		UserData: req.UserData,
	})

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Creating team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Team member created successfully",
	})
}

// UpdateTeamMember godoc
// @Summary Update team member
// @Description Update team member user data and group assignment
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param userGroupId path string true "User Group ID"
// @Param request body UpdateTeamMemberRequest true "UpdateTeamMember details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /update/:userGroupId [put]
func (h *TeamMemberHandler) UpdateTeamMember(c *gin.Context) {
	userGroupId, err := uuid.Parse(c.Param("userGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid user group ID",
			},
		})
		return
	}

	var req UpdateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	merchantId, exists := ginn.GetMerchantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	// Extract user ID from request to create entity request
	// Since we don't have user ID in the route, we need to get it from the user group
	err = h.teamMemberService.UpdateTeamMember(c, &merchantId, userGroupId, &entity.UpdateTeamMemberRequest{
		GroupId:  req.GroupId,
		UserData: req.UserData,
	})

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Updating team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team member updated successfully",
	})
}

// UpdateAdminTeamMember godoc
// @Summary Update admin team member
// @Description Update admin team member user data and group assignment
// @Tags Team-member-management
// @Accept json
// @Produce json
// @Param userGroupId path string true "User Group ID"
// @Param request body UpdateTeamMemberRequest true "UpdateTeamMember details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/update/:userGroupId [put]
func (h *TeamMemberHandler) UpdateAdminTeamMember(c *gin.Context) {
	userGroupId, err := uuid.Parse(c.Param("userGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid user group ID",
			},
		})
		return
	}

	var req UpdateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid request format",
			},
		})
		return
	}

	// Extract user ID from request to create entity request
	// Since we don't have user ID in the route, we need to get it from the user group
	err = h.teamMemberService.UpdateTeamMember(c, nil, userGroupId, &entity.UpdateTeamMemberRequest{
		GroupId:  req.GroupId,
		UserData: req.UserData,
	})

	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Updating admin team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin Team member updated successfully",
	})
}

// GetTeamMembers godoc
// @Summary Get team members
// @Description Get list of team members
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /all [get]
func (h *TeamMemberHandler) GetTeamMembers(c *gin.Context) {
	var groupId *uuid.UUID
	var err error

	gId := c.Param("groupId")

	if gId == "all" {
		groupId = nil
	} else {
		id, err := uuid.Parse(c.Param("groupId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "INVALID_REQUEST",
					"message": "Invalid group ID",
				},
			})
			return
		}
		groupId = &id
	}

	merchantId, exists := ginn.GetMerchantIDFromContext(c)

	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	members, err := h.teamMemberService.GetTeamMembers(c, &merchantId, groupId)
	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Getting team members failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    members,
		"message": "Team members fetched successfully",
	})
}

// GetAdminTeamMembers godoc
// @Summary Get admin team members
// @Description Get list of admin team members
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/all [get]
func (h *TeamMemberHandler) GetAdminTeamMembers(c *gin.Context) {
	var groupId *uuid.UUID
	var err error

	gId := c.Param("groupId")

	if gId == "all" {
		groupId = nil
	} else {
		id, err := uuid.Parse(c.Param("groupId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "INVALID_REQUEST",
					"message": "Invalid group ID",
				},
			})
			return
		}
		groupId = &id
	}

	members, err := h.teamMemberService.GetTeamMembers(c, nil, groupId)
	if err != nil {
		log.Println(err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Getting team members failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    members,
		"message": "Team members fetched successfully",
	})
}

// RemoveTeamMember godoc
// @Summary Remove team member
// @Description Remove team member by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /remove/:userGroupId [delete]
func (h *TeamMemberHandler) RemoveTeamMember(c *gin.Context) {
	userGroupId, err := uuid.Parse(c.Param("userGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid user group ID",
			},
		})
		return
	}

	merchantId, exists := ginn.GetMerchantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error: ApiError{
				Type:    "INVALID_REQUEST",
				Message: "Merchant ID not found in context",
			},
		})
		return
	}

	err = h.teamMemberService.RemoveTeamMember(c, &merchantId, userGroupId)

	if err != nil {
		log.Println("Error after service is called -> ", err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Removing team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Team member removed successfully",
	})
}

// RemoveAdminTeamMember godoc
// @Summary Remove admin team member
// @Description Remove admin team member by id
// @Tags Team-member-management
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/remove/:userGroupId [delete]
func (h *TeamMemberHandler) RemoveAdminTeamMember(c *gin.Context) {
	userGroupId, err := uuid.Parse(c.Param("userGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"type":    entity.ErrInvalidRequest,
				"message": "Invalid user group ID",
			},
		})
		return
	}

	err = h.teamMemberService.RemoveTeamMember(c, nil, userGroupId)

	if err != nil {
		log.Println("Error after service is called -> ", err)
		if authErr, ok := err.(*entity.TeamManagementError); ok {
			statusCode := http.StatusBadRequest
			if authErr.Type == entity.ErrInternalServer {
				statusCode = http.StatusInternalServerError
			}
			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"type":    authErr.Type,
					"message": authErr.Message,
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":    entity.ErrInternalServer,
					"message": "Removing admin team member failed",
				},
			})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Admin Team member removed successfully",
	})
}

// RegisterRoutes registers all team management routes
func (h *TeamMemberHandler) RegisterRoutes(router *gin.RouterGroup) {
	team := router.Group("/team-member")

	// Create JWT config
	jwtConfig := ginn.JWTAuthMiddlewareConfig{
		AuthService: h.authService,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Public:      false,
	}

	// Group management routes
	team.POST("/group/create",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_CREATE),
		h.CreateGroup)

	team.GET("/groups",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetGroups)

	team.GET("/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetGroup)

	team.PUT("/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_UPDATE),
		h.UpdateGroup)

	team.DELETE("/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_DELETE),
		h.DeleteGroup)

	// Team member management route
	team.POST("/create",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_CREATE),
		h.CreateTeamMember)

	team.POST("/admin/create",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_CREATE),
		h.CreateAdminTeamMember)

	team.GET("/all/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetTeamMembers)

	team.GET("/admin/all/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetAdminTeamMembers)

	team.PUT("/update/:userGroupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_UPDATE),
		h.UpdateTeamMember)

	team.PUT("/admin/update/:userGroupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_UPDATE),
		h.UpdateAdminTeamMember)

	team.DELETE("/remove/:userGroupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForMerchant(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_DELETE),
		h.RemoveTeamMember)

	team.DELETE("/admin/remove/:userGroupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_DELETE),
		h.RemoveAdminTeamMember)

	// Admin Group management routes
	team.POST("/admin/group/create",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_CREATE),
		h.CreateAdminGroup)

	team.GET("/admin/groups",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetAdminGroups)

	team.GET("/admin/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_READ),
		h.GetAdminGroup)

	team.PUT("/admin/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_UPDATE),
		h.UpdateAdminGroup)

	team.DELETE("/admin/group/:groupId",
		ginn.JWTAuthMiddleware(jwtConfig),
		h.rbac.RequirePermissionForAdmin(auth_entity.RESOURCE_TEAM, auth_entity.OPERATION_DELETE),
		h.DeleteAdminGroup)
}
