package models

import (
	"errors"
	"time"

	"github.com/Bendomey/goutilities/pkg/hashpassword"
	"gorm.io/gorm"
)

// User defines the two typa users in the system (Lawyers and normal clients)
type User struct {
	BaseModelSoftDelete
	Type            string `gorm:"not null;"` // customer, lawyer
	LastName        *string
	FirstName       *string
	OtherNames      *string
	Email           string  `gorm:"not null;unique"`
	Password        string  `gorm:"not null;"`
	Phone           *string `gorm:"unique"`
	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time
}

// BeforeCreate hook is called before the data is persisted to db
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	//hashes password
	hashed, err := hashpassword.HashPassword(user.Password)
	user.Password = hashed
	if err != nil {
		err = errors.New("CannotHashUserPassword")
	}
	return
}
