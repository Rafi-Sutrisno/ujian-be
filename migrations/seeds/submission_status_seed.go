package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"mods/domain/entity"
	"os"

	"gorm.io/gorm"
)

func SubmissionStatusSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/submission_status.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var statuses []entity.SubmissionStatus
	if err := json.Unmarshal(jsonData, &statuses); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.SubmissionStatus{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.SubmissionStatus{}); err != nil {
			return err
		}
	}

	for _, status := range statuses {
		var existingStatus entity.SubmissionStatus
		err := db.Where("name = ?", status.Name).First(&existingStatus).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&status).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
