package models

import (
	"time"
)

// Customer defines customers in the system. They could be businesses or individuals
type Customer struct {
	BaseModelSoftDelete
	UserID string `gorm:"not null;"`
	User   User
	Type   string `gorm:"not null;"` // business, individual
	TIN    *string

	// For address
	DigitalAddress    *string
	AddressCountry    *string
	AddressCity       *string
	AddressStreetName *string
	AddressNumber     *string

	// Company
	CompanyName                  *string
	CompanyEntityType            *string
	CompanyEntityTypeOther       *string
	CompanyCountryOfRegistration *string
	CompanyDateOfRegistration    *time.Time
	CompanyRegistrationNumber    *string

	SuspendAt   *time.Time
	SuspendByID *string
	SuspendBy   *Admin
}
