package schemas

import "github.com/graphql-go/graphql"

//EnumTypeUserType defines enum for useer type
var EnumTypeUserType = graphql.NewEnum(graphql.EnumConfig{
	Name: "UserType",
	Values: graphql.EnumValueConfigMap{
		"Customer": &graphql.EnumValueConfig{
			Value: "Customer",
		},
		"Lawyer": &graphql.EnumValueConfig{
			Value: "Laywer",
		},
	},
})

// UserType defines typings for users
var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(EnumTypeUserType),
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"otherNames": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"emailVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"setupAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"phoneVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)

//LoginUserType defines the response on successfull login
var LoginUserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "UserLoginResult",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(UserType),
			},
			"lawyer": &graphql.Field{
				Type: LawyerType,
			},
			"customer": &graphql.Field{
				Type: CustomerType,
			},
		},
	},
)

//UpdateCustomerInput defines input for updating
var UpdateCustomerInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UpdateCustomerInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"lastName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"firstName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"otherNames": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"phone": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"type": &graphql.InputObjectFieldConfig{
				Type: EnumTypeUserType,
			},
			"tin": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"digitalAddress": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressCountry": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressCity": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressStreetName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressNumber": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyEntityType": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyEntityTypeOther": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyCountryOfRegistration": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"companyDateOfRegistration": &graphql.InputObjectFieldConfig{
				Type: graphql.DateTime,
			},
			"companyRegistrationNumber": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

//UpdateLawyerInput defines input for updating
var UpdateLawyerInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UpdateLawyerInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"lastName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"firstName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"otherNames": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"phone": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"tin": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"digitalAddress": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressCountry": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressCity": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressStreetName": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"addressNumber": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"firstYearOfBarAdmission": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"licenseNumber": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"nationalIDFront": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"nationalIDBack": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"BARMembershipCard": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"lawCertificate": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"CV": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"coverLetter": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)
