package entity

import "github.com/google/uuid"

type ExamLang struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ExamID  string    `gorm:"type:uuid;not null" json:"exam_id"`
	LangID  uint      `gorm:"not null" json:"lang_id"`

	// Relationships
	Exam     Exam          `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
	Language Language      `gorm:"foreignKey:LangID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"language"`
}