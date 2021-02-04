package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/Bendomey/goutilities/pkg/generatecode"
	"github.com/Bendomey/goutilities/pkg/signjwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserService inteface holds the user-databse transactions of this controller
type UserService interface {
	CreateUser(ctx context.Context, userType string, email string, password string) (*models.User, error)
	ResendUserCode(ctx context.Context, userID string) (*models.User, error)
	VerifyUserEmail(ctx context.Context, userID string, code string) (*LoginResultUser, error)
}

//LoginResultUser is the typing for returning login successful data to user
type LoginResultUser struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// NewUserSvc exposed the ORM to the user functions in the module
func NewUserSvc(db *orm.ORM, rdb *redis.Client) UserService {
	return &ORM{db, rdb}
}

// CreateUser creates a user when invoked
func (orm *ORM) CreateUser(ctx context.Context, userType string, email string, password string) (*models.User, error) {
	var _User models.User

	//check if email is table or not
	_Response := orm.DB.DB.First(&_User, "email = ?", email)
	if _Response.RowsAffected != 0 {
		return nil, errors.New("UserAlreadyExists")
	}

	_User = models.User{
		Type:     userType,
		Email:    email,
		Password: password,
	}

	//create user
	saveErr := orm.DB.DB.Select("Type", "Email", "Password").Create(&_User).Error
	if saveErr != nil {
		return nil, saveErr
	}

	//create newsletter
	var _NewsLetter models.NewsletterSubscribers
	// if it already exists, then it means they've already appliced as a newletter subscriber
	_NewsLetterResponseFetch := orm.DB.DB.First(&_NewsLetter, "email = ?", email)
	if _NewsLetterResponseFetch.RowsAffected != 0 {
		// update their type to user, since anyone created with the newsletter field is anonymous
		_NewsLetter.Type = "User"
		orm.DB.DB.Save(&_NewsLetter)
	} else {
		_NewSubscriber := models.NewsletterSubscribers{
			Type:  "User",
			Email: email,
		}
		_ = orm.DB.DB.Select("Type", "Email").Create(&_NewSubscriber).Error
	}

	//create corresponding user from laywer or customer table
	if _User.Type == "Customer" {
		//create a customer record
		_Customer := models.Customer{
			UserID: _User.ID.String(),
		}
		_ = orm.DB.DB.Select("UserID").Create(&_Customer).Error
	} else {
		//create a lawyer record
		_Lawyer := models.Lawyer{
			UserID: _User.ID.String(),
		}
		_ = orm.DB.DB.Select("UserID").Create(&_Lawyer).Error
	}

	//generate an otp here
	code := generatecode.GenerateCode(6)

	//save in redis and expire in an hours time
	redisErr := orm.rdb.Set(ctx, fmt.Sprintf("%s", _User.ID), code, 1*time.Hour).Err()
	if redisErr != nil {
		return nil, redisErr
	}

	// send code to email
	log.Println("Generated code :: ", code)

	return &_User, nil
}

//ResendUserCode helps to resend a new code
func (orm *ORM) ResendUserCode(ctx context.Context, userID string) (*models.User, error) {
	var _User models.User

	// check if admin exists
	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("UserNotFound")
	}

	//generate code
	code := generatecode.GenerateCode(6)

	//save in redis and expire in an hours time
	redisErr := orm.rdb.Set(ctx, fmt.Sprintf("%s", _User.ID), code, 1*time.Hour).Err()
	if redisErr != nil {
		return nil, redisErr
	}

	// send code to email
	log.Println("Generated code :: ", code)

	return &_User, nil
}

// VerifyUserEmail compares the user code sent by user
func (orm *ORM) VerifyUserEmail(ctx context.Context, userID string, code string) (*LoginResultUser, error) {
	//check in redis to see if its the same and not expired
	value, err := orm.rdb.Get(ctx, fmt.Sprintf("%s", userID)).Result()
	if err == redis.Nil {
		return nil, errors.New("CodeHasExpired")
	} else if err != nil {
		return nil, err
	}

	if value != code {
		return nil, errors.New("CodeIncorrect")
	}

	// update user emailVerifiedAt record
	var _User models.User
	fetchUserErr := orm.DB.DB.First(&_User, "id = ?", userID).Error

	if fetchUserErr != nil {
		return nil, fetchUserErr
	}

	currentTime := time.Now()
	_User.EmailVerifiedAt = &currentTime
	orm.DB.DB.Save(&_User)

	//sign token
	token, signTokenErrr := signjwt.SignJWT(jwt.MapClaims{
		"id":   _User.ID,
		"type": _User.Type,
	}, utils.MustGet("USER_SECRET"))

	if signTokenErrr != nil {
		return nil, signTokenErrr
	}

	//invalidate the redis data pertaining to this user
	redisErr := orm.rdb.Set(ctx, fmt.Sprintf("%s", userID), nil, 1*time.Second).Err()
	if redisErr != nil {
		return nil, redisErr
	}

	return &LoginResultUser{
		Token: token,
		User:  _User,
	}, nil
}
