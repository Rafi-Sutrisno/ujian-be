package entity

import (
	"github.com/google/uuid"
)

type Class struct {
    ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    Name        string        `json:"name"`
    ShortName   string        `json:"short_name"`           
	Timestamp

	// Relationships
    UserClass []UserClass `gorm:"foreignKey:ClassID" json:"user_class"`
	Exams     []Exam      `gorm:"foreignKey:ClassID" json:"exams"`
}
