package models

// Faq defines structure of frequently asked questions
type Faq struct {
	BaseModel
	Question    string `gorm:"not null;"`
	Answer      string `gorm:"not null;"`
	CreatedByID *string
	CreatedBy   *Admin
}
