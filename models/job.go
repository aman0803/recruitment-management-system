package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	JobID       string         `gorm:"type:char(36);primaryKey" json:"job_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	CompanyName string         `gorm:"size:255;not null" json:"company_name"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	PostedByID  string         `gorm:"type:char(36);not null" json:"-"`                  // Foreign key
	PostedBy    User           `gorm:"foreignKey:PostedByID;references:UserID" json:"-"` // Reference to User
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
