package entity

import "github.com/google/uuid"

type UserClass struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ClassID uuid.UUID `gorm:"type:uuid;not null" json:"class_id"`
	RoleID  uint `gorm:"type:uuid;not null" json:"role_id"`

	// Relationships
	User  User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Class Class         `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"class"`
	Role  UserClassRole `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"role"`
}
