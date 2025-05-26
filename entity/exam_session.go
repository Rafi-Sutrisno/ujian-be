package entity

import "github.com/google/uuid"

type ExamSesssion struct {
	ID          		uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      		string      `gorm:"not null" json:"user_id"`
	ExamID      		string      `gorm:"not null" json:"exam_id"`
	SessionID			string		`gorm:"not null" json:"session_id"`
	IPAddress  			string		`gorm:"not null" json:"ip_address"`
	UserAgent  			string		`gorm:"not null" json:"user_agent"`
	Device     			string		`gorm:"not null" json:"device"`
	Status     			uint		`gorm:"not null" json:"status"`
	TotalCorrect		uint		`gorm:"not null" json:"total_correct"`

	User      			User    	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Exam      			Exam    	`gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"exam"`
	Timestamp
}