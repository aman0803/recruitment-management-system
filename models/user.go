package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID          string           `gorm:"type:char(36);primaryKey" json:"user_id"`
	Name            string           `gorm:"size:255;not null" json:"name"`
	Email           string           `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash    string           `gorm:"size:255;not null" json:"password_hash"`
	IsEmployer      bool             `gorm:"default:false" json:"is_employer"`
	ProfileHeadline string           `gorm:"size:255" json:"profile_headline"`
	Address         string           `gorm:"type:text" json:"address"`
	Resumes         []Resume         `gorm:"foreignKey:UserID" json:"-"`           // Relationship to Resumes
	Applications    []JobApplication `gorm:"foreignKey:ApplicantID" json:"-"` // Relationship to Job Applications
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
}
