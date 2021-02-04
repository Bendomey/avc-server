package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/graphql-go/graphql"
)

var userQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{}
}

var userMutation = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"createUser": {
			Type:        graphql.NewNonNull(schemas.UserType),
			Description: "Create User in the system",
			Args: graphql.FieldConfigArgument{
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeUserType),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				password := p.Args["password"].(string)
				email := p.Args["email"].(string)
				userType := p.Args["type"].(string)

				_Response, err := svcs.UserServices.CreateUser(p.Context, userType, email, password)
				if err != nil {
					return nil, err
				}
				return transformations.DBUserToGQLUser(_Response), nil
			},
		},
		"resendUserCode": {
			Type:        graphql.NewNonNull(schemas.UserType),
			Description: "Generates a new code and sends to user's email",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Args["userId"].(string)

				_Response, err := svcs.UserServices.ResendUserCode(p.Context, userID)
				if err != nil {
					return nil, err
				}
				return transformations.DBUserToGQLUser(_Response), nil
			},
		},
		"verifyUserEmail": {
			Type:        graphql.NewNonNull(schemas.LoginUserType),
			Description: "Verifies the user's email by comparing save code and user code",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"code": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Args["userId"].(string)
				code := p.Args["code"].(string)

				_Response, err := svcs.UserServices.VerifyUserEmail(p.Context, userID, code)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"user":  transformations.DBUserToGQLUser(&_Response.User),
					"token": _Response.Token,
				}, nil
			},
		},
	}
}

// ExposeUserResolver exposes the admin resolver
func ExposeUserResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    userQuery(services),
		Mutation: userMutation(services),
	}
}
