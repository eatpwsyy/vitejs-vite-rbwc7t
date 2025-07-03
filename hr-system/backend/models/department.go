package models

import (
	"time"
	"gorm.io/gorm"
)

type Department struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	Name        string `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	
	// Head of Department
	HeadID *uint     `json:"head_id,omitempty"`
	Head   *Employee `json:"head,omitempty" gorm:"foreignKey:HeadID"`
	
	// Budget and Goals
	Budget       float64 `json:"budget" gorm:"default:0"`
	Goals        string  `json:"goals" gorm:"type:text"`
	
	// Relationships
	Employees []Employee `json:"employees,omitempty" gorm:"foreignKey:DepartmentID"`
}

type DepartmentCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	HeadID      *uint   `json:"head_id,omitempty"`
	Budget      float64 `json:"budget"`
	Goals       string  `json:"goals"`
}

type DepartmentUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
	HeadID      *uint    `json:"head_id,omitempty"`
	Budget      *float64 `json:"budget,omitempty"`
	Goals       *string  `json:"goals,omitempty"`
}

type DepartmentResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	Budget       float64   `json:"budget"`
	Goals        string    `json:"goals"`
	Head         *Employee `json:"head,omitempty"`
	EmployeeCount int      `json:"employee_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}