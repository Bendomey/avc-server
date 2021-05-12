package schemas

import "github.com/graphql-go/graphql"

// PackageServiceType defines typings for Package Services
var PackageServiceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PackageService",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"service": &graphql.Field{
				Type: graphql.NewNonNull(ServiceType),
			},
			"package": &graphql.Field{
				Type: graphql.NewNonNull(PackageType),
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(EnumTypeForPackageService),
			},
			"quantity": &graphql.Field{
				Type: graphql.Int,
			},
			"isActive": &graphql.Field{
				Type: graphql.Boolean,
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
var FilterPackageServicesType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetPackageServicesFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"package": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			},
			"service": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
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
var EnumTypeForPackageService = graphql.NewEnum(graphql.EnumConfig{
	Name: "PackageServiceType",
	Values: graphql.EnumValueConfigMap{
		"BOOLEAN": &graphql.EnumValueConfig{
			Value: "BOOLEAN",
		},
		"NUMBER": &graphql.EnumValueConfig{
			Value: "NUMBER",
		},
	},
})
