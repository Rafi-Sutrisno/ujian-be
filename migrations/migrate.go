package migrations

import (
	"mods/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{}, &entity.Exam{}, &entity.UserExam{},
	); err != nil {
		return err
	}

	return nil
}