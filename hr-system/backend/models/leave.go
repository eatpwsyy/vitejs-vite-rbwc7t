package models

import (
	"time"
	"gorm.io/gorm"
)

type Leave struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Employee Information
	EmployeeID  uint     `json:"employee_id" validate:"required"`
	Employee    Employee `json:"employee" gorm:"foreignKey:EmployeeID"`
	
	// Leave Details
	Type        string    `json:"type" validate:"required"` // Annual, Sick, Maternity, Personal, etc.
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Days        int       `json:"days" validate:"required,min=1"`
	Reason      string    `json:"reason" validate:"required"`
	Status      string    `json:"status" gorm:"default:'Pending'" validate:"required"` // Pending, Approved, Rejected
	
	// Approval Information
	ApprovedBy  *uint     `json:"approved_by,omitempty"`
	Approver    *Employee `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	Comments    string    `json:"comments"`
	
	// Documentation
	AttachmentURL string `json:"attachment_url,omitempty"`
}

type LeaveCreateRequest struct {
	EmployeeID    uint      `json:"employee_id" validate:"required"`
	Type          string    `json:"type" validate:"required"`
	StartDate     time.Time `json:"start_date" validate:"required"`
	EndDate       time.Time `json:"end_date" validate:"required"`
	Reason        string    `json:"reason" validate:"required"`
	AttachmentURL string    `json:"attachment_url,omitempty"`
}

type LeaveUpdateRequest struct {
	Type          *string    `json:"type,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Reason        *string    `json:"reason,omitempty"`
	Status        *string    `json:"status,omitempty"`
	Comments      *string    `json:"comments,omitempty"`
	AttachmentURL *string    `json:"attachment_url,omitempty"`
}

type LeaveApprovalRequest struct {
	Status   string `json:"status" validate:"required,oneof=Approved Rejected"`
	Comments string `json:"comments"`
}

type LeaveResponse struct {
	ID            uint              `json:"id"`
	Employee      EmployeeResponse  `json:"employee"`
	Type          string            `json:"type"`
	StartDate     time.Time         `json:"start_date"`
	EndDate       time.Time         `json:"end_date"`
	Days          int               `json:"days"`
	Reason        string            `json:"reason"`
	Status        string            `json:"status"`
	Approver      *EmployeeResponse `json:"approver,omitempty"`
	ApprovedAt    *time.Time        `json:"approved_at,omitempty"`
	Comments      string            `json:"comments"`
	AttachmentURL string            `json:"attachment_url,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// Leave Types
var LeaveTypes = []string{
	"Annual",
	"Sick",
	"Maternity",
	"Paternity",
	"Personal",
	"Emergency",
	"Bereavement",
	"Study",
	"Unpaid",
}

// Leave Statuses
var LeaveStatuses = []string{
	"Pending",
	"Approved",
	"Rejected",
	"Cancelled",
}