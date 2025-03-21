package entity

import "github.com/google/uuid"

type UserExam struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID string `gorm:"type:uuid;not null" json:"user_id"`
	ExamID string `gorm:"type:uuid;not null" json:"exam_id"`
	Role   string    `json:"role"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Exam Exam `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exam"`
}
