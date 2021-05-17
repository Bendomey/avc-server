package models

import "time"

type Business struct {
	CountryID      *string
	Country        Country
	EntityType     *string
	Name           *string
	Owners         *string
	Directors      *string
	Address        *string
	NumberOfShares *string
	InitialCapital *float32
	Industry       *string
}

type Trademark struct {
	CountryID                 *string
	Country                   Country
	OwnershipType             *string
	Owners                    *string
	Address                   *string
	ClassificationOfTrademark *string
	Uploads                   *string
}

type Document struct {
	Type              *string
	NatureOfDoc       *string
	Deadline          *time.Time
	ExistingDocuments *string
	NewDocuments      *string
}

type ServicingField struct {
	BaseModelSoftDelete
	Business  Business  `gorm:"embedded; embeddedPrefix:business_"`
	Trademark Trademark `gorm:"embedded; embeddedPrefix:trademark_"`
	Document  Document  `gorm:"embedded; embeddedPrefix:document_"`
}
