package models

type ServicingStatus string

const (
	ServicingPending ServicingStatus = "PENDING"
	ServicingPaid    ServicingStatus = "PAID"
	ServicingActive  ServicingStatus = "ACTIVE"
	ServicingDone    ServicingStatus = "DONE"
)

// Servicing defines each single services user has made
type Servicing struct {
	BaseModelSoftDelete
	ServiceID      string `gorm:"not null;"`
	Service        Service
	PaymentID      *string
	Payment        Payment
	Status         ServicingStatus `gorm:"default:PENDING"`
	SubscriptionID *string
	Subscription   Subscription
	LawyerID       *string
	Lawyer         User
	CreatedByID    string
	CreatedBy      User
}
