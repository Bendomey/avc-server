package models

import (
	"errors"
	"time"

	"github.com/Bendomey/goutilities/pkg/hashpassword"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Admin defines a administrators who are to manage the app
type Admin struct {
	BaseModelSoftDelete
	FullName        string  `gorm:"not null;"`
	Email           string  `gorm:"not null;unique"`
	Password        string  `gorm:"not null;"`
	Phone           *string `gorm:"unique"`
	EmailVerifiedAt *time.Time
	CreatedByID     *uuid.UUID
	CreatedBy       *Admin
}

// BeforeSave hook is called before the data is persisted to db
func (admin *Admin) BeforeSave(tx *gorm.DB) (err error) {
	//hashes password
	hashed, err := hashpassword.HashPassword(admin.Password)
	admin.Password = hashed
	if err != nil {
		err = errors.New("CannotHashAdminPassword")
	}
	return
}

// BeforeDelete hook is called before the data is delete so that we dont delete super admin
func (admin *Admin) BeforeDelete(tx *gorm.DB) (err error) {
	if admin.CreatedByID == nil {
		err = errors.New("CannotDeleteSuperAdmin")
	}
	return
}
