package models

import (
	"time"
)

type File struct {
	ID          string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Tag         string     `gorm:"type:varchar(255)" json:"tag"`
	Filename    string     `gorm:"type:varchar(255)" json:"filename"`
	ContentType string     `gorm:"type:varchar(255)" json:"content_type"`
	Size        int        `gorm:"type:int" json:"size"`
	Deleted     *bool      `gorm:"type:boolean;default:false" json:"deleted,omitempty"`
	CreatedAt   time.Time  `gorm:"default:current_timestamp" json:"created_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
