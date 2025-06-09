package entity

import "github.com/google/uuid"

type Submission struct {
	ID          		uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      		string      `gorm:"not null" json:"user_id"`
	ExamID      		string      `gorm:"not null" json:"exam_id"`
	ProblemID   		string      `gorm:"not null" json:"problem_id"`
	Code     			string      `json:"code" binding:"required"`
	LangID 				uint      `json:"lang_id" binding:"required"`
	SubmissionTime    	string      `json:"submission_time" binding:"required"`
	Judge0Token			string		`json:"judge0_token" binding:"required"`
	StatusId     	    uint      `json:"status_id" binding:"required"`
	Time 				string		`json:"time"`
	Memory				string		`json:"memory"`

	User      			User    	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Exam      			Exam    	`gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
	Problem      		Problem    	`gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"problem"`
	Language      		Language    `gorm:"foreignKey:LangID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"language"`
    Status        		SubmissionStatus `gorm:"foreignKey:StatusId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"role"`
	Timestamp
}