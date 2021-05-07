package schemas

import "github.com/graphql-go/graphql"

// ServiceType defines typings for Services - categories
var ServiceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Service",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(EnumTypeForService),
			},
			"createdBy": &graphql.Field{
				Type: AdminType,
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

//FilterServicesType  for filtering Service type
var FilterServicesType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetServicesFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
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

//enumTypeForService for ordering
var EnumTypeForService = graphql.NewEnum(graphql.EnumConfig{
	Name: "ServiceType",
	Values: graphql.EnumValueConfigMap{
		"BOTH": &graphql.EnumValueConfig{
			Value: "BOTH",
		},
		"SUBSCRIBE": &graphql.EnumValueConfig{
			Value: "SUBSCRIBE",
		},
		"UNSUBSCRIBE": &graphql.EnumValueConfig{
			Value: "UNSUBSCRIBE",
		},
	},
})
