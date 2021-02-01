package jobs

import (
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/go-gormigrate/gormigrate/v2"
)

var (
	fullname = utils.MustGet("ADMIN_NAME")
	email    = utils.MustGet("ADMIN_EMAIL")
	password = utils.MustGet("ADMIN_PASSWORD")
	phone    = utils.MustGet("ADMIN_PHONE")
	role     = utils.MustGet("ADMIN_ROLE")
	// superAdmin *models.Admin = &models.Admin{
	// 	FullName: fullname,
	// 	Email:    email,
	// 	Password: password,
	// 	Role:     role,
	// 	Phone:    &phone,
	// }
)

// SeedSuperAdmin inserts the super user
var SeedSuperAdmin *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_SUPER_ADMIN",
	// Migrate: func(db *gorm.DB) error {
	// 	// return db.Create(&superAdmin).Error
	// },
	// Rollback: func(db *gorm.DB) error {
	// 	return db.Delete(&superAdmin).Error
	// },
}
