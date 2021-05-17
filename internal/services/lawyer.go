package services

import (
	"context"
	"errors"
	"time"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// LawyerService inteface holds the lawyer-databse transactions of this controller
type LawyerService interface {
	ApproveLawyer(ctx context.Context, lawyerID string, adminID string) (bool, error)
}

// NewLawyerSvc exposed the ORM to the lawyer functions in the module
func NewLawyerSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) LawyerService {
	return &ORM{db, rdb, mg}
}

//ApproveLawyer approves the lawyer in question
func (orm *ORM) ApproveLawyer(ctx context.Context, lawyerID string, adminID string) (bool, error) {
	var _Lawyer models.Lawyer

	//find
	err := orm.DB.DB.First(&_Lawyer, "id = ?", lawyerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("LawyerNotFound")
		}
		return false, err
	}

	//update suspendedAt
	updateError := orm.DB.DB.Model(&_Lawyer).Updates(map[string]interface{}{
		"approved_at":    time.Now(),
		"approved_by_id": adminID,
	}).Error

	if updateError != nil {
		return false, updateError
	}

	var _User models.User
	_ = orm.DB.DB.First(&_User, "id = ?", _Lawyer.UserID).Error

	//send mail
	subject := "Welcome To African Venture Counsel - Account Approved"
	body := "Congratulations, Your account has been approved. You can now receive jobs. Here's to your next chapter."
	go orm.mg.SendTransactionalMail(ctx, subject, body, _User.Email)
	return true, nil
}
