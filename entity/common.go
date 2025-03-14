package entity

import (
	"time"
)

type Timestamp struct {
	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

type Authorization struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}