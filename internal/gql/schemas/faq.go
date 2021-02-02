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
				Type: AdminType,
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
