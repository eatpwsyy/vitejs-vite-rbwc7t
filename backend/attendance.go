package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func checkIn(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	
	// Check if user already checked in today
	var existingAttendance Attendance
	if err := db.Where("user_id = ? AND date = ?", userID, today).First(&existingAttendance).Error; err == nil {
		if existingAttendance.CheckIn != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked in today"})
			return
		}
	}

	now := time.Now()
	
	// Determine status based on check-in time (assuming work starts at 9:00 AM)
	workStartTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	status := "present"
	if now.After(workStartTime.Add(15 * time.Minute)) { // 15 minutes grace period
		status = "late"
	}

	if existingAttendance.ID != 0 {
		// Update existing record
		existingAttendance.CheckIn = &now
		existingAttendance.Status = status
		existingAttendance.Notes = req.Notes
		
		if err := db.Save(&existingAttendance).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update check-in"})
			return
		}
		
		c.JSON(http.StatusOK, existingAttendance)
	} else {
		// Create new attendance record
		attendance := Attendance{
			UserID:  userID.(uint),
			Date:    today,
			CheckIn: &now,
			Status:  status,
			Notes:   req.Notes,
		}

		if err := db.Create(&attendance).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record check-in"})
			return
		}

		// Load user data
		if err := db.Preload("User").First(&attendance, attendance.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load attendance data"})
			return
		}

		c.JSON(http.StatusCreated, attendance)
	}
}

func checkOut(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	
	// Find today's attendance record
	var attendance Attendance
	if err := db.Where("user_id = ? AND date = ?", userID, today).First(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No check-in found for today"})
		return
	}

	if attendance.CheckIn == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must check in before checking out"})
		return
	}

	if attendance.CheckOut != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked out today"})
		return
	}

	now := time.Now()
	attendance.CheckOut = &now
	
	// Update notes if provided
	if req.Notes != "" {
		if attendance.Notes != "" {
			attendance.Notes += " | Checkout: " + req.Notes
		} else {
			attendance.Notes = "Checkout: " + req.Notes
		}
	}

	// Calculate work duration
	workDuration := now.Sub(*attendance.CheckIn)
	
	// Determine if it's half day (less than 4 hours)
	if workDuration.Hours() < 4 && attendance.Status != "late" {
		attendance.Status = "half_day"
	}

	if err := db.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record check-out"})
		return
	}

	// Load user data
	if err := db.Preload("User").First(&attendance, attendance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load attendance data"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func getTodayAttendance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	today := time.Now().Truncate(24 * time.Hour)

	var attendance Attendance
	err := db.Preload("User").Where("user_id = ? AND date = ?", userID, today).First(&attendance).Error
	
	if err != nil {
		// No attendance record for today
		c.JSON(http.StatusOK, gin.H{
			"date": today,
			"status": "not_checked_in",
			"check_in": nil,
			"check_out": nil,
		})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func getAttendanceHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Get query parameters
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

	var attendances []Attendance
	var total int64

	// Count total records
	db.Model(&Attendance{}).Where("user_id = ?", userID).Count(&total)

	// Get paginated records
	if err := db.Preload("User").
		Where("user_id = ?", userID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance history"})
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

func getAttendanceStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Get query parameters for date range (default to current month)
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, -1) // Last day of current month

	if start := c.Query("start_date"); start != "" {
		if parsed, err := time.Parse("2006-01-02", start); err == nil {
			startDate = parsed
		}
	}
	
	if end := c.Query("end_date"); end != "" {
		if parsed, err := time.Parse("2006-01-02", end); err == nil {
			endDate = parsed
		}
	}

	var stats AttendanceStats

	// Count total working days in the period (excluding weekends)
	totalDays := 0
	for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			totalDays++
		}
	}

	// Count present days
	var presentCount int64
	db.Model(&Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ? AND status = ?", userID, startDate, endDate, "present").
		Count(&presentCount)

	// Count late days
	var lateCount int64
	db.Model(&Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ? AND status = ?", userID, startDate, endDate, "late").
		Count(&lateCount)

	// Count half days
	var halfDayCount int64
	db.Model(&Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ? AND status = ?", userID, startDate, endDate, "half_day").
		Count(&halfDayCount)

	// Count total attendance records
	var totalAttendance int64
	db.Model(&Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ? AND check_in IS NOT NULL", userID, startDate, endDate).
		Count(&totalAttendance)

	stats.TotalDays = totalDays
	stats.PresentDays = int(presentCount)
	stats.LateDays = int(lateCount)
	stats.AbsentDays = totalDays - int(totalAttendance)
	
	if totalDays > 0 {
		stats.AttendanceRate = float64(totalAttendance) / float64(totalDays) * 100
	}

	c.JSON(http.StatusOK, stats)
}