package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAllUsers(c *gin.Context) {
	page := 1
	limit := 10
	
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	var users []User
	var total int64

	// Count total users
	db.Model(&User{}).Count(&total)

	// Get paginated users
	if err := db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func getAllAttendance(c *gin.Context) {
	page := 1
	limit := 10
	
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Filter by user ID if provided
	userIDParam := c.Query("user_id")
	
	// Filter by date range if provided
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	offset := (page - 1) * limit

	var attendances []Attendance
	var total int64

	query := db.Model(&Attendance{})

	// Apply filters
	if userIDParam != "" {
		if userID, err := strconv.ParseUint(userIDParam, 10, 32); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		query = query.Where("date >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	// Count total records with filters applied
	query.Count(&total)

	// Get paginated records with user data
	if err := query.Preload("User").
		Order("date DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": attendances,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func updateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Position != "" {
		user.Position = req.Position
	}
	if req.Department != "" {
		user.Department = req.Department
	}
	if req.Role != "" {
		// Validate role
		if req.Role == "admin" || req.Role == "employee" {
			user.Role = req.Role
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be 'admin' or 'employee'"})
			return
		}
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	// Check if user exists
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent deletion of admin users (optional safety check)
	if user.Role == "admin" {
		var adminCount int64
		db.Model(&User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the last admin user"})
			return
		}
	}

	// Soft delete the user (GORM will set deleted_at timestamp)
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}