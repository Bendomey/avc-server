package schemas

import "github.com/graphql-go/graphql"

// ServicingType defines typings for Servicing Fields - categories
var PaymentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Payment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"servicing": &graphql.Field{
				Type: ServicingType,
			},
			"code": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"authorizationUrl": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"accessCode": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"amount": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
			},
			"subscription": &graphql.Field{
				Type: SubscriptionType,
			},
			"status": &graphql.Field{
				Type: graphql.NewNonNull(EnumForPaymentStatus),
			},
			"createdBy": &graphql.Field{
				Type: UserType,
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

//FilterServicingType  for filtering servicing type
var FilterPaymentType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetPaymentsFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"status": &graphql.InputObjectFieldConfig{
				Type: EnumForPaymentStatus,
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

//EnumServicingStatus  for types of servicing
var EnumForPaymentStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "PaymentStatus",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value: "PENDING",
		},
		"SUCCESS": &graphql.EnumValueConfig{
			Value: "SUCCESS",
		},
		"FAILED": &graphql.EnumValueConfig{
			Value: "FAILED",
		},
	},
})
