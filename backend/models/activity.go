package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Date          string         `json:"date" binding:"required"` // Format: YYYY-MM-DD
	LearningHours float64        `json:"learning_hours"`
	SleepHours    float64        `json:"sleep_hours"`
	OfficeHours   float64        `json:"office_hours"`
	Notes         string         `json:"notes"`
}

// Made with Bob
