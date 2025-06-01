package entity

import "github.com/google/uuid"

type UserClass struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID  string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_class" json:"user_id"`
	ClassID string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_class" json:"class_id"`

	// Relationships
	User  User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Class Class         `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"class"`
}
