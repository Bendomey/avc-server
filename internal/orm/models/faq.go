package models

import "github.com/gofrs/uuid"

// Faq defines structure of frequently asked questions
type Faq struct {
	BaseModel
	Question    string `gorm:"not null;"`
	Answer      string `gorm:"not null;"`
	CreatedByID *uuid.UUID
	CreatedBy   *Admin
}
