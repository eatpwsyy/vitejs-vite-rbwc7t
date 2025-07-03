package handlers

import (
	"net/http"
	"time"

	"hr-backend/config"
	"hr-backend/middleware"
	"hr-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Update last login and refresh token
	now := time.Now()
	user.LastLogin = &now
	user.RefreshToken = refreshToken
	config.DB.Save(&user)

	// Get employee data if exists
	var employee *models.Employee
	if user.EmployeeID != nil {
		var emp models.Employee
		if err := config.DB.Preload("Department").Where("id = ?", *user.EmployeeID).First(&emp).Error; err == nil {
			employee = &emp
		}
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	var employeeResponse *models.EmployeeResponse
	if employee != nil {
		employeeResponse = &models.EmployeeResponse{
			ID:          employee.ID,
			FirstName:   employee.FirstName,
			LastName:    employee.LastName,
			Email:       employee.Email,
			Phone:       employee.Phone,
			EmployeeID:  employee.EmployeeID,
			Position:    employee.Position,
			Department: models.DepartmentResponse{
				ID:   employee.Department.ID,
				Name: employee.Department.Name,
			},
			Status:    employee.Status,
			CreatedAt: employee.CreatedAt,
			UpdatedAt: employee.UpdatedAt,
		}
	}

	response := models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userResponse,
		Employee:     employeeResponse,
	}

	c.JSON(http.StatusOK, response)
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "employee"
	}

	user := models.User{
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       role,
		EmployeeID: req.EmployeeID,
		IsActive:   true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    userResponse,
	})
}

func ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := config.DB.Preload("Employee.Department").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	var employeeResponse *models.EmployeeResponse
	if user.Employee != nil {
		employeeResponse = &models.EmployeeResponse{
			ID:          user.Employee.ID,
			FirstName:   user.Employee.FirstName,
			LastName:    user.Employee.LastName,
			Email:       user.Employee.Email,
			Phone:       user.Employee.Phone,
			EmployeeID:  user.Employee.EmployeeID,
			Position:    user.Employee.Position,
			Department: models.DepartmentResponse{
				ID:   user.Employee.Department.ID,
				Name: user.Employee.Department.Name,
			},
			Status:    user.Employee.Status,
			CreatedAt: user.Employee.CreatedAt,
			UpdatedAt: user.Employee.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user":     userResponse,
		"employee": employeeResponse,
	})
}