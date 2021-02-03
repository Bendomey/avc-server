package models

import (
	"errors"
	"time"

	"github.com/Bendomey/goutilities/pkg/hashpassword"
	"gorm.io/gorm"
)

// Admin defines a administrators who are to manage the app
type Admin struct {
	BaseModelSoftDelete
	FullName        string  `gorm:"not null;"`
	Email           string  `gorm:"not null;unique"`
	Password        string  `gorm:"not null;"`
	Role            string  `gorm:"not null;"`
	Phone           *string `gorm:"unique"`
	PhoneVerifiedAt *time.Time
	CreatedByID     *string
	CreatedBy       *Admin
	SuspendedAt     *time.Time
	SuspendedReason *string
	SuspendByID     *string
	SuspendBy       *Admin
}

// BeforeCreate hook is called before the data is persisted to db
func (admin *Admin) BeforeCreate(tx *gorm.DB) (err error) {
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
