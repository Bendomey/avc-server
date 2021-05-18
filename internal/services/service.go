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
type ServiceService interface {
	CreateService(context context.Context, name string, price *float64, description *string, serviceType models.ServiceType, variant string, createdBy string) (*models.Service, error)
	UpdateService(context context.Context, serviceID string, name *string, price *float64, description *string, serviceType *string, variant *string) (bool, error)
	DeleteService(context context.Context, serviceeID string) (bool, error)
	ReadService(ctx context.Context, serviceID string) (*models.Service, error)
	ReadServices(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Service, error)
	ReadServicesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error)
}

func ServiceSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) ServiceService {
	return &ORM{db, rdb, mg}
}

// CreateService adds a new service to the system
func (orm *ORM) CreateService(context context.Context, name string, price *float64, description *string, serviceType models.ServiceType, variant string, createdBy string) (*models.Service, error) {
	__Service := models.Service{
		Name:        name,
		CreatedByID: createdBy,
		Price:       price,
		Description: description,
		Type:        serviceType,
		Variant:     variant,
	}

	if err := orm.DB.DB.Select("Name", "CreatedByID", "Price", "Description", "Type", "Variant").Create(&__Service).Error; err != nil {
		return nil, err
	}

	return &__Service, nil
}

//UpdateService updates a service
func (orm *ORM) UpdateService(context context.Context, serviceID string, name *string, price *float64, description *string, serviceType *string, variant *string) (bool, error) {
	var __Service models.Service

	err := orm.DB.DB.First(&__Service, "id = ?", serviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("ServiceNotFound")
	}

	if name != nil {
		__Service.Name = *name
	}

	if price != nil {
		__Service.Price = price
	}

	if description != nil {
		__Service.Description = description
	}

	if serviceType != nil {
		__Service.Type = models.ServiceType(*serviceType)
	}

	if variant != nil {
		__Service.Variant = *variant
	}

	orm.DB.DB.Save(__Service)

	//return success
	return true, nil
}
func (orm *ORM) DeleteService(context context.Context, serviceID string) (bool, error) {
	var __Service models.Service

	err := orm.DB.DB.First(&__Service, "id = ?", serviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("ServiceNotFound")
	}

	//remove
	if delErr := orm.DB.DB.Delete(&__Service).Error; delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

func (orm *ORM) ReadService(ctx context.Context, serviceID string) (*models.Service, error) {
	var __Service models.Service

	err := orm.DB.DB.First(&__Service, "id = ?", serviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("ServiceNotFound")
	}

	return &__Service, nil
}

func (orm *ORM) ReadServices(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Service, error) {
	var __Services []*models.Service

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("services.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("services.Name = ?", name)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&__Services)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return __Services, nil
}

func (orm *ORM) ReadServicesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error) {
	var __ServicesLength int64

	_Results := orm.DB.DB.Model(&models.Service{})

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("services.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if name != nil {
		_Results = _Results.Where("services.Name = ?", name)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("services.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Count(&__ServicesLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &__ServicesLength, nil
}
