package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:employee"` // employee, admin
	Position  string         `json:"position"`
	Department string        `json:"department"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Attendance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Date      time.Time      `json:"date" gorm:"type:date;not null"`
	CheckIn   *time.Time     `json:"check_in"`
	CheckOut  *time.Time     `json:"check_out"`
	Notes     string         `json:"notes"`
	Status    string         `json:"status" gorm:"default:present"` // present, absent, late, half_day
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	Position   string `json:"position"`
	Department string `json:"department"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CheckInRequest struct {
	Notes string `json:"notes"`
}

type CheckOutRequest struct {
	Notes string `json:"notes"`
}

type AttendanceStats struct {
	TotalDays     int     `json:"total_days"`
	PresentDays   int     `json:"present_days"`
	AbsentDays    int     `json:"absent_days"`
	LateDays      int     `json:"late_days"`
	AttendanceRate float64 `json:"attendance_rate"`
}

type UpdateUserRequest struct {
	Name       string `json:"name"`
	Position   string `json:"position"`
	Department string `json:"department"`
	Role       string `json:"role"`
}

type UpdateProfileRequest struct {
	Name       string `json:"name"`
	Position   string `json:"position"`
	Department string `json:"department"`
}