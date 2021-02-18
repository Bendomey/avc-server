package services

import (
	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/go-redis/redis/v8"
)

// Services responsible for exposing all services to resolvers
type Services struct {
	AdminServices     AdminService
	CountryServices   CountryService
	UserServices      UserService
	LawyerServices    LawyerService
	LegalAreaServices LegalAreaService
}

//ORM gets orm connection
type ORM struct {
	DB  *orm.ORM
	rdb *redis.Client
	mg  mail.MailingService
}

//Factory activates all services
func Factory(orm *orm.ORM, rdb *redis.Client, mg mail.MailingService) Services {
	//activate admin service
	adminService := NewAdminSvc(orm, rdb, mg)
	countryService := NewCountrySvc(orm, rdb, mg)
	userService := NewUserSvc(orm, rdb, mg)
	lawyerService := NewLawyerSvc(orm, rdb, mg)
	legalAreaService := NewLegalAreaSvc(orm, rdb, mg)

	return Services{
		AdminServices:     adminService,
		CountryServices:   countryService,
		UserServices:      userService,
		LawyerServices:    lawyerService,
		LegalAreaServices: legalAreaService,
	}
}
