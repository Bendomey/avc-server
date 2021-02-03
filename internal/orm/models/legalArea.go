package models

// LegalArea defines the legal areas in the system
type LegalArea struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	Image       *string
	Description *string
	CreatedByID *string
	CreatedBy   *Admin
}
