package resolvers

import (
	"log"

	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var countriesQuery = map[string]*graphql.Field{
	"countries": {
		Type:        graphql.NewNonNull(graphql.NewList(schemas.CountryType)),
		Description: "Get country list",
		Resolve: utils.AuthenticateAdmin(
			func(params graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
				log.Print(adminData.ID)
				return nil, nil
			},
		),
	},
	"country": {
		Type:        schemas.CountryType,
		Description: "Get single country",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	},
}

var countriesMutation = map[string]*graphql.Field{
	"createCountry": {
		Type:        graphql.NewNonNull(schemas.CountryType),
		Description: "Create country",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	},
}

var country = ResolverLoader{
	Query:    countriesQuery,
	Mutation: countriesMutation,
}
