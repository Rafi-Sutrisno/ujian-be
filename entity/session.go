package entity

import (
	"time"
)

type DatabaseData struct {
	Id        string                 `gorm:"primaryKey"`
	Data      map[string]interface{} `gorm:"serializer:json"` // Stores session data as JSON
	ExpiredAt time.Time              `gorm:"index"`
	CSRFToken string                 // Store CSRF token per session
}
