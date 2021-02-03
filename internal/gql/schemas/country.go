package schemas

import "github.com/graphql-go/graphql"

//CountryType defines typings for country
var CountryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Country",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"currency": &graphql.Field{
				Type: graphql.String,
			},
			"createdBy": &graphql.Field{
				Type: graphql.NewNonNull(AdminType),
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

//FilterCountryType  for filtering countries type
var FilterCountryType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Filter",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"description": &graphql.InputObjectFieldConfig{
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
