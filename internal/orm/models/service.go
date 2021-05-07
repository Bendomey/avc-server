package models

type ServiceType string

const (
	Subscribe   ServiceType = "SUBSCRIBE"
	Unsubscribe ServiceType = "UNSUBSCRIBE"
	Both        ServiceType = "BOTH"
)

// Service defines services a user can subscribe to
type Service struct {
	BaseModelSoftDelete
	Name        string `gorm:"not null;"`
	Price       *string
	Description *string
	Type        ServiceType `gorm:"default:SUBSCRIBE"` // subscribe/unsubscribe/both
	CreatedByID string
	CreatedBy   Admin
}
