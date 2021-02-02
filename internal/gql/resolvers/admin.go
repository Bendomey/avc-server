package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/graphql-go/graphql"
)

var adminQuery = map[string]*graphql.Field{
	"admins": {
		Type:        graphql.NewNonNull(graphql.NewList(schemas.AdminType)),
		Description: "Get Administrators in the system",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	},
}

var adminMutation = map[string]*graphql.Field{
	"createAdmin": {
		Type:        graphql.NewNonNull(schemas.AdminType),
		Description: "Create Admin in the system",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	},
}

var admin = ResolverLoader{
	Query:    adminQuery,
	Mutation: adminMutation,
}
