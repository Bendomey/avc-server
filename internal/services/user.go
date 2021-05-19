package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/internal/sms"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/Bendomey/goutilities/pkg/generatecode"
	"github.com/Bendomey/goutilities/pkg/signjwt"
	"github.com/Bendomey/goutilities/pkg/validatehash"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserService inteface holds the user-databse transactions of this controller
type UserService interface {
	CreateUser(ctx context.Context, userType string, email string, password string) (*models.User, error)
	LoginUser(ctx context.Context, email string, password string) (*LoginResultUser, error)
	GetMe(ctx context.Context, userID string) (*GetMeType, error)
	ResendUserCode(ctx context.Context, userID string) (*models.User, error)
	VerifyUserEmail(ctx context.Context, userID string, code string) (*LoginResultUser, error)
	DeleteUser(ctx context.Context, userID string) (bool, error)
	SuspendUser(ctx context.Context, userID string, adminID string, reason string) (bool, error)
	RestoreUser(ctx context.Context, userID string) (bool, error)
	SendPhoneVerificationCode(ctx context.Context, phone string) (bool, error)
	VerifyPhoneVerificationCode(ctx context.Context, phone string, code string) (bool, error)
	UpdateUserAndCustomer(
		ctx context.Context,
		userID string,
		lastName *string,
		firstName *string,
		otherNames *string,
		phone *string,
		email *string,
		customerType *string,
		tin *string,
		digitalAddress *string,
		addressCountry *string,
		addressCity *string,
		addressStreetName *string,
		addressNumber *string,
		companyName *string,
		companyEntityType *string,
		companyEntityTypeOther *string,
		companyCountryOfRegistration *string,
		companyDateOfRegistration *time.Time,
		companyRegistrationNumber *string,
	) (*LoginResultUser, error)
	UpdateUserAndLawyer(
		ctx context.Context,
		userID string,
		lastName *string,
		firstName *string,
		otherNames *string,
		phone *string,
		email *string,
		tin *string,
		digitalAddress *string,
		addressCountry *string,
		addressCity *string,
		addressStreetName *string,
		addressNumber *string,
		firstYearOfBarAdmission *string,
		licenseNumber *string,
		nationalIDFront *string,
		nationalIDBack *string,
		bARMembershipCard *string,
		lawCertificate *string,
		CV *string,
		coverLetter *string,
	) (*LoginResultUser, error)
	ReadCustomers(ctx context.Context, filterQuery *utils.FilterQuery, customerType *string, isSuspended *bool) ([]models.Customer, error)
	ReadCustomersLength(ctx context.Context, filterQuery *utils.FilterQuery, customerType *string, isSuspended *bool) (*int64, error)
	ReadLawyers(ctx context.Context, filterQuery *utils.FilterQuery, isApproved *bool) ([]models.Lawyer, error)
	ReadLawyersLength(ctx context.Context, filterQuery *utils.FilterQuery, isApproved *bool) (*int64, error)
}

//LoginResultUser is the typing for returning login successful data to user
type LoginResultUser struct {
	Token    string           `json:"token"`
	User     models.User      `json:"user"`
	Lawyer   *models.Lawyer   `json:"lawyer"`
	Customer *models.Customer `json:"customer"`
}

type GetMeType struct {
	User     models.User      `json:"user"`
	Lawyer   *models.Lawyer   `json:"lawyer"`
	Customer *models.Customer `json:"customer"`
}

// NewUserSvc exposed the ORM to the user functions in the module
func NewUserSvc(db *orm.ORM, rdb *redis.Client, mg mail.MailingService) UserService {
	return &ORM{db, rdb, mg}
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
		raven.CaptureError(saveErr, nil)
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
	subject := "Welcome To African Venture Counsel - Verify Your Account"
	body := fmt.Sprintf("Use this code '%s' as your verification code on our platform ", code)
	go orm.mg.SendTransactionalMail(ctx, subject, body, email)

	return &_User, nil
}

//LoginUser logs in the user
func (orm *ORM) LoginUser(ctx context.Context, email string, password string) (*LoginResultUser, error) {
	var _User models.User

	//check if email is in db
	err := orm.DB.DB.First(&_User, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("UserNotFound")
		}
		return nil, err
	}

	//check if the person is suspended or not
	if _User.SuspendedAt != nil {
		return nil, errors.New("AccountSuspended")
	}

	//since email in db, lets validate hash and then send back
	isSame := validatehash.ValidateCipher(password, _User.Password)
	if isSame == false {
		return nil, errors.New("PasswordIncorrect")
	}

	//sign token
	token, signTokenErrr := signjwt.SignJWT(jwt.MapClaims{
		"id":   _User.ID,
		"type": _User.Type,
	}, utils.MustGet("USER_SECRET"))

	if signTokenErrr != nil {
		return nil, signTokenErrr
	}

	var _Customer models.Customer
	var _Lawyer models.Lawyer

	custErr := orm.DB.DB.First(&_Customer, "user_id = ?", _User.ID).Error
	_ = orm.DB.DB.First(&_Lawyer, "user_id = ?", _User.ID).Error

	if errors.Is(custErr, gorm.ErrRecordNotFound) {

		//if lawyer and approved
		if _User.SetupAt != nil && _Lawyer.ApprovedAt == nil {
			return nil, errors.New("Your application is under review. We will notify you soon ")
		}

		return &LoginResultUser{
			Token:    token,
			User:     _User,
			Lawyer:   &_Lawyer,
			Customer: nil,
		}, nil
	}

	return &LoginResultUser{
		Token:    token,
		User:     _User,
		Lawyer:   nil,
		Customer: &_Customer,
	}, nil

}

//GetMe helps to get all details needed for a single user
func (orm *ORM) GetMe(ctx context.Context, userID string) (*GetMeType, error) {
	var _User models.User

	//check if email is in db
	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("UserNotFound")
		}
		return nil, err
	}
	fmt.Println(_User)
	//check if the person is suspended or not
	if _User.SuspendedAt != nil {
		return nil, errors.New("AccountSuspended")
	}

	var _Customer models.Customer
	var _Lawyer models.Lawyer

	custErr := orm.DB.DB.First(&_Customer, "user_id = ?", _User.ID).Error
	_ = orm.DB.DB.First(&_Lawyer, "user_id = ?", _User.ID).Error

	if errors.Is(custErr, gorm.ErrRecordNotFound) {
		return &GetMeType{
			User:     _User,
			Lawyer:   &_Lawyer,
			Customer: nil,
		}, nil
	}

	return &GetMeType{
		User:     _User,
		Lawyer:   nil,
		Customer: &_Customer,
	}, nil
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
	subject := "Welcome To African Venture Counsel - Verify Your Account"
	body := fmt.Sprintf("Use this code '%s' as your verification code on our platform ", code)
	go orm.mg.SendTransactionalMail(ctx, subject, body, _User.Email)

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
		Token:    token,
		User:     _User,
		Lawyer:   nil,
		Customer: nil,
	}, nil
}

//SendPhoneVerificationCode sends phone code
func (orm *ORM) SendPhoneVerificationCode(ctx context.Context, phone string) (bool, error) {
	var _User models.User

	//check if in db
	_Results := orm.DB.DB.First(&_User, "phone = ?", phone)
	if _Results.RowsAffected > 0 {
		return false, errors.New("PhoneNumberAlreadyExists")
	}

	//generate code
	code := generatecode.GenerateCode(6)

	//save in redis and expire in an hours time
	redisErr := orm.rdb.Set(ctx, fmt.Sprintf("%s", phone), code, 1*time.Hour).Err()
	if redisErr != nil {
		return false, redisErr
	}

	// send code to phone
	log.Println("Generated code for phone :: ", code)
	body := fmt.Sprintf("Use this code %s to verify your phone number on our website", code)
	sms.Send(phone, body)

	return true, nil
}

//VerifyPhoneVerificationCode sends phone code
func (orm *ORM) VerifyPhoneVerificationCode(ctx context.Context, phone string, code string) (bool, error) {

	value, err := orm.rdb.Get(ctx, fmt.Sprintf("%s", phone)).Result()
	if err == redis.Nil {
		return false, errors.New("CodeHasExpired")
	} else if err != nil {
		return false, err
	}

	if value != code {
		return false, errors.New("CodeIncorrect")
	}

	//invalidate the redis data pertaining to this phone
	redisErr := orm.rdb.Set(ctx, fmt.Sprintf("%s", phone), nil, 1*time.Second).Err()
	if redisErr != nil {
		return false, redisErr
	}

	return true, nil
}

//UpdateUserAndCustomer udpates the user and customer
func (orm *ORM) UpdateUserAndCustomer(
	ctx context.Context,
	userID string,
	lastName *string,
	firstName *string,
	otherNames *string,
	phone *string,
	email *string,

	customerType *string,
	tin *string,
	digitalAddress *string,
	addressCountry *string,
	addressCity *string,
	addressStreetName *string,
	addressNumber *string,
	companyName *string,
	companyEntityType *string,
	companyEntityTypeOther *string,
	companyCountryOfRegistration *string,
	companyDateOfRegistration *time.Time,
	companyRegistrationNumber *string,
) (*LoginResultUser, error) {

	var _User models.User

	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("UserNotFound")
	}

	if _User.SetupAt == nil {
		nowTime := time.Now()
		_User.SetupAt = &nowTime
	}

	if lastName != nil {
		_User.LastName = lastName
	}
	if firstName != nil {
		_User.FirstName = firstName
	}
	if otherNames != nil {
		_User.OtherNames = otherNames
	}

	if phone != nil {
		_User.Phone = phone
		nowTime := time.Now()
		_User.PhoneVerifiedAt = &nowTime
	}
	if email != nil {
		_User.Email = *email
		nowTime := time.Now()
		_User.EmailVerifiedAt = &nowTime
	}

	//for customer
	var _Customer models.Customer

	custFetcherr := orm.DB.DB.First(&_Customer, "user_id = ?", userID).Error
	if errors.Is(custFetcherr, gorm.ErrRecordNotFound) {
		return nil, errors.New("CustomerNotFound")
	}

	if customerType != nil {
		_Customer.Type = customerType
	}

	if tin != nil {
		_Customer.TIN = tin
	}

	if digitalAddress != nil {
		_Customer.DigitalAddress = digitalAddress
	}

	if addressCountry != nil {
		_Customer.AddressCountry = addressCountry
	}
	if addressCity != nil {
		_Customer.AddressCity = addressCity
	}
	if addressStreetName != nil {
		_Customer.AddressStreetName = addressStreetName
	}
	if addressNumber != nil {
		_Customer.AddressNumber = addressNumber
	}
	if companyName != nil {
		_Customer.CompanyName = companyName
	}
	if companyEntityType != nil {
		_Customer.CompanyEntityType = companyEntityType
	}
	if companyEntityTypeOther != nil {
		_Customer.CompanyEntityTypeOther = companyEntityTypeOther
	}
	if companyCountryOfRegistration != nil {
		_Customer.CompanyCountryOfRegistration = companyCountryOfRegistration
	}
	if companyDateOfRegistration != nil {
		_Customer.CompanyDateOfRegistration = companyDateOfRegistration
	}
	if companyRegistrationNumber != nil {
		_Customer.CompanyRegistrationNumber = companyRegistrationNumber
	}

	//save em
	orm.DB.DB.Save(&_User)
	orm.DB.DB.Save(&_Customer)

	//sign token
	token, signTokenErrr := signjwt.SignJWT(jwt.MapClaims{
		"id":   _User.ID,
		"type": _User.Type,
	}, utils.MustGet("USER_SECRET"))

	if signTokenErrr != nil {
		return nil, signTokenErrr
	}

	return &LoginResultUser{
		Token:    token,
		User:     _User,
		Lawyer:   nil,
		Customer: &_Customer,
	}, nil
}

//UpdateUserAndLawyer updates user details plus lawyer details
func (orm *ORM) UpdateUserAndLawyer(
	ctx context.Context,
	userID string,
	lastName *string,
	firstName *string,
	otherNames *string,
	phone *string,
	email *string,

	tin *string,
	digitalAddress *string,
	addressCountry *string,
	addressCity *string,
	addressStreetName *string,
	addressNumber *string,
	firstYearOfBarAdmission *string,
	licenseNumber *string,
	nationalIDFront *string,
	nationalIDBack *string,
	bARMembershipCard *string,
	lawCertificate *string,
	CV *string,
	coverLetter *string,
) (*LoginResultUser, error) {
	var _User models.User

	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("UserNotFound")
	}

	if _User.SetupAt == nil {
		nowTime := time.Now()
		_User.SetupAt = &nowTime
	}

	if lastName != nil {
		_User.LastName = lastName
	}

	if otherNames != nil {
		_User.OtherNames = otherNames
	}
	if firstName != nil {
		_User.FirstName = firstName
	}
	if phone != nil {
		_User.Phone = phone
		nowTime := time.Now()
		_User.PhoneVerifiedAt = &nowTime
	}
	if email != nil {
		_User.Email = *email
		nowTime := time.Now()
		_User.EmailVerifiedAt = &nowTime
	}

	//for customer
	var _Lawyer models.Lawyer

	lawyerFetcherr := orm.DB.DB.First(&_Lawyer, "user_id = ?", userID).Error
	if errors.Is(lawyerFetcherr, gorm.ErrRecordNotFound) {
		return nil, errors.New("LawyerNotFound")
	}

	if tin != nil {
		_Lawyer.TIN = tin
	}
	if digitalAddress != nil {
		_Lawyer.DigitalAddress = digitalAddress
	}
	if addressCountry != nil {
		_Lawyer.AddressCountry = addressCountry
	}
	if addressCity != nil {
		_Lawyer.AddressCity = addressCity
	}
	if addressStreetName != nil {
		_Lawyer.AddressStreetName = addressStreetName
	}
	if addressNumber != nil {
		_Lawyer.AddressNumber = addressNumber
	}
	if firstYearOfBarAdmission != nil {
		_Lawyer.FirstYearOfBarAdmission = firstYearOfBarAdmission
	}
	if licenseNumber != nil {
		_Lawyer.LicenseNumber = licenseNumber
	}
	if nationalIDFront != nil {
		_Lawyer.NationalIDFront = nationalIDFront
	}
	if nationalIDBack != nil {
		_Lawyer.NationalIDBack = nationalIDBack
	}
	if bARMembershipCard != nil {
		_Lawyer.BARMembershipCard = bARMembershipCard
	}
	if lawCertificate != nil {
		_Lawyer.LawCertificate = lawCertificate
	}
	if CV != nil {
		_Lawyer.CV = CV
	}
	if coverLetter != nil {
		_Lawyer.CoverLetter = coverLetter
	}

	//save em
	orm.DB.DB.Save(&_User)
	orm.DB.DB.Save(&_Lawyer)

	//sign token
	token, signTokenErrr := signjwt.SignJWT(jwt.MapClaims{
		"id":   _User.ID,
		"type": _User.Type,
	}, utils.MustGet("USER_SECRET"))

	if signTokenErrr != nil {
		return nil, signTokenErrr
	}

	return &LoginResultUser{
		Token:    token,
		User:     _User,
		Lawyer:   &_Lawyer,
		Customer: nil,
	}, nil
}

//DeleteUser deletes a user
func (orm *ORM) DeleteUser(ctx context.Context, userID string) (bool, error) {
	var _User models.User

	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("UserNotFound")
	}
	if _User.Type == "Customer" {
		var _Customer models.Customer
		_ = orm.DB.DB.Where("user_id = ?", userID).Delete(&_Customer).Error
	} else {
		var _Lawyer models.Lawyer
		_ = orm.DB.DB.Where("user_id = ?", userID).Delete(&_Lawyer).Error
	}

	//delete
	delErr := orm.DB.DB.Delete(&_User).Error
	if delErr != nil {
		return false, delErr
	}

	// return success
	return true, nil
}

//SuspendUser suspends the user in question
func (orm *ORM) SuspendUser(ctx context.Context, userID string, adminID string, reason string) (bool, error) {
	var _User models.User

	//find
	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("UserNotFound")
		}
		return false, err
	}

	//update suspendedAt
	updateError := orm.DB.DB.Model(&_User).Updates(map[string]interface{}{
		"SuspendedAt":     time.Now(),
		"SuspendedReason": reason,
		"SuspendedByID":   adminID,
	}).Error
	if updateError != nil {
		return false, updateError
	}

	//send mail plus reason
	subject := "Welcome To African Venture Counsel - Account Suspension"
	body := fmt.Sprintf("Your account has been suspended.\n Reason: %s", reason)
	go orm.mg.SendTransactionalMail(ctx, subject, body, _User.Email)
	return true, nil
}

//RestoreUser restores the user in question
func (orm *ORM) RestoreUser(ctx context.Context, userID string) (bool, error) {
	var _User models.User

	//find
	err := orm.DB.DB.First(&_User, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("UserNotFound")
		}
		return false, err
	}

	//update suspendedAt
	updateError := orm.DB.DB.Model(&_User).Updates(map[string]interface{}{
		"SuspendedAt":     nil,
		"SuspendedReason": nil,
		"SuspendedByID":   nil,
	}).Error
	if updateError != nil {
		return false, updateError
	}

	subject := "Welcome To African Venture Counsel - Account Restoration"
	body := "Your account has been restored."
	go orm.mg.SendTransactionalMail(ctx, subject, body, _User.Email)
	return true, nil
}

//ReadCustomers reads the customers in the system
func (orm *ORM) ReadCustomers(ctx context.Context, filterQuery *utils.FilterQuery, customerType *string, isSuspended *bool) ([]models.Customer, error) {
	var _Customers []models.Customer

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("\"User\".created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if customerType != nil {
		_Results = _Results.Where("customers.Type = ?", customerType)
	}

	if isSuspended != nil {
		suspended := *isSuspended
		if suspended {
			_Results = _Results.Where("\"User\".suspended_at IS NOT NULL")
		} else {
			_Results = _Results.Where("\"User\".suspended_at IS NULL")
		}
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("User").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Customers)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Customers, nil
}

//ReadCustomersLength reads the length customers in the system
func (orm *ORM) ReadCustomersLength(ctx context.Context, filterQuery *utils.FilterQuery, customerType *string, isSuspended *bool) (*int64, error) {
	var _CustomersLength int64

	_Results := orm.DB.DB.Model(&models.Customer{})

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("\"User\".created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if customerType != nil {
		_Results = _Results.Where("customers.Type = ?", customerType)
	}

	if isSuspended != nil {
		suspended := *isSuspended
		if suspended {
			_Results = _Results.Where("\"User\".suspended_at IS NOT NULL")
		} else {
			_Results = _Results.Where("\"User\".suspended_at IS NULL")
		}
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("User").Count(&_CustomersLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_CustomersLength, nil
}

//ReadLawyers fetches lawyers
func (orm *ORM) ReadLawyers(ctx context.Context, filterQuery *utils.FilterQuery, isApproved *bool) ([]models.Lawyer, error) {
	var _Lawyers []models.Lawyer

	_Results := orm.DB.DB

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("\"User\".created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if isApproved != nil {
		approved := *isApproved
		if approved {
			_Results = _Results.Where("lawyers.approved_at IS NOT NULL")
		} else {
			_Results = _Results.Where("lawyers.approved_at IS NULL")
		}
	}

	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("User").
		Order(fmt.Sprintf("%s %s", filterQuery.OrderBy, filterQuery.Order)).
		Limit(filterQuery.Limit).Offset(filterQuery.Skip).
		Find(&_Lawyers)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return _Lawyers, nil
}

//ReadLawyersLength gets the length of lawyers table
func (orm *ORM) ReadLawyersLength(ctx context.Context, filterQuery *utils.FilterQuery, isApproved *bool) (*int64, error) {
	var _LawyersLength int64
	_Results := orm.DB.DB.Model(&models.Lawyer{})

	//add date range if added
	if filterQuery.DateRange != nil {
		_Results = _Results.Where("\"User\".created_at BETWEEN ? AND ?", filterQuery.DateRange.StartTime, filterQuery.DateRange.EndTime)
	}

	if isApproved != nil {
		approved := *isApproved
		if approved {
			_Results = _Results.Where("lawyers.approved_at IS NOT NULL")
		} else {
			_Results = _Results.Where("lawyers.approved_at IS NULL")
		}
	}
	if filterQuery.Search != nil {
		for index, singleCriteria := range filterQuery.Search.SearchFields {
			//if index is o, start to filter
			if index == 0 {
				_Results = _Results.Where(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
				continue
			}
			//more than one make it or so either ways it comes
			_Results = _Results.Or(fmt.Sprintf("\"User\".%s LIKE ?", singleCriteria), fmt.Sprintf("%%%s%%", filterQuery.Search.Criteria))
		}
	}

	//continue the filtration
	_Results = _Results.Joins("User").Count(&_LawyersLength)

	if _Results.Error != nil {
		return nil, _Results.Error
	}
	return &_LawyersLength, nil
}
