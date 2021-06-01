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
	"github.com/go-redis/redis/v8"
	"github.com/kehindesalaam/go-paystack/paystack"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	SubscribeToPackage(context context.Context, packageID string, numberOfMonths int, createdBy string) (*models.Payment, error)
}

func SubscriptionSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) SubscriptionService {
	return &ORM{db, rdb, mg}
}

//SubscribeToPackage subscribes a user to a package
func (orm *ORM) SubscribeToPackage(context context.Context, packageID string, numberOfMonths int, createdBy string) (*models.Payment, error) {
	//save subscription
	__subscription := models.Subscription{
		PackageID:   packageID,
		CreatedByID: createdBy,
		SubscribeAt: time.Now(),
		ExpiresAt:   time.Now().AddDate(0, numberOfMonths, 0), // add the number of months (1 or 12)
	}

	//find package
	var __package models.Package
	packageErr := orm.DB.DB.First(&__package, "id = ?", packageID).Select("AmountPerMonth", "AmountPerYear").Error
	if errors.Is(packageErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("PackageNotFound")
	}

	//find user
	var __Customer models.User

	err := orm.DB.DB.First(&__Customer, "id = ?", createdBy).Select("Email").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("UserNotFound")
	}

	//generate payment link
	__payment := models.Payment{
		CreatedByID: createdBy,
	}

	var amount float64
	if numberOfMonths == 1 {
		amount = float64(*__package.AmountPerMonth)
	} else {
		amount = float64(*__package.AmountPerYear)
	}
	__payment.Amount = amount

	serv := __subscription.ID.String()
	__payment.SubscriptionID = &serv

	//initialize the payment
	currency := "GHS"
	amountHere := fmt.Sprintf("%f", __payment.Amount)
	ref := __payment.Code.String()
	response, payErr := utils.InitializePayment(context, paystack.TransactionRequest{
		Amount:    &amountHere,
		Currency:  &currency,
		Reference: &ref,
		Email:     &__Customer.Email,
		Metadata:  paystack.Metadata{},
		Channels:  []string{"card"},
		// CallbackURL:       "",
	})
	if payErr != nil {
		// raven.CaptureError(payErr, nil)
		fmt.Print(payErr)
		return nil, payErr
	}

	fmt.Print("Payment response", response)
	__payment.AuthorizationUrl = *response.AuthorizationUrl
	__payment.AccessCode = *response.AccessCode

	//save
	if subErr := orm.DB.DB.Select("PackageID", "CreatedByID", "SubscribeAt", "ExpiresAt").Create(&__subscription).Error; subErr != nil {
		return nil, subErr
	}

	if err := orm.DB.DB.Select("CreatedByID", "Amount", "SubscriptionID", "AuthorizationUrl", "AccessCode").Create(&__payment).Error; err != nil {
		// raven.CaptureError(err, nil)
		fmt.Print(payErr)
		return nil, err
	}

	// return payment
	return &__payment, nil
}
