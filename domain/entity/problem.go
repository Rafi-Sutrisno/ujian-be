package entity

import (
	"github.com/google/uuid"
)

type Problem struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Title         string    `json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Constraints   string    `json:"constraints"`
	SampleInput   string    `json:"sample_input"`
	SampleOutput  string    `json:"sample_output"`
	CpuTimeLimit      float64 `json:"cpu_time_limit,omitempty"`
    MemoryLimit       int     `json:"memory_limit,omitempty"`
    Timestamp

    // Relationships
    TestCase        []TestCase `gorm:"foreignKey:ProblemID" json:"test_case"`
    Submission      []Submission  `gorm:"foreignKey:ProblemID" json:"submission"`
}
