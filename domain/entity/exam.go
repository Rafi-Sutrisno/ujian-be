package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exam struct {
    ID                      uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    ClassID                 string        `gorm:"type:uuid;not null" json:"class_id"`
    Name                    string        `json:"name"`
    ShortName               string        `json:"short_name"`
	IsPublished             bool          `json:"is_published"`
    StartTime               time.Time     `json:"start_time"`           
    Duration                time.Duration `json:"duration"`             
	EndTime                 time.Time     `json:"end_time"`
    IsSEBRestricted         bool          `gorm:"default:false" json:"is_seb_restricted"`
    SEBBrowserKey           string        `gorm:"type:text" json:"seb_browser_key"`   
    SEBConfigKey            string        `gorm:"type:text" json:"seb_config_key"`                      
    SEBQuitURL              string        `gorm:"type:text" json:"seb_quit_url"` 
	Timestamp

    Class       Class         `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"class"`
    ExamLang    []ExamLang    `gorm:"foreignKey:ExamID" json:"exam_lang"`
    Submission  []Submission  `gorm:"foreignKey:ExamID" json:"submission"`
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
