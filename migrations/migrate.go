package migrations

import (
	"mods/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Migrator().DropTable(
		&entity.UserRole{},
		&entity.Class{},
		&entity.Exam{},
		&entity.Problem{},  
		&entity.TestCase{},            
		&entity.User{},         
		&entity.UserClass{}, 
		&entity.Submission{}, 
		&entity.Language{},    
		&entity.ExamLang{}, )

	if err := db.AutoMigrate(
		&entity.UserRole{},    
		&entity.Class{},
		&entity.Exam{},
		&entity.Problem{},  
		&entity.TestCase{},            
		&entity.User{},         
		&entity.UserClass{}, 
		&entity.Submission{}, 
		&entity.Language{},    
		&entity.ExamLang{}, 
	); err != nil {
		return err
	}

	return nil
}
