package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"mods/entity"
	"os"

	"gorm.io/gorm"
)

func UserClassRoleSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/user_class_roles.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var roles []entity.UserClassRole
	if err := json.Unmarshal(jsonData, &roles); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.UserClassRole{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.UserClassRole{}); err != nil {
			return err
		}
	}

	for _, role := range roles {
		var existingRole entity.UserClassRole
		err := db.Where(&entity.UserClassRole{Name: role.Name}).First(&existingRole).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if db.Find(&existingRole, "name = ?", role.Name).RowsAffected == 0 {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
