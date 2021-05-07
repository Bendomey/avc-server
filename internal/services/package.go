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

// PackageService inteface holds the Package-databse transactions of this controller
type PackageService interface {
	CreatePackage(context context.Context, name string, createdBy string, amountPerMonth *int, amountPerYear *int) (*models.Package, error)
	UpdatePackage(context context.Context, packageID string, name *string, amountPerMonth *int, amountPerYear *int) (bool, error)
	DeletePackage(context context.Context, packageID string) (bool, error)
	ReadPackage(ctx context.Context, packageID string) (*models.Package, error)
	ReadPackages(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Package, error)
	ReadPackagesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error)
}

func PackageSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) PackageService {
	return &ORM{db, rdb, mg}
}

// CreatePackage adds a new package to the system
func (orm *ORM) CreatePackage(context context.Context, name string, createdBy string, amountPerMonth *int, amountPerYear *int) (*models.Package, error) {
	__Package := models.Package{
		Name:           name,
		CreatedByID:    createdBy,
		AmountPerMonth: amountPerMonth,
		AmountPerYear:  amountPerYear,
	}

	if err := orm.DB.DB.Select("Name", "CreatedByID", "AmountPerMonth", "AmountPerYear").Create(&__Package).Error; err != nil {
		return nil, err
	}

	return &__Package, nil
}

func (orm *ORM) UpdatePackage(context context.Context, packageID string, name *string, amountPerMonth *int, amountPerYear *int) (bool, error) {
	var __Package models.Package

	err := orm.DB.DB.First(&__Package, "id = ?", packageID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("PackageNotFound")
	}

	if name != nil {
		__Package.Name = *name
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

func (orm *ORM) ReadPackages(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Package, error) {
	var __Packages []*models.Package

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("packages.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("packages.Name = ?", name)
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

func (orm *ORM) ReadPackagesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error) {
	var __PackagesLength int64

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("packages.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("packages.Name = ?", name)
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
