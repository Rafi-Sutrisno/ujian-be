package entity

import (
	"github.com/google/uuid"
)

type Problem struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ExamID        string    `gorm:"type:uuid;not null" json:"exam_id"`
	Title         string    `json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Constraints   string    `json:"constraints"`
	SampleInput   string    `json:"sample_input"`
	SampleOutput  string    `json:"sample_output"`

    Timestamp

    // Relationships
    Exam            Exam `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
    TestCase        []TestCase `gorm:"foreignKey:ProblemID" json:"test_case"`
    Submission      []Submission  `gorm:"foreignKey:ProblemID" json:"submission"`
}
