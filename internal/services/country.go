package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// CountryService inteface holds the country-databse transactions of this controller
type CountryService interface {
	CreateCountry(ctx context.Context, name string, description *string, currency *string, image *string, adminID string) (*models.Country, error)
	UpdateCountry(ctx context.Context, countryID string, name *string, description *string, currency *string, image *string) (bool, error)
	DeleteCountry(ctx context.Context, countryID string) (bool, error)
	ReadCountry(ctx context.Context, countryID string) (*models.Country, error)
	ReadCountries(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) ([]*models.Country, error)
	ReadCountriesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) (*int64, error)
}

// NewCountrySvc exposed the ORM to the country functions in the module
func NewCountrySvc(db *orm.ORM, rdb *redis.Client) CountryService {
	return &ORM{db, rdb}
}

//CreateCountry creates a country
func (orm *ORM) CreateCountry(ctx context.Context, name string, description *string, currency *string, image *string, adminID string) (*models.Country, error) {
	_Country := models.Country{
		Name:        name,
		Description: description,
		Currency:    currency,
		Image:       image,
		CreatedByID: adminID,
	}

	err := orm.DB.DB.Select("Name", "Description", "Currency", "Image", "CreatedByID").Create(&_Country).Error
	if err != nil {
		return nil, err
	}

	// return success
	return &_Country, nil
}

//UpdateCountry updates the country in question
func (orm *ORM) UpdateCountry(ctx context.Context, countryID string, name *string, description *string, currency *string, image *string) (bool, error) {
	var _Country models.Country

	err := orm.DB.DB.First(&_Country, "id = ?", countryID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("CountryNotFound")
	}

	if name != nil {
		_Country.Name = *name
	}

	if description != nil {
		_Country.Description = description
	}

	if currency != nil {
		_Country.Currency = currency
	}

	if image != nil {
		_Country.Image = image
	}

	orm.DB.DB.Save(&_Country)

	// return success
	return true, nil
}

//DeleteCountry deletes the country in question
func (orm *ORM) DeleteCountry(ctx context.Context, countryID string) (bool, error) {
	var _Country models.Country

	err := orm.DB.DB.First(&_Country, "id = ?", countryID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("CountryNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_Country).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadCountry calls the single in country
func (orm *ORM) ReadCountry(ctx context.Context, countryID string) (*models.Country, error) {
	var _Country models.Country

	err := orm.DB.DB.Joins("CreatedBy").First(&_Country, "countries.id = ?", countryID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("CountryNotFound")
	}

	// return success
	return &_Country, nil
}

//ReadCountries based on a query
func (orm *ORM) ReadCountries(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) ([]*models.Country, error) {
	var _Countries []*models.Country

	_Results := orm.DB.DB
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("countries.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("countries.Name = ?", name)
	}

	if description != nil {
		_Results = _Results.Where("countries.Description = ?", description)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("countries.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("countries.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Countries)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Countries, nil
}

//ReadCountriesLength retirieved the count based on a query
func (orm *ORM) ReadCountriesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) (*int64, error) {
	var _CountriesLength *int64

	_Results := orm.DB.DB
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("countries.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("countries.Name = ?", name)
	}

	if description != nil {
		_Results = _Results.Where("countries.Description = ?", description)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("countries.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("countries.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Count(_CountriesLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _CountriesLength, nil
}
