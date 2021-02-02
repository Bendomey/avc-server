package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Lawyer defines the registred laywers in the system
type Lawyer struct {
	BaseModelSoftDelete
	UserID uuid.UUID `gorm:"not null;"`
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

	SuspendAt   *time.Time
	SuspendByID *uuid.UUID
	SuspendBy   *Admin

	ApprovedAt   *time.Time
	ApprovedByID *uuid.UUID
	ApprovedBy   *Admin
}
