package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/graphql-go/graphql"
)

// QueryType is the main querytype implementation
var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"countries": &graphql.Field{
				Type:        graphql.NewList(schemas.CountryType),
				Description: "Get country list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	},
)
