package models

import (
	"time"
)

// Lawyer defines the registred laywers in the system
type Lawyer struct {
	BaseModelSoftDelete
	UserID string `gorm:"not null;"`
	User   User

	// For address
	DigitalAddress    *string
	AddressCountry    *string
	AddressCity       *string
	AddressStreetName *string
	AddressNumber     *string

	// lawyer stuff
	FirstYearOfBarAdmission *string
	LicenseNumber           *string
	TIN                     *string

	// Uploads
	NationalIDFront   *string
	NationalIDBack    *string
	BARMembershipCard *string
	LawCertificate    *string
	CV                *string
	CoverLetter       *string

	ApprovedAt   *time.Time
	ApprovedByID *string
	ApprovedBy   *Admin
}
