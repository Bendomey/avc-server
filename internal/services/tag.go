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

// TagService inteface holds the Tag-databse transactions of this controller
type TagService interface {
	CreateTag(context context.Context, name string, createdBy string) (*models.Tag, error)
	UpdateTag(context context.Context, tagID string, name *string) (bool, error)
	DeleteTag(context context.Context, tagID string) (bool, error)
	ReadTag(ctx context.Context, tagID string) (*models.Tag, error)
}

// TagSvc exposed the ORM to the Tag functions in the module
func TagSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) TagService {
	return &ORM{db, rdb, mg}
}

//CreateTag allow admins to create Tags
func (orm *ORM) CreateTag(context context.Context, name string, createdBy string) (*models.Tag, error) {
	_Tag := models.Tag{
		Name:        name,
		CreatedByID: createdBy,
	}

	err := orm.DB.DB.Select("Name", "CreatedByID").Create(&_Tag).Error
	if err != nil {
		return nil, err
	}

	return &_Tag, nil
}

//UpdateTag allow admins to update Tags
func (orm *ORM) UpdateTag(context context.Context, tagID string, name *string) (bool, error) {
	var _Tag models.Tag

	err := orm.DB.DB.First(&_Tag, "id = ?", tagID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("TagNotFound")
	}

	if name != nil {
		_Tag.Name = *name
	}

	orm.DB.DB.Save(&_Tag)

	// return success
	return true, nil
}

//DeleteTag allow admins to delete Tags
func (orm *ORM) DeleteTag(context context.Context, tagID string) (bool, error) {
	var _Tag models.Tag

	err := orm.DB.DB.First(&_Tag, "id = ?", tagID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("TagNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_Tag).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadTag calls the single in tag
func (orm *ORM) ReadTag(ctx context.Context, tagID string) (*models.Tag, error) {
	var _Tag models.Tag

	err := orm.DB.DB.Joins("CreatedBy").First(&_Tag, "tags.id = ?", tagID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("TagNotFound")
	}

	// return success
	return &_Tag, nil
}

//ReadTags based on a query
func (orm *ORM) ReadTags(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Tag, error) {
	var _Tag []*models.Tag

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("tags.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("tags.Name = ?", name)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("tags.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("tags.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Tag)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Tag, nil
}

//ReadTagsLength retirieved the count based on a query
func (orm *ORM) ReadTagsLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error) {
	var _TagsLength int64

	_Results := orm.DB.DB.Model(&models.LegalArea{})
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("tags.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("tags.Name = ?", name)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("tags.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("tags.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Count(&_TagsLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_TagsLength, nil
}
