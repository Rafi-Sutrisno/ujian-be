package entity

import (
	"mods/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	RoleID   uint      `gorm:"not null" json:"role_id"`
	Noid     string    `json:"noid" binding:"required"`

	UserClass   []UserClass   `gorm:"foreignKey:UserID" json:"user_class"`
	Submission  []Submission  `gorm:"foreignKey:UserID" json:"submission"`
	Role        UserRole      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"role"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error
	// u.ID = uuid.New()
	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
