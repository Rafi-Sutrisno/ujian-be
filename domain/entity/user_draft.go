package entity

import "github.com/google/uuid"

type UserCodeDraft struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	ExamID    string    `gorm:"not null" json:"exam_id"`
	ProblemID string    `gorm:"not null" json:"problem_id"`
	Language  string    `gorm:"not null" json:"language"`
	Code      string    `json:"code"`
}
