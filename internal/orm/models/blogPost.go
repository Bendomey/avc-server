package models

import "github.com/gofrs/uuid"

// BlogPost defines a blog post
type BlogPost struct {
	BaseModelSoftDelete
	Title       string `gorm:"not null;"`
	Image       *string
	Status      string `gorm:"not null;"` // active, draft
	TagID       uuid.UUID
	Tag         Tag
	Details     string `gorm:"not null;"`
	CreatedByID *uuid.UUID
	CreatedBy   *Admin
}
