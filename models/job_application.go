package models

import (
	"time"

	"gorm.io/gorm"
)

type JobApplication struct {
	ApplicationID string         `gorm:"type:char(36);primaryKey" json:"application_id"`
	ApplicantID   string         `gorm:"type:char(36);not null;index" json:"-"`    // Foreign key for User who applied
	JobID         string         `gorm:"type:char(36);not null;index" json:"-"`    // Foreign key for the Job applied to
	Status        string         `gorm:"size:255;default:'pending'" json:"status"` // E.g., pending, approved, rejected
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields
}

// Optionally, you can define constants for status values
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)
