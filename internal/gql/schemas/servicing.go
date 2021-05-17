package schemas

import "github.com/graphql-go/graphql"

// ServicingType defines typings for Servicing Fields - categories
var ServicingType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Servicing",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"service": &graphql.Field{
				Type: ServiceType,
			},
			"cost": &graphql.Field{
				Type: graphql.Float,
			},
			// "payment": &graphql.Field{
			// 	Type: PaymentType,
			// },
			"subscription": &graphql.Field{
				Type: SubscriptionType,
			},
			"status": &graphql.Field{
				Type: graphql.NewNonNull(EnumForServicingStatus),
			},
			"lawyer": &graphql.Field{
				Type: UserType,
			},
			"serviceFields": &graphql.Field{
				Type: ServicingFieldType,
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
var FilterServicingType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetServicingsFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"status": &graphql.InputObjectFieldConfig{
				Type: EnumForServicingStatus,
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
var EnumForServicingStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "ServicingStatus",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value: "PENDING",
		},
		"PAID": &graphql.EnumValueConfig{
			Value: "PAID",
		},
		"ACTIVE": &graphql.EnumValueConfig{
			Value: "ACTIVE",
		},
		"DONE": &graphql.EnumValueConfig{
			Value: "DONE",
		},
	},
})
