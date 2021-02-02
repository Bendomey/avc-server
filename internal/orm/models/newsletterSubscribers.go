package models

// NewsletterSubscribers defines those who have subscribed for newsletters
type NewsletterSubscribers struct {
	BaseModel
	Email string `gorm:"not null;unique"`
	Type  string `gorm:"not null;"` // anon, user
}
