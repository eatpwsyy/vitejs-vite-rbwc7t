package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	Email      string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password   string `json:"-" gorm:"not null" validate:"required,min=6"` // Hidden in JSON
	Role       string `json:"role" gorm:"default:'employee'" validate:"required"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`
	
	// Optional employee reference
	EmployeeID *uint     `json:"employee_id,omitempty"`
	Employee   *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	
	// Session management
	LastLogin    *time.Time `json:"last_login,omitempty"`
	RefreshToken string     `json:"-" gorm:"type:text"` // Hidden in JSON
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	Role       string `json:"role,omitempty"`
	EmployeeID *uint  `json:"employee_id,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token        string           `json:"token"`
	RefreshToken string           `json:"refresh_token"`
	User         UserResponse     `json:"user"`
	Employee     *EmployeeResponse `json:"employee,omitempty"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User Roles
var UserRoles = []string{
	"admin",
	"hr",
	"manager",
	"employee",
}