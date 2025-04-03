package migrations

import (
	"mods/migrations/seeds"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.UserRoleSeeder(db); err != nil {
		return err
	}
	if err := seeds.UserClassRoleSeeder(db); err != nil {
		return err
	}
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}
	if err := seeds.LanguageSeeder(db); err != nil {
		return err
	}

	return nil
}