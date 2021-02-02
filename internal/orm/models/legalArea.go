package models

import "github.com/gofrs/uuid"

// LegalArea defines the legal areas in the system
type LegalArea struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	Image       *string
	Description *string
	CreatedByID *uuid.UUID
	CreatedBy   *Admin
}
