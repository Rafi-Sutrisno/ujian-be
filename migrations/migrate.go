package migrations

import (
	"mods/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Migrator().DropTable(&entity.UserRole{}, &entity.UserClassRole{})

	if err := db.AutoMigrate(
		&entity.UserRole{},    
		&entity.UserClassRole{}, 
		&entity.Class{},
		&entity.Exam{},            
		&entity.User{},         
		&entity.UserClass{},    
	); err != nil {
		return err
	}

	return nil
}
