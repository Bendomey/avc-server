package schemas

import "github.com/graphql-go/graphql"

// AdminType defines typings for administrators
var AdminType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Administrator",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"fullname": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"emailVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			// "createdBy": &graphql.Field{
			// 	Type: &graphql.NonNull{s},
			// },
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)
