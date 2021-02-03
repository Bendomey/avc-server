package models

// Tag defines categoru of a certain blog post
type Tag struct {
	BaseModelSoftDelete
	name        string `gorm:"not null;"`
	CreatedByID *string
	CreatedBy   *Admin
}
