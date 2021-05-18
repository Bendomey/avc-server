package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CustomPackageService struct {
	ServiceId string
	Quantity  *int
	IsActive  *bool
}

// PackageService inteface holds the Package-databse transactions of this controller
type PackageService interface {
	CreatePackage(context context.Context, name string, createdBy string, amountPerMonth *int, description *string, amountPerYear *int, customPackages []CustomPackageService) (*models.Package, error)
	CreateCustomPackage(context context.Context, createdBy string, packageId string, customPackages []CustomPackageService, name *string, description *string) (*models.Package, error)
	ApprovePackage(context context.Context, packageID string, approvedBy string, amountPerMonth int, amountPerYear int) (bool, error)
	UpdatePackage(context context.Context, packageID string, name *string, description *string, amountPerMonth *int, amountPerYear *int) (bool, error)
	DeletePackage(context context.Context, packageID string) (bool, error)
	ReadPackage(ctx context.Context, packageID string) (*models.Package, error)
	ReadPackages(ctx context.Context, filterQuery *utils.FilterQuery, name *string, packagesType *string) ([]*models.Package, error)
	ReadPackagesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, packagesType *string) (*int64, error)
}

func PackageSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) PackageService {
	return &ORM{db, rdb, mg}
}

// CreatePackage adds a new package to the system
func (orm *ORM) CreatePackage(context context.Context, name string, createdBy string, amountPerMonth *int, description *string, amountPerYear *int, customPackages []CustomPackageService) (*models.Package, error) {
	__Package := models.Package{
		Name:           name,
		CreatedByID:    &createdBy,
		AmountPerMonth: amountPerMonth,
		AmountPerYear:  amountPerYear,
		Description:    description,
		Status:         "APPROVED",
	}

	if err := orm.DB.DB.Select("Name", "CreatedByID", "AmountPerMonth", "AmountPerYear", "Description", "Status").Create(&__Package).Error; err != nil {
		return nil, err
	}

	//create individual package services
	for _, packageService := range customPackages {
		var typeOfPackageService models.PackageServiceType
		if packageService.IsActive != nil {
			typeOfPackageService = "BOOLEAN"
		} else {
			typeOfPackageService = "NUMBER"
		}
		__newPackageService := models.PackageService{
			ServiceID: packageService.ServiceId,
			PackageID: __Package.ID.String(),
			Type:      typeOfPackageService,
			Quantity:  packageService.Quantity,
			IsActive:  packageService.IsActive,
		}
		if err := orm.DB.DB.Select("ServiceID", "PackageID", "Type", "Quantity", "IsActive").Create(&__newPackageService).Error; err != nil {
			return nil, err
		}
	}

	return &__Package, nil
}

// CreateCustomPackage adds a new package to the system
func (orm *ORM) CreateCustomPackage(context context.Context, createdBy string, packageId string, customPackages []CustomPackageService, name *string, description *string) (*models.Package, error) {
	var nameForCustomPackage string

	if name != nil {
		nameForCustomPackage = *name
	} else {
		// fetch custom packages created by user and count it
		var __PackagesLength int64 // hold length

		__Results := orm.DB.DB.Model(&models.Package{}).Where("packages.requested_by_id = ?", createdBy).Count(&__PackagesLength)
		if __Results.Error != nil {
			return nil, __Results.Error
		}
		nameForCustomPackage = fmt.Sprintf("Custom Package %d", __PackagesLength+1)
	}

	__Package := models.Package{
		Name:          nameForCustomPackage,
		RequestedByID: &createdBy,
		Status:        "PENDING",
		Description:   description,
	}

	if err := orm.DB.DB.Select("Name", "RequestedByID", "Status").Create(&__Package).Error; err != nil {
		return nil, err
	}

	//after creating package
	//create individual package services
	for _, packageService := range customPackages {
		var typeOfPackageService models.PackageServiceType
		if packageService.IsActive != nil {
			typeOfPackageService = "BOOLEAN"
		} else {
			typeOfPackageService = "NUMBER"
		}
		__newPackageService := models.PackageService{
			ServiceID:     packageService.ServiceId,
			PackageID:     __Package.ID.String(),
			Type:          typeOfPackageService,
			Quantity:      packageService.Quantity,
			IsActive:      packageService.IsActive,
			RequestedByID: &createdBy,
		}
		if err := orm.DB.DB.Select("ServiceID", "PackageID", "Type", "Quantity", "IsActive", "RequestedByID").Create(&__newPackageService).Error; err != nil {
			return nil, err
		}
	}

	return &__Package, nil
}

//Approve custom package
func (orm *ORM) ApprovePackage(context context.Context, packageID string, approvedBy string, amountPerMonth int, amountPerYear int) (bool, error) {
	var __Package models.Package

	err := orm.DB.DB.First(&__Package, "id = ?", packageID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageNotFound")
	}

	__Package.CreatedByID = &approvedBy
	__Package.AmountPerMonth = &amountPerMonth
	__Package.AmountPerYear = &amountPerYear
	__Package.Status = "APPROVED"
	orm.DB.DB.Save(__Package)

	//find user
	var __User models.User

	userErr := orm.DB.DB.First(&__User, "id = ?", __Package.RequestedByID).Error
	if errors.Is(userErr, gorm.ErrRecordNotFound) {
		return false, errors.New("UsereNotFound")
	}

	subject := "Package Approved - African Venture Counsel"
	name := "User"
	if __User.LastName != nil {
		name = *__User.LastName
	}
	body := fmt.Sprintf("Dear %s, your requested package has been approved. Visit our platform and then subscribe to this package", name)
	go orm.mg.SendTransactionalMail(context, subject, body, __User.Email)

	return true, nil
}

func (orm *ORM) UpdatePackage(context context.Context, packageID string, name *string, description *string, amountPerMonth *int, amountPerYear *int) (bool, error) {
	var __Package models.Package

	err := orm.DB.DB.First(&__Package, "id = ?", packageID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageNotFound")
	}

	if name != nil {
		__Package.Name = *name
	}
	if description != nil {
		__Package.Description = description
	}

	if amountPerMonth != nil {
		__Package.AmountPerMonth = amountPerMonth
	}

	if amountPerYear != nil {
		__Package.AmountPerYear = amountPerYear
	}

	orm.DB.DB.Save(__Package)

	//return success
	return true, nil
}

func (orm *ORM) DeletePackage(context context.Context, packageID string) (bool, error) {
	var __Package models.Package

	err := orm.DB.DB.First(&__Package, "id = ?", packageID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageNotFound")
	}

	//remove
	if delErr := orm.DB.DB.Delete(&__Package).Error; delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

func (orm *ORM) ReadPackage(ctx context.Context, packageID string) (*models.Package, error) {
	var __Package models.Package

	err := orm.DB.DB.First(&__Package, "id = ?", packageID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("PackageNotFound")
	}

	return &__Package, nil
}

func (orm *ORM) ReadPackages(ctx context.Context, filterQuery *utils.FilterQuery, name *string, packagesType *string) ([]*models.Package, error) {
	var __Packages []*models.Package

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("packages.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("packages.Name = ?", name)
	}

	if packagesType != nil {
		mainPackages := "MAIN"
		if packagesType == &mainPackages {
			_Results = _Results.Where("packages.RequestedByID = ?", nil)
		} else {
			_Results = _Results.Where("packages.CreatedByID = ? AND packages.Status = ?", nil, "PENDING")
		}
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("packages.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("packages.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&__Packages)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return __Packages, nil
}

func (orm *ORM) ReadPackagesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, packagesType *string) (*int64, error) {
	var __PackagesLength int64

	_Results := orm.DB.DB.Model(&models.Package{})

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("packages.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("packages.Name = ?", name)
	}

	if packagesType != nil {
		mainPackages := "MAIN"
		if packagesType == &mainPackages {
			_Results = _Results.Where("packages.RequestedByID = ?", nil)
		} else {
			_Results = _Results.Where("packages.CreatedByID = ? AND packages.Status = ?", nil, "PENDING")
		}
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("packages.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("packages.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Count(&__PackagesLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &__PackagesLength, nil
}
