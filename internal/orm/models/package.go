package models

// Package defines packages a user can subscribe to
type Package struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	AmountPerMonth        uint64 `gorm:"not null;"`
	AmountPerYear        uint64 `gorm:"not null;"`
	CreatedByID string
	CreatedBy   Admin
}
