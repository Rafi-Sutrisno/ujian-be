package migrations

import (
	"mods/domain/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Migrator().DropTable(
		// &entity.UserRole{},
		// &entity.Class{},
		// &entity.Exam{},
		// &entity.Problem{},  
		// &entity.TestCase{},            
		// &entity.User{},         
		// &entity.UserClass{}, 
		// &entity.Submission{}, 
		// &entity.SubmissionStatus{}, 
		// &entity.Language{},    
		// &entity.ExamLang{},
		// &entity.ExamSesssion{},
		&entity.ExamProblem{},
		// &entity.UserCodeDraft{},
	 )

	if err := db.AutoMigrate(
		// &entity.UserRole{},    
		// &entity.Class{},
		// &entity.Exam{},
		// &entity.Problem{},  
		// &entity.TestCase{},            
		// &entity.User{},         
		// &entity.UserClass{}, 
		// &entity.Submission{}, 
		// &entity.SubmissionStatus{}, 
		// &entity.Language{},    
		// &entity.ExamLang{}, 
		// &entity.ExamSesssion{},
		&entity.ExamProblem{},
		// &entity.UserCodeDraft{},
	); err != nil {
		return err
	}

	return nil
}
