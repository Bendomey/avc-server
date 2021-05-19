package schemas

import "github.com/graphql-go/graphql"

// PackageType defines typings for Packages - categories
var PackageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Package",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"amountPerYear": &graphql.Field{
				Type: graphql.Int,
			},
			"amountPerMonth": &graphql.Field{
				Type: graphql.Int,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: EnumForPackageStatus,
			},
			"createdBy": &graphql.Field{
				Type: AdminType,
			},
			"requestedBy": &graphql.Field{
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

//EnumPackageStatus  for types of servicing
var EnumForPackageStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "PackageStatus",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value: "PENDING",
		},
		"APPROVED": &graphql.EnumValueConfig{
			Value: "APPROVED",
		},
	},
})

//CustomPackageServices for adding custom Package services
var CustomPackageServices = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "CustomPackageServices",
		Fields: graphql.InputObjectConfigFieldMap{
			"serviceId": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			},
			"quantity": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"isActive": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
	},
)

//FilterPackagesType  for filtering Package type
var FilterPackagesType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetPackagesFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"type": &graphql.InputObjectFieldConfig{
				Type: EnumForPackageType,
			},
			"user": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			},
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

//EnumPackageType  for types of servicing
var EnumForPackageType = graphql.NewEnum(graphql.EnumConfig{
	Name: "PackageType",
	Values: graphql.EnumValueConfigMap{
		"MAIN": &graphql.EnumValueConfig{
			Value: "MAIN",
		},
		"REQUESTED": &graphql.EnumValueConfig{
			Value: "REQUESTED",
		},
	},
})
