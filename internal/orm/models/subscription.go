package models

import "time"

type SubscriptionStatus string

const (
	SubscriptionPending SubscriptionStatus = "PENDING"
	SubscriptionActive  SubscriptionStatus = "ACTIVE"
	SubscriptionExpired SubscriptionStatus = "EXPIRED"
)

// Subscription defines what the user has subscribed
type Subscription struct {
	BaseModelSoftDelete
	PackageID   string `gorm:"not null;"`
	Package     Package
	PaymentID   *string
	Payment     Payment
	Status      SubscriptionStatus `gorm:"default:PENDING"`
	SubscribeAt time.Time
	ExpiresAt   time.Time
	CreatedByID string
	CreatedBy   User
}
