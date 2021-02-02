package models

import "github.com/gofrs/uuid"

// Country defines countries we are working with
type Country struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	Description *string
	Currency    *string
	Image       *string
	CreatedByID *uuid.UUID
	CreatedBy   *Admin
}
