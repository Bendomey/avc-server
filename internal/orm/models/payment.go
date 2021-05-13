package models

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

// Payment defines payments a user creates
type Payment struct {
	BaseModelSoftDelete
	Code             string  `gorm:"not null;"`
	Amount           float64 `gorm:"not null;"`
	ServicingID      *string `gorm:"not null;"`
	Servicing        *Servicing
	SubscriptionID   *string `gorm:"not null;"`
	Subscription     *Subscription
	AuthorizationUrl string        `gorm:"not null;"`
	AccessCode       string        `gorm:"not null;"`
	Status           PaymentStatus `gorm:"default:PENDING"`
	CreatedByID      string
	CreatedBy        User
}
