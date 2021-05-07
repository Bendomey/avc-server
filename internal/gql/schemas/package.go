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

//FilterPackagesType  for filtering Package type
var FilterPackagesType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetPackagesFilter",
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
