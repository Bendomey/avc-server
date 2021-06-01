package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis/v8"
	"github.com/kehindesalaam/go-paystack/paystack"
	"gorm.io/gorm"
)

// Servicing inteface holds the Service-databse transactions of this controller
type ServicingService interface {
	CreateServicing(context context.Context, serviceID string, createdBy string, businessCountry *string, businessEntityType *string, businessName *string, businessOwners *string, businessDirectors *string, businessAddress *string, businessNumberOfShares *string, businessInitialCapital *float32, businessIndustry *string, trademarkCountry *string, trademarkOwnershipType *string, trademarkOwners *string, trademarkAddress *string, trademarkClassification *string, trademarkUploads *string, documentType *string, natureOfDocument *string, documentDeadline *time.Time, existingDocuments *string, newDocuments *string) (*models.Servicing, error)
	// UpdateService(context context.Context, serviceID string, name *string, price *float64, description *string, serviceType *string) (bool, error)
	// DeleteService(context context.Context, serviceeID string) (bool, error)
	// ReadService(ctx context.Context, serviceID string) (*models.Service, error)
	// ReadServices(ctx context.Context, filterQuery *utils.FilterQuery, name *string) ([]*models.Service, error)
	// ReadServicesLength(ctx context.Context, filterQuery *utils.FilterQuery, name *string) (*int64, error)
}

func ServicingSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) ServicingService {
	return &ORM{db, rdb, mg}
}

// CreateServicing adds a new servicing to the user's records
func (orm *ORM) CreateServicing(context context.Context, serviceID string, createdBy string, businessCountry *string, businessEntityType *string, businessName *string, businessOwners *string, businessDirectors *string, businessAddress *string, businessNumberOfShares *string, businessInitialCapital *float32, businessIndustry *string, trademarkCountry *string, trademarkOwnershipType *string, trademarkOwners *string, trademarkAddress *string, trademarkClassification *string, trademarkUploads *string, documentType *string, natureOfDocument *string, documentDeadline *time.Time, existingDocuments *string, newDocuments *string) (*models.Servicing, error) {

	//find the service
	var __Service models.Service

	serviceErr := orm.DB.DB.First(&__Service, "created_by_id = ?", createdBy).Error
	if errors.Is(serviceErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("ServiceNotFound")
	}

	//find user
	var __Customer models.Customer

	err := orm.DB.DB.Preload("User").First(&__Customer, "user_id = ?", createdBy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("CustomerNotFound")
	}

	//save service fields
	__business := models.Business{
		CountryID:      businessCountry,
		EntityType:     businessEntityType,
		Name:           businessName,
		Owners:         businessOwners,
		Directors:      businessDirectors,
		Address:        businessAddress,
		NumberOfShares: businessNumberOfShares,
		InitialCapital: businessInitialCapital,
		Industry:       businessIndustry,
	}

	__trademark := models.Trademark{
		CountryID:                 trademarkCountry,
		OwnershipType:             trademarkOwnershipType,
		Owners:                    trademarkOwners,
		Address:                   trademarkAddress,
		ClassificationOfTrademark: trademarkClassification,
		Uploads:                   trademarkUploads,
	}

	__document := models.Document{
		Type:              documentType,
		NatureOfDoc:       natureOfDocument,
		Deadline:          documentDeadline,
		ExistingDocuments: existingDocuments,
		NewDocuments:      newDocuments,
	}

	__ServiceFields := models.ServicingField{
		Business:  __business,
		Trademark: __trademark,
		Document:  __document,
	}

	//create a servicing record
	__Servicing := models.Servicing{
		ServiceID:       serviceID,
		CreatedByID:     createdBy,
		ServiceFieldsID: __ServiceFields.ID.String(),
		LawyerID:        __Customer.LawyerID,
		Cost:            __Service.Price,
	}

	var __UserSubscription models.Subscription
	errSub := orm.DB.DB.First(&__UserSubscription, models.Subscription{
		CreatedByID: createdBy,
		Status:      "ACTIVE",
	}).Error

	//set subscription_id if available
	if errors.Is(errSub, gorm.ErrRecordNotFound) {
		__Servicing.SubscriptionID = nil
	} else {
		sub := __UserSubscription.ID.String()
		__Servicing.SubscriptionID = &sub
	}

	// save the servicing_field record
	if err := orm.DB.DB.Select("Business", "Trademark", "Document").Create(&__ServiceFields).Error; err != nil {
		raven.CaptureError(err, nil)
		return nil, err
	}

	// save the servicing record
	if err := orm.DB.DB.Select("ServiceID", "CreatedByID", "LawyerID", "SubscriptionID").Create(&__Servicing).Error; err != nil {
		raven.CaptureError(err, nil)
		return nil, err
	}

	if errors.Is(errSub, gorm.ErrRecordNotFound) {
		//payment
		__Payment := models.Payment{
			Amount:      *__Service.Price,
			CreatedByID: createdBy,
		}

		serv := __Servicing.ID.String()
		__Payment.ServicingID = &serv

		//initialize the payment
		currency := " USD"
		amountHere := fmt.Sprintf("%f", *__Service.Price)
		ref := __Payment.Code.String()
		response, payErr := utils.InitializePayment(context, paystack.TransactionRequest{
			Amount:    &amountHere,
			Currency:  &currency,
			Reference: &ref,
			Email:     &__Customer.User.Email,
			Metadata:  paystack.Metadata{},
			Channels:  []string{"card"},
			// CallbackURL:       "",
		})
		if payErr != nil {
			raven.CaptureError(payErr, nil)
			fmt.Print(payErr)
		}

		fmt.Print("Payment response", response)
		// __Payment.AuthorizationUrl = response.authorization_url
		__Payment.AuthorizationUrl = ""
		__Payment.AccessCode = ""

		// save the payment record
		if err := orm.DB.DB.Select("Amount", "CreatedByID", "ServicingID").Create(&__Payment).Error; err != nil {
			raven.CaptureError(err, nil)
			return nil, err
		}

	}

	return &__Servicing, nil
}
