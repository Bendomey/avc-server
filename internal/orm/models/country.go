package models

// Country defines countries we are working with
type Country struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	Description *string
	Currency    *string
	Image       *string
	// CreatedByID *string
	// CreatedBy   *Admin
}
