package schemas

import "github.com/graphql-go/graphql"

// ServicingType defines typings for Servicing Fields - categories
var SubscriptionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"package": &graphql.Field{
				Type: PackageType,
			},
			// "payment": &graphql.Field{
			// 	Type: PaymentType,
			// },
			"status": &graphql.Field{
				Type: graphql.NewNonNull(EnumForSubscriptionStatus),
			},
			"subscribedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"expiresAt": &graphql.Field{
				Type: graphql.DateTime,
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
var FilterSubscriptionType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetSubscriptionsFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"status": &graphql.InputObjectFieldConfig{
				Type: EnumForSubscriptionStatus,
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
var EnumForSubscriptionStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "SubscriptionStatus",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value: "PENDING",
		},
		"ACTIVE": &graphql.EnumValueConfig{
			Value: "ACTIVE",
		},
		"EXPIRED": &graphql.EnumValueConfig{
			Value: "EXPIRED",
		},
	},
})
