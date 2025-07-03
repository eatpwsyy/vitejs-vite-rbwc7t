package handlers

import (
	"net/http"
	"strconv"
	"time"

	"hr-backend/config"
	"hr-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLeaves(c *gin.Context) {
	var leaves []models.Leave
	query := config.DB.Preload("Employee.Department").Preload("Approver")

	// Add filters
	if empID := c.Query("employee"); empID != "" {
		if employeeID, err := strconv.Atoi(empID); err == nil {
			query = query.Where("employee_id = ?", employeeID)
		}
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if leaveType := c.Query("type"); leaveType != "" {
		query = query.Where("type = ?", leaveType)
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("start_date >= ?", date)
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("end_date <= ?", date)
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Model(&models.Leave{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&leaves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaves"})
		return
	}

	// Convert to response format
	var leaveResponses []models.LeaveResponse
	for _, leave := range leaves {
		empResponse := models.EmployeeResponse{
			ID:         leave.Employee.ID,
			FirstName:  leave.Employee.FirstName,
			LastName:   leave.Employee.LastName,
			Email:      leave.Employee.Email,
			EmployeeID: leave.Employee.EmployeeID,
			Position:   leave.Employee.Position,
			Department: models.DepartmentResponse{
				ID:   leave.Employee.Department.ID,
				Name: leave.Employee.Department.Name,
			},
		}

		var approverResponse *models.EmployeeResponse
		if leave.Approver != nil {
			approverResponse = &models.EmployeeResponse{
				ID:         leave.Approver.ID,
				FirstName:  leave.Approver.FirstName,
				LastName:   leave.Approver.LastName,
				Email:      leave.Approver.Email,
				EmployeeID: leave.Approver.EmployeeID,
				Position:   leave.Approver.Position,
			}
		}

		leaveResponses = append(leaveResponses, models.LeaveResponse{
			ID:            leave.ID,
			Employee:      empResponse,
			Type:          leave.Type,
			StartDate:     leave.StartDate,
			EndDate:       leave.EndDate,
			Days:          leave.Days,
			Reason:        leave.Reason,
			Status:        leave.Status,
			Approver:      approverResponse,
			ApprovedAt:    leave.ApprovedAt,
			Comments:      leave.Comments,
			AttachmentURL: leave.AttachmentURL,
			CreatedAt:     leave.CreatedAt,
			UpdatedAt:     leave.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"leaves": leaveResponses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetLeave(c *gin.Context) {
	id := c.Param("id")
	var leave models.Leave

	if err := config.DB.Preload("Employee.Department").Preload("Approver").Where("id = ?", id).First(&leave).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leave"})
		return
	}

	empResponse := models.EmployeeResponse{
		ID:         leave.Employee.ID,
		FirstName:  leave.Employee.FirstName,
		LastName:   leave.Employee.LastName,
		Email:      leave.Employee.Email,
		EmployeeID: leave.Employee.EmployeeID,
		Position:   leave.Employee.Position,
		Department: models.DepartmentResponse{
			ID:   leave.Employee.Department.ID,
			Name: leave.Employee.Department.Name,
		},
	}

	var approverResponse *models.EmployeeResponse
	if leave.Approver != nil {
		approverResponse = &models.EmployeeResponse{
			ID:         leave.Approver.ID,
			FirstName:  leave.Approver.FirstName,
			LastName:   leave.Approver.LastName,
			Email:      leave.Approver.Email,
			EmployeeID: leave.Approver.EmployeeID,
			Position:   leave.Approver.Position,
		}
	}

	response := models.LeaveResponse{
		ID:            leave.ID,
		Employee:      empResponse,
		Type:          leave.Type,
		StartDate:     leave.StartDate,
		EndDate:       leave.EndDate,
		Days:          leave.Days,
		Reason:        leave.Reason,
		Status:        leave.Status,
		Approver:      approverResponse,
		ApprovedAt:    leave.ApprovedAt,
		Comments:      leave.Comments,
		AttachmentURL: leave.AttachmentURL,
		CreatedAt:     leave.CreatedAt,
		UpdatedAt:     leave.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func CreateLeave(c *gin.Context) {
	var req models.LeaveCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
		return
	}

	// Calculate days
	days := int(req.EndDate.Sub(req.StartDate).Hours()/24) + 1

	// Check for overlapping leaves
	var existingLeave models.Leave
	if err := config.DB.Where("employee_id = ? AND status IN ? AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?))",
		req.EmployeeID, []string{"Pending", "Approved"}, req.StartDate, req.StartDate, req.EndDate, req.EndDate).First(&existingLeave).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Employee has overlapping leave during this period"})
		return
	}

	leave := models.Leave{
		EmployeeID:    req.EmployeeID,
		Type:          req.Type,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Days:          days,
		Reason:        req.Reason,
		Status:        "Pending",
		AttachmentURL: req.AttachmentURL,
	}

	if err := config.DB.Create(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leave"})
		return
	}

	// Load relationships
	config.DB.Preload("Employee.Department").First(&leave, leave.ID)

	empResponse := models.EmployeeResponse{
		ID:         leave.Employee.ID,
		FirstName:  leave.Employee.FirstName,
		LastName:   leave.Employee.LastName,
		Email:      leave.Employee.Email,
		EmployeeID: leave.Employee.EmployeeID,
		Position:   leave.Employee.Position,
		Department: models.DepartmentResponse{
			ID:   leave.Employee.Department.ID,
			Name: leave.Employee.Department.Name,
		},
	}

	response := models.LeaveResponse{
		ID:            leave.ID,
		Employee:      empResponse,
		Type:          leave.Type,
		StartDate:     leave.StartDate,
		EndDate:       leave.EndDate,
		Days:          leave.Days,
		Reason:        leave.Reason,
		Status:        leave.Status,
		AttachmentURL: leave.AttachmentURL,
		CreatedAt:     leave.CreatedAt,
		UpdatedAt:     leave.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateLeave(c *gin.Context) {
	id := c.Param("id")
	var leave models.Leave

	if err := config.DB.Where("id = ?", id).First(&leave).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leave"})
		return
	}

	// Check if leave can be updated (only pending leaves)
	if leave.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update leave that is not pending"})
		return
	}

	var req models.LeaveUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Type != nil {
		leave.Type = *req.Type
	}
	if req.StartDate != nil {
		leave.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		leave.EndDate = *req.EndDate
	}
	if req.Reason != nil {
		leave.Reason = *req.Reason
	}
	if req.AttachmentURL != nil {
		leave.AttachmentURL = *req.AttachmentURL
	}

	// Recalculate days if dates changed
	if req.StartDate != nil || req.EndDate != nil {
		if leave.EndDate.Before(leave.StartDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
			return
		}
		leave.Days = int(leave.EndDate.Sub(leave.StartDate).Hours()/24) + 1
	}

	if err := config.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave"})
		return
	}

	// Load relationships
	config.DB.Preload("Employee.Department").Preload("Approver").First(&leave, leave.ID)

	empResponse := models.EmployeeResponse{
		ID:         leave.Employee.ID,
		FirstName:  leave.Employee.FirstName,
		LastName:   leave.Employee.LastName,
		Email:      leave.Employee.Email,
		EmployeeID: leave.Employee.EmployeeID,
		Position:   leave.Employee.Position,
		Department: models.DepartmentResponse{
			ID:   leave.Employee.Department.ID,
			Name: leave.Employee.Department.Name,
		},
	}

	response := models.LeaveResponse{
		ID:            leave.ID,
		Employee:      empResponse,
		Type:          leave.Type,
		StartDate:     leave.StartDate,
		EndDate:       leave.EndDate,
		Days:          leave.Days,
		Reason:        leave.Reason,
		Status:        leave.Status,
		AttachmentURL: leave.AttachmentURL,
		CreatedAt:     leave.CreatedAt,
		UpdatedAt:     leave.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func ApproveLeave(c *gin.Context) {
	id := c.Param("id")
	var leave models.Leave

	if err := config.DB.Where("id = ?", id).First(&leave).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leave"})
		return
	}

	if leave.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave is not pending approval"})
		return
	}

	var req models.LeaveApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	approverID := c.GetUint("user_id")
	now := time.Now()

	leave.Status = req.Status
	leave.Comments = req.Comments
	leave.ApprovedBy = &approverID
	leave.ApprovedAt = &now

	if err := config.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave " + req.Status + " successfully"})
}

func DeleteLeave(c *gin.Context) {
	id := c.Param("id")
	var leave models.Leave

	if err := config.DB.Where("id = ?", id).First(&leave).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leave"})
		return
	}

	// Check if leave can be deleted (only pending leaves)
	if leave.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete leave that is not pending"})
		return
	}

	if err := config.DB.Delete(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete leave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave deleted successfully"})
}

func GetLeaveTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"leave_types": models.LeaveTypes})
}

func GetLeaveStatuses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"leave_statuses": models.LeaveStatuses})
}