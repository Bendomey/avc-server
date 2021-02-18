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

// LegalAreaService inteface holds the legalarea-databse transactions of this controller
type LegalAreaService interface {
	CreateLegalArea(ctx context.Context, name string, description *string, image *string, adminID string) (*models.LegalArea, error)
	UpdateLegalArea(ctx context.Context, legalAreaID string, name *string, description *string, image *string) (bool, error)
	DeleteLegalArea(ctx context.Context, legalAreaID string) (bool, error)
	ReadLegalArea(ctx context.Context, legalAreaID string) (*models.LegalArea, error)
	ReadLegalAreas(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) ([]*models.LegalArea, error)
	ReadLegalAreasLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) (*int64, error)
}

// NewLegalAreaSvc exposed the ORM to the LegalArea functions in the module
func NewLegalAreaSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) LegalAreaService {
	return &ORM{db, rdb, mg}
}

//CreateLegalArea creates a legal area
func (orm *ORM) CreateLegalArea(ctx context.Context, name string, description *string, image *string, adminID string) (*models.LegalArea, error) {
	_LegalArea := models.LegalArea{
		Name:        name,
		Description: description,
		Image:       image,
		CreatedByID: adminID,
	}

	err := orm.DB.DB.Select("Name", "Description", "Image", "CreatedByID").Create(&_LegalArea).Error
	if err != nil {
		return nil, err
	}
	// return success
	return &_LegalArea, nil
}

//UpdateLegalArea updates a legal are in question
func (orm *ORM) UpdateLegalArea(ctx context.Context, legalAreaID string, name *string, description *string, image *string) (bool, error) {
	var _LegalArea models.LegalArea

	err := orm.DB.DB.First(&_LegalArea, "id = ?", legalAreaID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("LegalAreaNotFound")
	}

	if name != nil {
		_LegalArea.Name = *name
	}

	if description != nil {
		_LegalArea.Description = description
	}

	if image != nil {
		_LegalArea.Image = image
	}

	orm.DB.DB.Save(&_LegalArea)

	// return success
	return true, nil
}

//DeleteLegalArea deletes the legal area in question
func (orm *ORM) DeleteLegalArea(ctx context.Context, legalAreaID string) (bool, error) {
	var _LegalArea models.LegalArea

	err := orm.DB.DB.First(&_LegalArea, "id = ?", legalAreaID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("LegalAreaNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_LegalArea).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadLegalArea calls the single in LegalArea
func (orm *ORM) ReadLegalArea(ctx context.Context, legalAreaID string) (*models.LegalArea, error) {
	var _LegalArea models.LegalArea

	err := orm.DB.DB.Joins("CreatedBy").First(&_LegalArea, "legal_areas.id = ?", legalAreaID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("LegalAreaNotFound")
	}

	// return success
	return &_LegalArea, nil
}

//ReadLegalAreas based on a query
func (orm *ORM) ReadLegalAreas(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) ([]*models.LegalArea, error) {
	var _LegalAreas []*models.LegalArea

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("legal_areas.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("legal_areas.Name = ?", name)
	}

	if description != nil {
		_Results = _Results.Where("legal_areas.Description = ?", description)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("legal_areas.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("legal_areas.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_LegalAreas)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _LegalAreas, nil
}

//ReadLegalAreasLength retirieved the count based on a query
func (orm *ORM) ReadLegalAreasLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string, description *string) (*int64, error) {
	var _LegalAreasLength int64

	_Results := orm.DB.DB.Model(&models.LegalArea{})
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("legal_areas.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("legal_areas.Name = ?", name)
	}

	if description != nil {
		_Results = _Results.Where("legal_areas.Description = ?", description)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("legal_areas.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("legal_areas.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Count(&_LegalAreasLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_LegalAreasLength, nil
}
