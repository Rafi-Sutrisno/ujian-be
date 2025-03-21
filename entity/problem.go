package entity

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    ExamID       uuid.UUID `gorm:"type:uuid;not null" json:"exam_id"`
    Title        string    `gorm:"type:varchar(255);not null" json:"title"`
    Description  string    `gorm:"type:text" json:"description"`     // Content of the problem
    Constraints  string    `gorm:"type:varchar(255)" json:"constraints"`
    SampleInput  string    `gorm:"type:varchar(255)" json:"sample_input"`
    SampleOutput string    `gorm:"type:varchar(255)" json:"sample_output"`
    CreatedBy    string    `gorm:"type:varchar(255)" json:"created_by"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

    // Relationships
    Exam Exam `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"exam"`
}
