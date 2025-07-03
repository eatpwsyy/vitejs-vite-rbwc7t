package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hr-backend/config"
	"hr-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	query := config.DB.Preload("Department").Preload("Manager")

	// Add filters
	if dept := c.Query("department"); dept != "" {
		if deptID, err := strconv.Atoi(dept); err == nil {
			query = query.Where("department_id = ?", deptID)
		}
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR employee_id ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Model(&models.Employee{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}

	// Convert to response format
	var employeeResponses []models.EmployeeResponse
	for _, emp := range employees {
		deptResponse := models.DepartmentResponse{
			ID:          emp.Department.ID,
			Name:        emp.Department.Name,
			Description: emp.Department.Description,
			IsActive:    emp.Department.IsActive,
		}

		var managerResponse *models.EmployeeResponse
		if emp.Manager != nil {
			managerResponse = &models.EmployeeResponse{
				ID:         emp.Manager.ID,
				FirstName:  emp.Manager.FirstName,
				LastName:   emp.Manager.LastName,
				Email:      emp.Manager.Email,
				EmployeeID: emp.Manager.EmployeeID,
				Position:   emp.Manager.Position,
			}
		}

		employeeResponses = append(employeeResponses, models.EmployeeResponse{
			ID:                    emp.ID,
			FirstName:             emp.FirstName,
			LastName:              emp.LastName,
			Email:                 emp.Email,
			Phone:                 emp.Phone,
			DateOfBirth:           emp.DateOfBirth,
			Gender:                emp.Gender,
			Address:               emp.Address,
			EmployeeID:            emp.EmployeeID,
			Position:              emp.Position,
			Department:            deptResponse,
			HireDate:              emp.HireDate,
			Salary:                emp.Salary,
			Status:                emp.Status,
			Manager:               managerResponse,
			EmergencyContact:      emp.EmergencyContact,
			EmergencyContactPhone: emp.EmergencyContactPhone,
			CreatedAt:             emp.CreatedAt,
			UpdatedAt:             emp.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employeeResponses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.Preload("Department").Preload("Manager").Preload("Leaves").Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	deptResponse := models.DepartmentResponse{
		ID:          employee.Department.ID,
		Name:        employee.Department.Name,
		Description: employee.Department.Description,
		IsActive:    employee.Department.IsActive,
	}

	var managerResponse *models.EmployeeResponse
	if employee.Manager != nil {
		managerResponse = &models.EmployeeResponse{
			ID:         employee.Manager.ID,
			FirstName:  employee.Manager.FirstName,
			LastName:   employee.Manager.LastName,
			Email:      employee.Manager.Email,
			EmployeeID: employee.Manager.EmployeeID,
			Position:   employee.Manager.Position,
		}
	}

	response := models.EmployeeResponse{
		ID:                    employee.ID,
		FirstName:             employee.FirstName,
		LastName:              employee.LastName,
		Email:                 employee.Email,
		Phone:                 employee.Phone,
		DateOfBirth:           employee.DateOfBirth,
		Gender:                employee.Gender,
		Address:               employee.Address,
		EmployeeID:            employee.EmployeeID,
		Position:              employee.Position,
		Department:            deptResponse,
		HireDate:              employee.HireDate,
		Salary:                employee.Salary,
		Status:                employee.Status,
		Manager:               managerResponse,
		EmergencyContact:      employee.EmergencyContact,
		EmergencyContactPhone: employee.EmergencyContactPhone,
		CreatedAt:             employee.CreatedAt,
		UpdatedAt:             employee.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func CreateEmployee(c *gin.Context) {
	var req models.EmployeeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingEmployee models.Employee
	if err := config.DB.Where("email = ?", req.Email).First(&existingEmployee).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Employee with this email already exists"})
		return
	}

	// Generate employee ID if not provided
	employeeID := generateEmployeeID()

	employee := models.Employee{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Email:                 req.Email,
		Phone:                 req.Phone,
		DateOfBirth:           req.DateOfBirth,
		Gender:                req.Gender,
		Address:               req.Address,
		EmployeeID:            employeeID,
		Position:              req.Position,
		DepartmentID:          req.DepartmentID,
		HireDate:              req.HireDate,
		Salary:                req.Salary,
		Status:                "Active",
		ManagerID:             req.ManagerID,
		EmergencyContact:      req.EmergencyContact,
		EmergencyContactPhone: req.EmergencyContactPhone,
	}

	if err := config.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	// Load relationships
	config.DB.Preload("Department").Preload("Manager").First(&employee, employee.ID)

	deptResponse := models.DepartmentResponse{
		ID:          employee.Department.ID,
		Name:        employee.Department.Name,
		Description: employee.Department.Description,
		IsActive:    employee.Department.IsActive,
	}

	response := models.EmployeeResponse{
		ID:                    employee.ID,
		FirstName:             employee.FirstName,
		LastName:              employee.LastName,
		Email:                 employee.Email,
		Phone:                 employee.Phone,
		DateOfBirth:           employee.DateOfBirth,
		Gender:                employee.Gender,
		Address:               employee.Address,
		EmployeeID:            employee.EmployeeID,
		Position:              employee.Position,
		Department:            deptResponse,
		HireDate:              employee.HireDate,
		Salary:                employee.Salary,
		Status:                employee.Status,
		EmergencyContact:      employee.EmergencyContact,
		EmergencyContactPhone: employee.EmergencyContactPhone,
		CreatedAt:             employee.CreatedAt,
		UpdatedAt:             employee.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	var req models.EmployeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Email != nil {
		// Check if email already exists for another employee
		var existingEmployee models.Employee
		if err := config.DB.Where("email = ? AND id != ?", *req.Email, id).First(&existingEmployee).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Employee with this email already exists"})
			return
		}
		employee.Email = *req.Email
	}
	if req.Phone != nil {
		employee.Phone = *req.Phone
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = *req.DateOfBirth
	}
	if req.Gender != nil {
		employee.Gender = *req.Gender
	}
	if req.Address != nil {
		employee.Address = *req.Address
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.DepartmentID != nil {
		employee.DepartmentID = *req.DepartmentID
	}
	if req.Salary != nil {
		employee.Salary = *req.Salary
	}
	if req.Status != nil {
		employee.Status = *req.Status
	}
	if req.ManagerID != nil {
		employee.ManagerID = req.ManagerID
	}
	if req.EmergencyContact != nil {
		employee.EmergencyContact = *req.EmergencyContact
	}
	if req.EmergencyContactPhone != nil {
		employee.EmergencyContactPhone = *req.EmergencyContactPhone
	}

	if err := config.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	// Load relationships
	config.DB.Preload("Department").Preload("Manager").First(&employee, employee.ID)

	deptResponse := models.DepartmentResponse{
		ID:          employee.Department.ID,
		Name:        employee.Department.Name,
		Description: employee.Department.Description,
		IsActive:    employee.Department.IsActive,
	}

	var managerResponse *models.EmployeeResponse
	if employee.Manager != nil {
		managerResponse = &models.EmployeeResponse{
			ID:         employee.Manager.ID,
			FirstName:  employee.Manager.FirstName,
			LastName:   employee.Manager.LastName,
			Email:      employee.Manager.Email,
			EmployeeID: employee.Manager.EmployeeID,
			Position:   employee.Manager.Position,
		}
	}

	response := models.EmployeeResponse{
		ID:                    employee.ID,
		FirstName:             employee.FirstName,
		LastName:              employee.LastName,
		Email:                 employee.Email,
		Phone:                 employee.Phone,
		DateOfBirth:           employee.DateOfBirth,
		Gender:                employee.Gender,
		Address:               employee.Address,
		EmployeeID:            employee.EmployeeID,
		Position:              employee.Position,
		Department:            deptResponse,
		HireDate:              employee.HireDate,
		Salary:                employee.Salary,
		Status:                employee.Status,
		Manager:               managerResponse,
		EmergencyContact:      employee.EmergencyContact,
		EmergencyContactPhone: employee.EmergencyContactPhone,
		CreatedAt:             employee.CreatedAt,
		UpdatedAt:             employee.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	if err := config.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

func generateEmployeeID() string {
	var count int64
	config.DB.Model(&models.Employee{}).Count(&count)
	return fmt.Sprintf("EMP%04d", count+1)
}