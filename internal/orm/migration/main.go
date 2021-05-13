package migration

import (
	"fmt"

	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/getsentry/raven-go"

	"github.com/Bendomey/avc-server/internal/orm/migration/jobs"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func updateMigration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Admin{},
		&models.Country{},
		&models.User{},
		&models.Customer{},
		&models.Lawyer{},
		&models.Faq{},
		&models.Tag{},
		&models.LegalArea{},
		&models.BlogPost{},
		&models.NewsletterSubscribers{},
		&models.Package{},
		&models.Service{},
		&models.PackageService{},
		&models.Subscription{},
		&models.Servicing{},
		&models.ServicingField{},
		&models.Payment{},
	)
	return err
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB, seedDB bool) error {
	// Keep a list of migrations here
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		log.Info("[Migration.InitSchema] Initializing database schema")
		db.Exec("create extension \"uuid-ossp\";")
		if err := updateMigration(db); err != nil {
			raven.CaptureError(err, nil)
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		// Add more jobs, etc here
		return nil
	})
	m.Migrate()

	if err := updateMigration(db); err != nil {
		raven.CaptureError(err, nil)
		return err
	}

	if seedDB {
		m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
			jobs.SeedSuperAdmin,
		})
	}

	return m.Migrate()
}
