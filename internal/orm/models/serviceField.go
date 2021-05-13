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
	Type                      *string
	Owners                    *string
	Address                   *string
	ClassificationOfTrademark *string
	Uploads                   *string
}

type DocumentType string

const (
	DocumentNew      DocumentType = "NEW"
	DocumentExisting DocumentType = "EXISTING"
)

type Document struct {
	Type              *DocumentType
	NatureOfDoc       *string
	Deadline          *time.Time
	ExistingDocuments *string
	NewDocuments      *string
}

type ServicingField struct {
	Business  Business  `gorm:"embedded; embeddedPrefix:business_"`
	Trademark Trademark `gorm:"embedded; embeddedPrefix:trademark_"`
	Document  Document  `gorm:"embedded; embeddedPrefix:document_"`
}
