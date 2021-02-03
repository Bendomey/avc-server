package schemas

import "github.com/graphql-go/graphql"

// EnumTypeAdminRole defines the enum type
var EnumTypeAdminRole = graphql.NewEnum(graphql.EnumConfig{
	Name: "AdminRole",
	Values: graphql.EnumValueConfigMap{
		"Admin": &graphql.EnumValueConfig{
			Value: "Admin",
		},
		"User": &graphql.EnumValueConfig{
			Value: "User",
		},
	},
})

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
			"role": &graphql.Field{
				Type: graphql.NewNonNull(EnumTypeAdminRole),
			},
			"phoneVerifiedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			// "createdBy": &graphql.Field{
			// 	Type: AdminType,
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

//LoginAdminType defines the response on successfull login
var LoginAdminType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdminLoginResult",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"admin": &graphql.Field{
				Type: graphql.NewNonNull(AdminType),
			},
		},
	},
)
