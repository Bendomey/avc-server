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

// FAQService inteface holds the FAQ-databse transactions of this controller
type FAQService interface {
	CreateFAQ(context context.Context, question string, answer string, createdBy string) (*models.Faq, error)
	UpdateFAQ(context context.Context, faqID string, question *string, answer *string) (bool, error)
	DeleteFAQ(context context.Context, faqID string) (bool, error)
	ReadFAQ(ctx context.Context, faqID string) (*models.Faq, error)
	ReadFAQs(ctx context.Context, filterQuery *utils.FilterQuery, question *string, answer *string) ([]*models.Faq, error)
	ReadFAQsLength(ctx context.Context, filterQuery *utils.FilterQuery, question *string, answer *string) (*int64, error)
}

// FAQSvc exposed the ORM to the FAQ functions in the module
func FAQSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) FAQService {
	return &ORM{db, rdb, mg}
}

//CreateFAQ allow admins to create Faqs
func (orm *ORM) CreateFAQ(context context.Context, question string, answer string, createdBy string) (*models.Faq, error) {
	_Faq := models.Faq{
		Question:    question,
		Answer:      answer,
		CreatedByID: createdBy,
	}

	err := orm.DB.DB.Select("Question", "Answer", "CreatedByID").Create(&_Faq).Error
	if err != nil {
		return nil, err
	}

	return &_Faq, nil
}

//UpdateFAQ allow admins to update Tags
func (orm *ORM) UpdateFAQ(context context.Context, faqID string, question *string, answer *string) (bool, error) {
	var _Faq models.Faq

	err := orm.DB.DB.First(&_Faq, "id = ?", faqID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("FAQNotFound")
	}

	if question != nil {
		_Faq.Question = *question
	}

	if answer != nil {
		_Faq.Answer = *answer
	}

	orm.DB.DB.Save(&_Faq)

	// return success
	return true, nil
}

//DeleteFAQ allow admins to delete Tags
func (orm *ORM) DeleteFAQ(context context.Context, faqID string) (bool, error) {
	var _Faq models.Faq

	err := orm.DB.DB.First(&_Faq, "id = ?", faqID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("FAQNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_Faq).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//ReadFAQ calls the single in tag
func (orm *ORM) ReadFAQ(ctx context.Context, faqID string) (*models.Faq, error) {
	var _Faq models.Faq

	err := orm.DB.DB.Joins("CreatedBy").First(&_Faq, "faqs.id = ?", faqID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("FAQNotFound")
	}

	// return success
	return &_Faq, nil
}

//ReadFAQs based on a query
func (orm *ORM) ReadFAQs(ctx context.Context, filterQuery *utils.FilterQuery, question *string, answer *string) ([]*models.Faq, error) {
	var _Faq []*models.Faq

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("faqs.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if question != nil {
		_Results = _Results.Where("faqs.Question = ?", question)
	}

	if answer != nil {
		_Results = _Results.Where("faqs.Answer = ?", answer)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("faqs.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("faqs.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Faq)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Faq, nil
}

//ReadFAQsLength retirieved the count based on a query
func (orm *ORM) ReadFAQsLength(ctx context.Context, filterQuery *utils.FilterQuery, question *string, answer *string) (*int64, error) {
	var _FaqsLength int64

	_Results := orm.DB.DB.Model(&models.Faq{})
	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("faqs.created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}
	if question != nil {
		_Results = _Results.Where("faqs.Question = ?", question)
	}

	if answer != nil {
		_Results = _Results.Where("faqs.Answer = ?", answer)
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("faqs.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("faqs.%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("CreatedBy").
		Count(&_FaqsLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_FaqsLength, nil
}
