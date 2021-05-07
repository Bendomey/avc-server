package models

// Package defines packages a user can subscribe to
type Package struct {
	BaseModelSoftDelete
	Name           string `gorm:"not null;"`
	AmountPerMonth *int
	AmountPerYear  *int
	CreatedByID    string
	CreatedBy      Admin
}
