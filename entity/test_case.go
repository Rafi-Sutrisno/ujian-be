package entity

import (
	"github.com/google/uuid"
)

type TestCase struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProblemID       string    `gorm:"type:uuid;not null" json:"problem_id"`
	InputData       string    `json:"input_data"`
	ExpectedOutput  string    `json:"expected_output"`
	CreatedBy       string    `json:"created_by"`

    Timestamp

    // Relationships
    Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"problem"`
}
