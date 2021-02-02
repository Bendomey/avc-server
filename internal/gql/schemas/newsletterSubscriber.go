package schemas

import "github.com/graphql-go/graphql"

// NewsletterSubscriber defines typings for newsletter subscriber
var NewsletterSubscriber = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "NewsletterSubscriber",
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