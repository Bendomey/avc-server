package models

// Country defines countries we are working with
type Country struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;unique"`
	Description *string
	Currency    *string
	Image       *string
	CreatedByID string
	CreatedBy   Admin
}
