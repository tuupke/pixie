package main

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/tuupke/pixie"
)

func migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "2022-10-05 setup tables",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&ExternalData{}, &Problem{}, &pixie.Setting{}, &Host{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("external_data", "problems", "settings", "hosts")
			},
		},
		{
			ID: "2022-10-05-insert-settings",
			Migrate: func(tx *gorm.DB) error {
				return tx.Create([]pixie.Setting{
					{Key: "contest", Value: []byte("nwerc18")},
					{Key: "domjudge", Value: []byte("https://www.domjudge.org/demoweb/api/v4/")},
				}).Error
			},
		},
	}
}
