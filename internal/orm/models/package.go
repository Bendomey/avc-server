package models

// Package defines packages a user can subscribe to
type Package struct {
	BaseModelSoftDelete
	Name           string `gorm:"not null;"`
	Description    *string
	AmountPerMonth *int
	AmountPerYear  *int
	Status         string // pending, approved
	CreatedByID    *string
	CreatedBy      Admin
	RequestedByID  *string //when they are creating a custom package
	RequestedBy    User
}
