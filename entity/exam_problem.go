package entity

import "github.com/google/uuid"

type ExamProblem struct {
	ID      	uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ExamID  	string    `gorm:"type:uuid;not null;uniqueIndex:idx_exam_problem" json:"exam_id"`
	ProblemID 	string    `gorm:"type:uuid;not null;uniqueIndex:idx_exam_problem" json:"problem_id"`

	// Relationships
	Exam  		Exam      `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
	Problem		Problem     `gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"problem"`
}
