package resolvers

import (
	"time"

	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var userQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{}
}

var userMutation = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"createUser": {
			Type:        graphql.NewNonNull(schemas.UserType),
			Description: "Create User in the system",
			Args: graphql.FieldConfigArgument{
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeUserType),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				password := p.Args["password"].(string)
				email := p.Args["email"].(string)
				userType := p.Args["type"].(string)

				_Response, err := svcs.UserServices.CreateUser(p.Context, userType, email, password)
				if err != nil {
					return nil, err
				}
				return transformations.DBUserToGQLUser(_Response), nil
			},
		},
		"loginUser": {
			Type:        graphql.NewNonNull(schemas.LoginUserType),
			Description: "Create User in the system",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				password := p.Args["password"].(string)
				email := p.Args["email"].(string)

				_Response, err := svcs.UserServices.LoginUser(p.Context, email, password)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"user":     transformations.DBUserToGQLUser(&_Response.User),
					"lawyer":   transformations.DBUserToGQLLawyer(_Response.Lawyer),
					"customer": transformations.DBUserToGQLCustomer(_Response.Customer),
					"token":    _Response.Token,
				}, nil
			},
		},
		"resendUserCode": {
			Type:        graphql.NewNonNull(schemas.UserType),
			Description: "Generates a new code and sends to user's email",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Args["userId"].(string)

				_Response, err := svcs.UserServices.ResendUserCode(p.Context, userID)
				if err != nil {
					return nil, err
				}
				return transformations.DBUserToGQLUser(_Response), nil
			},
		},
		"verifyUserEmail": {
			Type:        graphql.NewNonNull(schemas.LoginUserType),
			Description: "Verifies the user's email by comparing save code and user code",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"code": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Args["userId"].(string)
				code := p.Args["code"].(string)

				_Response, err := svcs.UserServices.VerifyUserEmail(p.Context, userID, code)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"user":  transformations.DBUserToGQLUser(&_Response.User),
					"token": _Response.Token,
				}, nil
			},
		},
		"sendPhoneVerificationCode": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Generates a new code and sends to user's phone",
			Args: graphql.FieldConfigArgument{
				"phone": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				phone := p.Args["phone"].(string)

				_Response, err := svcs.UserServices.SendPhoneVerificationCode(p.Context, phone)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"verifyPhoneCode": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Verifies if code sumbitted is same as what is saved",
			Args: graphql.FieldConfigArgument{
				"phone": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"code": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				phone := p.Args["phone"].(string)
				code := p.Args["code"].(string)

				_Response, err := svcs.UserServices.VerifyPhoneVerificationCode(p.Context, phone, code)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"updateUserAndCustomer": {
			Type:        graphql.NewNonNull(schemas.LoginUserType),
			Description: "Updates Customer Details Plus Users Details ... They Are Connected",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.UpdateCustomerInput),
				},
			},
			Resolve: utils.AuthenticateUser(
				func(p graphql.ResolveParams, loggedInUser *utils.UserFromToken) (interface{}, error) {
					takeInput, inputOk := p.Args["input"].(map[string]interface{})
					var lastName, firstName, otherNames, phone, email, customerType, tin,
						digitalAddress, addressCountry, addressCity, addressStreetName, addressNumber, companyName,
						companyEntityType, companyEntityTypeOther, companyCountryOfRegistration,
						companyRegistrationNumber *string

					var companyDateOfRegistration *time.Time
					if inputOk {
						takeLastName, lastNameOk := takeInput["lastName"].(string)
						if lastNameOk {
							lastName = &takeLastName
						} else {
							lastName = nil
						}
						takefirstName, firstNameOk := takeInput["firstName"].(string)
						if firstNameOk {
							firstName = &takefirstName
						} else {
							firstName = nil
						}
						takeOtherName, otherNameOk := takeInput["otherNames"].(string)
						if otherNameOk {
							otherNames = &takeOtherName
						} else {
							otherNames = nil
						}
						takePhone, phoneOk := takeInput["phone"].(string)
						if phoneOk {
							phone = &takePhone
						} else {
							phone = nil
						}
						takeEmail, emailOk := takeInput["email"].(string)
						if emailOk {
							email = &takeEmail
						} else {
							email = nil
						}
						takeCustomerType, customerTypeOk := takeInput["customerType"].(string)
						if customerTypeOk {
							customerType = &takeCustomerType
						} else {
							customerType = nil
						}
						takeTin, TinOk := takeInput["tin"].(string)
						if TinOk {
							tin = &takeTin
						} else {
							tin = nil
						}
						takeDigitalAddress, DigitalAddressOk := takeInput["digitalAddress"].(string)
						if DigitalAddressOk {
							digitalAddress = &takeDigitalAddress
						} else {
							digitalAddress = nil
						}
						takeaddressCountry, addressCountryOk := takeInput["addressCountry"].(string)
						if addressCountryOk {
							addressCountry = &takeaddressCountry
						} else {
							addressCountry = nil
						}
						takeaddressCity, addressCityOk := takeInput["addressCity"].(string)
						if addressCityOk {
							addressCity = &takeaddressCity
						} else {
							addressCity = nil
						}
						takeaddressStreetName, addressStreetNameOk := takeInput["addressStreetName"].(string)
						if addressStreetNameOk {
							addressStreetName = &takeaddressStreetName
						} else {
							addressStreetName = nil
						}
						takeaddressNumber, addressNumberOk := takeInput["addressNumber"].(string)
						if addressNumberOk {
							addressNumber = &takeaddressNumber
						} else {
							addressNumber = nil
						}
						takecompanyName, companyNameOk := takeInput["companyName"].(string)
						if companyNameOk {
							companyName = &takecompanyName
						} else {
							companyName = nil
						}
						takecompanyEntityType, companyEntityTypeOk := takeInput["companyEntityType"].(string)
						if companyEntityTypeOk {
							companyEntityType = &takecompanyEntityType
						} else {
							companyEntityType = nil
						}
						takecompanyEntityTypeOther, companyEntityTypeOtherOk := takeInput["companyEntityTypeOther"].(string)
						if companyEntityTypeOtherOk {
							companyEntityTypeOther = &takecompanyEntityTypeOther
						} else {
							companyEntityTypeOther = nil
						}
						takecompanyCountryOfRegistration, companyCountryOfRegistrationOk := takeInput["companyCountryOfRegistration"].(string)
						if companyCountryOfRegistrationOk {
							companyCountryOfRegistration = &takecompanyCountryOfRegistration
						} else {
							companyCountryOfRegistration = nil
						}
						takecompanyDateOfRegistration, companyDateOfRegistrationOk := takeInput["companyDateOfRegistration"].(time.Time)
						if companyDateOfRegistrationOk {
							companyDateOfRegistration = &takecompanyDateOfRegistration
						} else {
							companyDateOfRegistration = nil
						}
						takecompanyRegistrationNumber, companyRegistrationNumberOk := takeInput["companyRegistrationNumber"].(string)
						if companyRegistrationNumberOk {
							companyRegistrationNumber = &takecompanyRegistrationNumber
						} else {
							companyRegistrationNumber = nil
						}
					}

					_Response, err := svcs.UserServices.UpdateUserAndCustomer(p.Context, loggedInUser.ID, lastName, firstName,
						otherNames, phone, email, customerType, tin, digitalAddress, addressCountry, addressCity, addressStreetName, addressNumber,
						companyName, companyEntityType, companyEntityTypeOther, companyCountryOfRegistration, companyDateOfRegistration, companyRegistrationNumber,
					)
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"user":     transformations.DBUserToGQLUser(&_Response.User),
						"lawyer":   transformations.DBUserToGQLLawyer(_Response.Lawyer),
						"customer": transformations.DBUserToGQLCustomer(_Response.Customer),
						"token":    _Response.Token,
					}, nil

				},
			),
		},
		"updateUserAndLawyer": {
			Type:        graphql.NewNonNull(schemas.LoginUserType),
			Description: "Updates Lawyer Details Plus Users Details ... They Are Connected",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.UpdateLawyerInput),
				},
			},
			Resolve: utils.AuthenticateUser(
				func(p graphql.ResolveParams, loggedInUser *utils.UserFromToken) (interface{}, error) {
					takeInput, inputOk := p.Args["input"].(map[string]interface{})
					var lastName, firstName, otherNames, phone, email, tin,
						digitalAddress, addressCountry, addressCity, addressStreetName, addressNumber,
						firstYearOfBarAdmission, licenseNumber, nationalIDFront, nationalIDBack, BARMembershipCard, lawCertificate,
						CV, coverLetter *string

					if inputOk {
						takeLastName, lastNameOk := takeInput["lastName"].(string)
						if lastNameOk {
							lastName = &takeLastName
						} else {
							lastName = nil
						}
						takefirstName, firstNameOk := takeInput["firstName"].(string)
						if firstNameOk {
							firstName = &takefirstName
						} else {
							firstName = nil
						}
						takeOtherName, otherNameOk := takeInput["otherNames"].(string)
						if otherNameOk {
							otherNames = &takeOtherName
						} else {
							otherNames = nil
						}
						takePhone, phoneOk := takeInput["phone"].(string)
						if phoneOk {
							phone = &takePhone
						} else {
							phone = nil
						}
						takeEmail, emailOk := takeInput["email"].(string)
						if emailOk {
							email = &takeEmail
						} else {
							email = nil
						}
						takeTin, TinOk := takeInput["tin"].(string)
						if TinOk {
							tin = &takeTin
						} else {
							tin = nil
						}
						takeDigitalAddress, DigitalAddressOk := takeInput["digitalAddress"].(string)
						if DigitalAddressOk {
							digitalAddress = &takeDigitalAddress
						} else {
							digitalAddress = nil
						}
						takeaddressCountry, addressCountryOk := takeInput["addressCountry"].(string)
						if addressCountryOk {
							addressCountry = &takeaddressCountry
						} else {
							addressCountry = nil
						}
						takeaddressCity, addressCityOk := takeInput["addressCity"].(string)
						if addressCityOk {
							addressCity = &takeaddressCity
						} else {
							addressCity = nil
						}
						takeaddressStreetName, addressStreetNameOk := takeInput["addressStreetName"].(string)
						if addressStreetNameOk {
							addressStreetName = &takeaddressStreetName
						} else {
							addressStreetName = nil
						}
						takeaddressNumber, addressNumberOk := takeInput["addressNumber"].(string)
						if addressNumberOk {
							addressNumber = &takeaddressNumber
						} else {
							addressNumber = nil
						}
						takefirstYearOfBarAdmission, firstYearOfBarAdmissionOk := takeInput["firstYearOfBarAdmission"].(string)
						if firstYearOfBarAdmissionOk {
							firstYearOfBarAdmission = &takefirstYearOfBarAdmission
						} else {
							firstYearOfBarAdmission = nil
						}
						takelicenseNumber, licenseNumberOk := takeInput["licenseNumber"].(string)
						if licenseNumberOk {
							licenseNumber = &takelicenseNumber
						} else {
							licenseNumber = nil
						}
						takenationalIDFront, nationalIDFrontOk := takeInput["nationalIDFront"].(string)
						if nationalIDFrontOk {
							nationalIDFront = &takenationalIDFront
						} else {
							nationalIDFront = nil
						}
						takenationalIDBack, nationalIDBackOk := takeInput["nationalIDBack"].(string)
						if nationalIDBackOk {
							nationalIDBack = &takenationalIDBack
						} else {
							nationalIDBack = nil
						}
						takeBARMembershipCard, BARMembershipCardOk := takeInput["BARMembershipCard"].(string)
						if BARMembershipCardOk {
							BARMembershipCard = &takeBARMembershipCard
						} else {
							BARMembershipCard = nil
						}
						takelawCertificate, lawCertificateOk := takeInput["lawCertificate"].(string)
						if lawCertificateOk {
							lawCertificate = &takelawCertificate
						} else {
							lawCertificate = nil
						}
						takeCV, CVOk := takeInput["CV"].(string)
						if CVOk {
							CV = &takeCV
						} else {
							CV = nil
						}
						takecoverLetter, coverLetterOk := takeInput["coverLetter"].(string)
						if coverLetterOk {
							coverLetter = &takecoverLetter
						} else {
							coverLetter = nil
						}
					}

					_Response, err := svcs.UserServices.UpdateUserAndLawyer(p.Context, loggedInUser.ID, lastName, firstName, otherNames, phone, email, tin, digitalAddress,
						addressCountry, addressCity, addressStreetName, addressNumber, firstYearOfBarAdmission, licenseNumber,
						nationalIDFront, nationalIDBack, BARMembershipCard, lawCertificate, CV, coverLetter,
					)
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"user":     transformations.DBUserToGQLUser(&_Response.User),
						"lawyer":   transformations.DBUserToGQLLawyer(_Response.Lawyer),
						"customer": transformations.DBUserToGQLCustomer(_Response.Customer),
						"token":    _Response.Token,
					}, nil
				},
			),
		},
		"deleteUser": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Deletes user",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					userID := p.Args["userId"].(string)

					_Response, err := svcs.UserServices.DeleteUser(p.Context, userID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"suspendUser": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Suspend User",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"reason": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					userID := p.Args["userId"].(string)
					reason := p.Args["reason"].(string)

					_Response, err := svcs.UserServices.SuspendUser(p.Context, userID, adminData.ID, reason)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"restoreUser": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Restore User",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					userID := p.Args["userId"].(string)

					_Response, err := svcs.UserServices.RestoreUser(p.Context, userID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeUserResolver exposes the admin resolver
func ExposeUserResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    userQuery(services),
		Mutation: userMutation(services),
	}
}
