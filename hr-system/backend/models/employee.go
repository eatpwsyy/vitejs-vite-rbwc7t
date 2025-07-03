package models

import (
	"time"
	"gorm.io/gorm"
)

type Employee struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Personal Information
	FirstName    string    `json:"first_name" gorm:"not null" validate:"required"`
	LastName     string    `json:"last_name" gorm:"not null" validate:"required"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Phone        string    `json:"phone" validate:"required"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Gender       string    `json:"gender" gorm:"type:enum('Male','Female','Other')"`
	Address      string    `json:"address"`
	
	// Employment Information
	EmployeeID   string    `json:"employee_id" gorm:"uniqueIndex;not null"`
	Position     string    `json:"position" validate:"required"`
	DepartmentID uint      `json:"department_id" validate:"required"`
	Department   Department `json:"department" gorm:"foreignKey:DepartmentID"`
	HireDate     time.Time `json:"hire_date" validate:"required"`
	Salary       float64   `json:"salary" validate:"required,min=0"`
	Status       string    `json:"status" gorm:"default:'Active'" validate:"required"`
	
	// Manager Information
	ManagerID    *uint     `json:"manager_id,omitempty"`
	Manager      *Employee `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	
	// Additional Information
	EmergencyContact     string `json:"emergency_contact"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
	
	// Relationships
	Leaves []Leave `json:"leaves,omitempty" gorm:"foreignKey:EmployeeID"`
}

type EmployeeCreateRequest struct {
	FirstName             string    `json:"first_name" validate:"required"`
	LastName              string    `json:"last_name" validate:"required"`
	Email                 string    `json:"email" validate:"required,email"`
	Phone                 string    `json:"phone" validate:"required"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	Gender                string    `json:"gender"`
	Address               string    `json:"address"`
	Position              string    `json:"position" validate:"required"`
	DepartmentID          uint      `json:"department_id" validate:"required"`
	HireDate              time.Time `json:"hire_date" validate:"required"`
	Salary                float64   `json:"salary" validate:"required,min=0"`
	ManagerID             *uint     `json:"manager_id,omitempty"`
	EmergencyContact      string    `json:"emergency_contact"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
}

type EmployeeUpdateRequest struct {
	FirstName             *string    `json:"first_name,omitempty"`
	LastName              *string    `json:"last_name,omitempty"`
	Email                 *string    `json:"email,omitempty"`
	Phone                 *string    `json:"phone,omitempty"`
	DateOfBirth           *time.Time `json:"date_of_birth,omitempty"`
	Gender                *string    `json:"gender,omitempty"`
	Address               *string    `json:"address,omitempty"`
	Position              *string    `json:"position,omitempty"`
	DepartmentID          *uint      `json:"department_id,omitempty"`
	Salary                *float64   `json:"salary,omitempty"`
	Status                *string    `json:"status,omitempty"`
	ManagerID             *uint      `json:"manager_id,omitempty"`
	EmergencyContact      *string    `json:"emergency_contact,omitempty"`
	EmergencyContactPhone *string    `json:"emergency_contact_phone,omitempty"`
}

type EmployeeResponse struct {
	ID                    uint               `json:"id"`
	FirstName             string             `json:"first_name"`
	LastName              string             `json:"last_name"`
	Email                 string             `json:"email"`
	Phone                 string             `json:"phone"`
	DateOfBirth           time.Time          `json:"date_of_birth"`
	Gender                string             `json:"gender"`
	Address               string             `json:"address"`
	EmployeeID            string             `json:"employee_id"`
	Position              string             `json:"position"`
	Department            DepartmentResponse `json:"department"`
	HireDate              time.Time          `json:"hire_date"`
	Salary                float64            `json:"salary"`
	Status                string             `json:"status"`
	Manager               *EmployeeResponse  `json:"manager,omitempty"`
	EmergencyContact      string             `json:"emergency_contact"`
	EmergencyContactPhone string             `json:"emergency_contact_phone"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}