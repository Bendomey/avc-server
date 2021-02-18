package schemas

import "github.com/graphql-go/graphql"

var enumTypeCustomerType = graphql.NewEnum(graphql.EnumConfig{
	Name: "CustomerType",
	Values: graphql.EnumValueConfigMap{
		"Business": &graphql.EnumValueConfig{
			Value: "Business",
		},
		"Individual": &graphql.EnumValueConfig{
			Value: "Individual",
		},
	},
})

// CustomerType defines typings for customers
var CustomerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Customer",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(UserType),
			},
			"type": &graphql.Field{
				Type: enumTypeCustomerType,
			},
			"tin": &graphql.Field{
				Type: graphql.String,
			},
			"digitalAddress": &graphql.Field{
				Type: graphql.String,
			},
			"addressCountry": &graphql.Field{
				Type: graphql.String,
			},
			"addressCity": &graphql.Field{
				Type: graphql.String,
			},
			"addressStreetName": &graphql.Field{
				Type: graphql.String,
			},
			"addressNumber": &graphql.Field{
				Type: graphql.String,
			},
			"companyName": &graphql.Field{
				Type: graphql.String,
			},
			"companyEntityType": &graphql.Field{
				Type: graphql.String,
			},
			"companyEntityTypeOther": &graphql.Field{
				Type: graphql.String,
			},
			"companyCountryOfRegistration": &graphql.Field{
				Type: graphql.String,
			},
			"companyDateOfRegistration": &graphql.Field{
				Type: graphql.String,
			},
			"companyRegistrationNumber": &graphql.Field{
				Type: graphql.String,
			},

			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)

//FilterCustomerType  for filtering customers
var FilterCustomerType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetCustomersFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"suspended": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"type": &graphql.InputObjectFieldConfig{
				Type: enumTypeCustomerType,
			},
			"search": &graphql.InputObjectFieldConfig{
				Type: SearchType,
			},
			"order": &graphql.InputObjectFieldConfig{
				Type: enumTypeForOrder,
			},
			"orderBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"dateRange": &graphql.InputObjectFieldConfig{
				Type: DateRangeType,
			},
		},
	},
)
