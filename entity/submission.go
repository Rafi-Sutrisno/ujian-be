package entity

import "github.com/google/uuid"

type Submission struct {
	ID          		uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      		string      `gorm:"not null" json:"user_id"`
	ExamID      		string      `gorm:"not null" json:"exam_id"`
	ProblemID   		string      `gorm:"not null" json:"problem_id"`
	Code     			string      `json:"code" binding:"required"`
	Language 			string      `json:"language" binding:"required"`
	SubmissionTime    	string      `json:"submission_time" binding:"required"`
	
	Status     			string      `json:"status" binding:"required"`

	User      			User    	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Exam      			Exam    	`gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
	Problem      		Problem    	`gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"problem"`

	Timestamp
}