package schemas

import "github.com/graphql-go/graphql"

var enumTypeTestUserType = graphql.NewEnum(graphql.EnumConfig{
	Name: "UserType",
	Values: graphql.EnumValueConfigMap{
		"Customer": &graphql.EnumValueConfig{
			Value: 0,
		},
		"Lawyer": &graphql.EnumValueConfig{
			Value: 1,
		},
	},
})

// UserType defines typings for users
var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"type": &graphql.Field{
				Type: graphql.NewNonNull(enumTypeTestUserType),
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"otherNames": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"emailVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"phoneVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
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
