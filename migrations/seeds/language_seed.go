package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"mods/domain/entity"
	"os"

	"gorm.io/gorm"
)

func LanguageSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/language.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var languages []entity.Language
	if err := json.Unmarshal(jsonData, &languages); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Language{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Language{}); err != nil {
			return err
		}
	}

	for _, lang := range languages {
		var existingLang entity.Language
		err := db.Where(&entity.Language{Name: lang.Name}).First(&existingLang).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if db.Find(&existingLang, "name = ?", lang.Name).RowsAffected == 0 {
			if err := db.Create(&lang).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
