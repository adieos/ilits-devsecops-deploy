package migrations

import (
	"github.com/adieos/ilits-devsecops-deploy/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := db.AutoMigrate(
		&entity.User{},
	); err != nil {
		return err
	}

	return nil
}
