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

// Service inteface holds the Service-databse transactions of this controller
type PackageServiceService interface {
	CreatePackageService(context context.Context, serviceID string, packageID string, packageServiceType models.PackageServiceType, quantity *int, isActive *bool, createdBy string) (*models.PackageService, error)
	UpdatePackageService(context context.Context, packageServiceID string, serviceID *string, packageID *string, packageServiceType *string, quantity *int, isActive *bool) (bool, error)
	DeletePackageService(context context.Context, packageServiceID string) (bool, error)
	ReadPackageService(ctx context.Context, packageServiceID string) (*models.PackageService, error)
	ReadPackageServices(ctx context.Context, filterQuery *utils.FilterQuery, serviceID *string, packageID *string) ([]*models.PackageService, error)
	ReadPackageServicesLength(ctx context.Context, filterQuery *utils.FilterQuery, serviceID *string, packageID *string) (*int64, error)
}

func PackageServiceSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) PackageServiceService {
	return &ORM{db, rdb, mg}
}

// CreatePackageService adds a new package service to the system
func (orm *ORM) CreatePackageService(context context.Context, serviceID string, packageID string, packageServiceType models.PackageServiceType, quantity *int, isActive *bool, createdBy string) (*models.PackageService, error) {
	__PackageService := models.PackageService{
		ServiceID:   serviceID,
		PackageID:   packageID,
		Type:        packageServiceType,
		Quantity:    quantity,
		IsActive:    isActive,
		CreatedByID: createdBy,
	}

	if err := orm.DB.DB.Select("ServiceID", "PackageID", "Type", "Quantity", "IsActive", "CreatedByID").Create(&__PackageService).Error; err != nil {
		return nil, err
	}

	return &__PackageService, nil
}

//UpdatePackageService updates a package service
func (orm *ORM) UpdatePackageService(context context.Context, packageServiceID string, serviceID *string, packageID *string, packageServiceType *string, quantity *int, isActive *bool) (bool, error) {
	var __PackageService models.PackageService

	err := orm.DB.DB.First(&__PackageService, "id = ?", packageServiceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageServiceNotFound")
	}

	if serviceID != nil {
		__PackageService.ServiceID = *serviceID
	}

	if packageID != nil {
		__PackageService.PackageID = *packageID
	}

	if quantity != nil {
		__PackageService.Quantity = quantity
	}

	if isActive != nil {
		__PackageService.IsActive = isActive
	}

	if packageServiceType != nil {
		__PackageService.Type = models.PackageServiceType(*packageServiceType)
	}

	orm.DB.DB.Save(__PackageService)

	//return success
	return true, nil
}

//DeletePackageService deletes a package service
func (orm *ORM) DeletePackageService(context context.Context, packageServiceID string) (bool, error) {
	var __ServicePackage models.PackageService

	err := orm.DB.DB.First(&__ServicePackage, "id = ?", packageServiceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageServiceNotFound")
	}

	//remove
	if delErr := orm.DB.DB.Delete(&__ServicePackage).Error; delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadPackageService reads a package service
func (orm *ORM) ReadPackageService(ctx context.Context, packageServiceID string) (*models.PackageService, error) {
	var __PackageService models.PackageService

	err := orm.DB.DB.Preload("Service").Preload("Package").First(&__PackageService, "id = ?", packageServiceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("PackageServiceNotFound")
	}

	return &__PackageService, nil
}

//ReadPackageServices reads all package services depending on a filter
func (orm *ORM) ReadPackageServices(ctx context.Context, filterQuery *utils.FilterQuery, serviceID *string, packageID *string) ([]*models.PackageService, error) {
	var __PackageServices []*models.PackageService

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("package_services.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if serviceID != nil {
		_Results = _Results.Where("package_services.ServiceID = ?", serviceID)
	}

	if packageID != nil {
		_Results = _Results.Where("package_services.PackageID = ?", packageID)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("package_services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("package_services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").Joins("Service").Joins("Package").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&__PackageServices)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return __PackageServices, nil
}

func (orm *ORM) ReadPackageServicesLength(ctx context.Context, filterQuery *utils.FilterQuery, serviceID *string, packageID *string) (*int64, error) {
	var __PackageServicesLength int64

	_Results := orm.DB.DB.Model(&models.Service{})

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("package_services.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if serviceID != nil {
		_Results = _Results.Where("package_services.ServiceID = ?", serviceID)
	}

	if packageID != nil {
		_Results = _Results.Where("package_services.PackageID = ?", packageID)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("package_services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("package_services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Count(&__PackageServicesLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &__PackageServicesLength, nil
}
