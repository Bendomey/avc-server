package models

import "github.com/gofrs/uuid"

// Tag defines categoru of a certain blog post
type Tag struct {
	BaseModelSoftDelete
	name        string `gorm:"not null;"`
	CreatedByID *uuid.UUID
	CreatedBy   *Admin
}
