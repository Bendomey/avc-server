package models

type PackageServiceType string

const (
	Boolean PackageServiceType = "BOOLEAN"
	Number  PackageServiceType = "NUMBER"
)

type PackageService struct {
	BaseModelSoftDelete
	ServiceID     string `gorm:"not null;"`
	Service       Service
	PackageID     string `gorm:"not null;"`
	Package       Package
	Type          PackageServiceType
	Quantity      *int
	IsActive      *bool
	CreatedByID   *string
	CreatedBy     Admin
	RequestedByID *string
	RequestedBy   User
}
