package services

import (
	"context"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/go-redis/redis/v8"
)

// NewsletterService inteface holds the newsletter-databse transactions of this controller
type NewsletterService interface {
	SubscribeToNewsletter(context context.Context, email string) (bool, error)
}

// NewsletterSvc exposed the ORM to the newsletter functions in the module
func NewsletterSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) NewsletterService {
	return &ORM{db, rdb, mg}
}

//SubscribeToNewsletter allow users on the website to subscribe to newsletters
func (orm *ORM) SubscribeToNewsletter(context context.Context, email string) (bool, error) {
	_Newsletter := models.NewsletterSubscribers{
		Email: email,
		Type:  "Anon",
	}

	err := orm.DB.DB.Select("Email", "Type").Create(&_Newsletter).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
