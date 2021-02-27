package schemas

import "github.com/graphql-go/graphql"

// FaqType defines typings for frequently asked questions
var FaqType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Faq",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"question": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"answer": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdBy": &graphql.Field{
				Type: graphql.NewNonNull(AdminType),
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

//FilterFaqsType  for filtering faqs type
var FilterFaqsType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetFAQsFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"question": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"answer": &graphql.InputObjectFieldConfig{
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
