package services

import (
	"github.com/Bendomey/avc-server/internal/mail"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/go-redis/redis/v8"
)

// Services responsible for exposing all services to resolvers
type Services struct {
	AdminServices          AdminService
	CountryServices        CountryService
	UserServices           UserService
	LawyerServices         LawyerService
	LegalAreaServices      LegalAreaService
	NewsletterServices     NewsletterService
	TagServices            TagService
	FaqServices            FAQService
	BlogPostServices       BlogPostService
	PackageServices        PackageService
	ServiceServices        ServiceService
	PackageServiceServices PackageServiceService
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
	newsletterService := NewsletterSvc(orm, rdb, mg)
	tagService := TagSvc(orm, rdb, mg)
	faqService := FAQSvc(orm, rdb, mg)
	postService := BlogPostSvc(orm, rdb, mg)
	packageService := PackageSvc(orm, rdb, mg)
	serviceService := ServiceSvc(orm, rdb, mg)
	packageServiceService := PackageServiceSvc(orm, rdb, mg)

	return Services{
		AdminServices:          adminService,
		CountryServices:        countryService,
		UserServices:           userService,
		LawyerServices:         lawyerService,
		LegalAreaServices:      legalAreaService,
		NewsletterServices:     newsletterService,
		TagServices:            tagService,
		FaqServices:            faqService,
		BlogPostServices:       postService,
		PackageServices:        packageService,
		ServiceServices:        serviceService,
		PackageServiceServices: packageServiceService,
	}
}
