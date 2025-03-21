package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exam struct {
    ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Name        string        `json:"name"`
    ShortName   string        `json:"short_name"`
	IsPublished bool          `json:"is_published"`
    StartTime   time.Time     `json:"start_time"`           
    Duration    time.Duration `json:"duration"`             
	EndTime     time.Time     `json:"end_time"`           
    CreatedBy   string        `json:"created_by"`           
    
	Timestamp

    // Relationships
    UserExams []UserExam `gorm:"foreignKey:ExamID" json:"user_exams"`
    Problems  []Problem  `gorm:"foreignKey:ExamID" json:"problems"`
}


func (e *Exam) CalculateEndTime() {
    e.EndTime = e.StartTime.Add(e.Duration)
}

func (e *Exam) BeforeCreate(tx *gorm.DB) (err error) {
    e.CalculateEndTime()
    return
}

func (e *Exam) BeforeUpdate(tx *gorm.DB) (err error) {
    e.CalculateEndTime()
    return
}
