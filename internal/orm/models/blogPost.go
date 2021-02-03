package models

// BlogPost defines a blog post
type BlogPost struct {
	BaseModelSoftDelete
	Title       string `gorm:"not null;"`
	Image       *string
	Status      string `gorm:"not null;"` // active, draft
	TagID       string
	Tag         Tag
	Details     string `gorm:"not null;"`
	CreatedByID *string
	CreatedBy   *Admin
}
