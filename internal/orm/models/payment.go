package models

import "github.com/gofrs/uuid"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

// Payment defines payments a user creates
type Payment struct {
	BaseModelSoftDelete
	Code             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Amount           float64   `gorm:"not null;"`
	ServicingID      *string   `gorm:"not null;"`
	Servicing        *Servicing
	SubscriptionID   *string `gorm:"not null;"`
	Subscription     *Subscription
	AuthorizationUrl string        `gorm:"not null;"`
	AccessCode       string        `gorm:"not null;"`
	Status           PaymentStatus `gorm:"default:PENDING"`
	CreatedByID      string
	CreatedBy        User
}
