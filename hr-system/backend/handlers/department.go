package handlers

import (
	"net/http"
	"strconv"

	"hr-backend/config"
	"hr-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDepartments(c *gin.Context) {
	var departments []models.Department
	query := config.DB.Preload("Head")

	// Add filters
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if active := c.Query("active"); active != "" {
		if isActive, err := strconv.ParseBool(active); err == nil {
			query = query.Where("is_active = ?", isActive)
		}
	}

	if err := query.Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}

	// Convert to response format and get employee count
	var departmentResponses []models.DepartmentResponse
	for _, dept := range departments {
		var employeeCount int64
		config.DB.Model(&models.Employee{}).Where("department_id = ?", dept.ID).Count(&employeeCount)

		var head *models.Employee
		if dept.Head != nil {
			head = dept.Head
		}

		departmentResponses = append(departmentResponses, models.DepartmentResponse{
			ID:            dept.ID,
			Name:          dept.Name,
			Description:   dept.Description,
			IsActive:      dept.IsActive,
			Budget:        dept.Budget,
			Goals:         dept.Goals,
			Head:          head,
			EmployeeCount: int(employeeCount),
			CreatedAt:     dept.CreatedAt,
			UpdatedAt:     dept.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"departments": departmentResponses})
}

func GetDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.Preload("Head").Preload("Employees").Where("id = ?", id).First(&department).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department"})
		return
	}

	var head *models.Employee
	if department.Head != nil {
		head = department.Head
	}

	response := models.DepartmentResponse{
		ID:            department.ID,
		Name:          department.Name,
		Description:   department.Description,
		IsActive:      department.IsActive,
		Budget:        department.Budget,
		Goals:         department.Goals,
		Head:          head,
		EmployeeCount: len(department.Employees),
		CreatedAt:     department.CreatedAt,
		UpdatedAt:     department.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func CreateDepartment(c *gin.Context) {
	var req models.DepartmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if department name already exists
	var existingDept models.Department
	if err := config.DB.Where("name = ?", req.Name).First(&existingDept).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Department with this name already exists"})
		return
	}

	department := models.Department{
		Name:        req.Name,
		Description: req.Description,
		HeadID:      req.HeadID,
		Budget:      req.Budget,
		Goals:       req.Goals,
		IsActive:    true,
	}

	if err := config.DB.Create(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	// Load relationships
	config.DB.Preload("Head").First(&department, department.ID)

	var head *models.Employee
	if department.Head != nil {
		head = department.Head
	}

	response := models.DepartmentResponse{
		ID:            department.ID,
		Name:          department.Name,
		Description:   department.Description,
		IsActive:      department.IsActive,
		Budget:        department.Budget,
		Goals:         department.Goals,
		Head:          head,
		EmployeeCount: 0,
		CreatedAt:     department.CreatedAt,
		UpdatedAt:     department.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.Where("id = ?", id).First(&department).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department"})
		return
	}

	var req models.DepartmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if department name already exists for another department
		var existingDept models.Department
		if err := config.DB.Where("name = ? AND id != ?", *req.Name, id).First(&existingDept).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Department with this name already exists"})
			return
		}
		department.Name = *req.Name
	}
	if req.Description != nil {
		department.Description = *req.Description
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}
	if req.HeadID != nil {
		department.HeadID = req.HeadID
	}
	if req.Budget != nil {
		department.Budget = *req.Budget
	}
	if req.Goals != nil {
		department.Goals = *req.Goals
	}

	if err := config.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	// Load relationships
	config.DB.Preload("Head").First(&department, department.ID)

	var employeeCount int64
	config.DB.Model(&models.Employee{}).Where("department_id = ?", department.ID).Count(&employeeCount)

	var head *models.Employee
	if department.Head != nil {
		head = department.Head
	}

	response := models.DepartmentResponse{
		ID:            department.ID,
		Name:          department.Name,
		Description:   department.Description,
		IsActive:      department.IsActive,
		Budget:        department.Budget,
		Goals:         department.Goals,
		Head:          head,
		EmployeeCount: int(employeeCount),
		CreatedAt:     department.CreatedAt,
		UpdatedAt:     department.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.Where("id = ?", id).First(&department).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department"})
		return
	}

	// Check if department has employees
	var employeeCount int64
	config.DB.Model(&models.Employee{}).Where("department_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete department with existing employees"})
		return
	}

	if err := config.DB.Delete(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
}