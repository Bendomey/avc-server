package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/Bendomey/goutilities/pkg/hashpassword"
	"github.com/Bendomey/goutilities/pkg/signjwt"
	"github.com/Bendomey/goutilities/pkg/validatehash"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// AdminService inteface holds the admin-databse transactions of this controller
type AdminService interface {
	CreateAdmin(ctx context.Context, name string, email string, role string, reatedBy *string) (*models.Admin, error)
	LoginAdmin(ctx context.Context, email string, password string) (*LoginResult, error)
	UpdateAdminPassword(ctx context.Context, adminID string, oldPassword string, newPassword string) (bool, error)
	UpdateAdmin(ctx context.Context, adminID string, fullname *string, email *string) (bool, error)
	UpdateAdminPhone(ctx context.Context, adminID string, phone string) (bool, error)
	DeleteAdmin(ctx context.Context, adminID string) (bool, error)
	SuspendAdmin(ctx context.Context, user string, admin string, reason string) (bool, error)
	RestoreAdmin(ctx context.Context, adminID string) (bool, error)
}

//ORM gets orm connection
type ORM struct {
	DB  *orm.ORM
	rdb *redis.Client
}

//LoginResult is the typing for returning login successful data to user
type LoginResult struct {
	Token string       `json:"token"`
	Admin models.Admin `json:"admin"`
}

// NewAdminSvc exposed the ORM to the admin functions in the module
func NewAdminSvc(db *orm.ORM, rdb *redis.Client) AdminService {
	return &ORM{db, rdb}
}

// LoginAdmin checks if the email is having valid credentials and returns them a unique, secured token to help them get resources from app
func (orm *ORM) LoginAdmin(ctx context.Context, email string, password string) (*LoginResult, error) {
	var _Admin models.Admin

	//check if email is in db
	_Result := orm.DB.DB.Joins("CreatedBy").First(&_Admin, "admins.email = ?", email)
	if errors.Is(_Result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("AdminNotFound")
	}

	//check if the person is suspended or not
	if _Admin.SuspendedAt != nil {
		return nil, errors.New("AccountSuspended")
	}

	//since email in db, lets validate hash and then send back
	isSame := validatehash.ValidateCipher(password, _Admin.Password)
	if isSame == false {
		return nil, errors.New("PasswordIncorrect")
	}

	//sign token
	token, signTokenErrr := signjwt.SignJWT(jwt.MapClaims{
		"id":   _Admin.ID,
		"role": _Admin.Role,
	}, utils.MustGet("ADMIN_SECRET"))

	if signTokenErrr != nil {
		return nil, signTokenErrr
	}
	log.Println(_Admin.Password)
	return &LoginResult{
		Token: token,
		Admin: _Admin,
	}, nil
}

// CreateAdmin creates an admin when invoked
func (orm *ORM) CreateAdmin(ctx context.Context, name string, email string, role string, createdBy *string) (*models.Admin, error) {
	//generate password
	password := utils.GenerateRandomString(10)
	log.Println(password)
	_Admin := models.Admin{
		FullName:    name,
		Email:       email,
		Password:    password,
		CreatedByID: createdBy,
		Role:        role,
	}

	// create admin
	_Result := orm.DB.DB.Select("FullName", "Email", "Password", "CreatedByID", "Role").Create(&_Admin)
	if _Result.Error != nil {
		return nil, _Result.Error
	}

	//send welcome message to email plus new generated password

	//return admin as response
	return &_Admin, nil
}

// UpdateAdminPassword updates password of an admin
func (orm *ORM) UpdateAdminPassword(ctx context.Context, adminID string, oldPassword string, newPassword string) (bool, error) {
	var _Admin models.Admin

	err := orm.DB.DB.First(&_Admin, "id = ?", adminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}

	isSame := validatehash.ValidateCipher(oldPassword, _Admin.Password)
	if isSame == false {
		return false, errors.New("OldPasswordIncorrect")
	}
	//hash new password
	hashed, hashErr := hashpassword.HashPassword(newPassword)
	if hashErr != nil {
		return false, hashErr
	}
	updateError := orm.DB.DB.Model(&_Admin).Update("password", hashed).Error
	if updateError != nil {
		return false, updateError
	}
	return true, nil

}

// UpdateAdmin updates data of an admin
func (orm *ORM) UpdateAdmin(ctx context.Context, adminID string, fullname *string, email *string) (bool, error) {
	var _Admin models.Admin

	err := orm.DB.DB.First(&_Admin, "id = ?", adminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}

	if fullname != nil {
		_Admin.FullName = *fullname
	}

	if email != nil {
		_Admin.Email = *email
	}
	orm.DB.DB.Save(&_Admin)

	return true, nil
}

// UpdateAdminPhone updates data of an admin
func (orm *ORM) UpdateAdminPhone(ctx context.Context, adminID string, phone string) (bool, error) {
	var _Admin models.Admin

	err := orm.DB.DB.First(&_Admin, "id = ?", adminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}
	updateError := orm.DB.DB.Model(&_Admin).Updates(map[string]interface{}{"phone": phone, "phone_verified_at": time.Now()}).Error
	if updateError != nil {
		return false, updateError
	}
	return true, nil
}

// DeleteAdmin deletes an admin
func (orm *ORM) DeleteAdmin(ctx context.Context, adminID string) (bool, error) {
	var _Admin models.Admin
	//find
	err := orm.DB.DB.First(&_Admin, "id = ?", adminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}

	//delete
	delErr := orm.DB.DB.Delete(&_Admin).Error
	if delErr != nil {
		return false, delErr
	}

	return true, nil
}

//SuspendAdmin suspends the admin in question
func (orm *ORM) SuspendAdmin(ctx context.Context, user string, admin string, reason string) (bool, error) {
	var _Admin models.Admin

	//find
	err := orm.DB.DB.First(&_Admin, "id = ?", user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}

	//update suspendedAt
	updateError := orm.DB.DB.Model(&_Admin).Updates(map[string]interface{}{
		"SuspendedAt":     time.Now(),
		"SuspendedReason": reason,
		"SuspendByID":     admin,
	}).Error
	if updateError != nil {
		return false, updateError
	}

	//send mail plus reason
	return true, nil
}

//RestoreAdmin suspends the admin in question
func (orm *ORM) RestoreAdmin(ctx context.Context, adminID string) (bool, error) {
	var _Admin models.Admin

	//find
	err := orm.DB.DB.First(&_Admin, "id = ?", adminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("AdminNotFound")
	}

	//update suspendedAt
	updateError := orm.DB.DB.Model(&_Admin).Updates(map[string]interface{}{
		"SuspendedAt":     nil,
		"SuspendedReason": nil,
		"SuspendByID":     nil,
	}).Error
	if updateError != nil {
		return false, updateError
	}
	return true, nil
}