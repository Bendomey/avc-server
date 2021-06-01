package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var subscriptionsQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{}
}

var subscriptionsMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"subscribeToPackage": {
			Type:        graphql.NewNonNull(schemas.PaymentType),
			Description: "Subscribe to package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"numberOfMonths": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: utils.AuthenticateUser(
				func(p graphql.ResolveParams, userData *utils.UserFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)
					numberOfMonths := p.Args["numberOfMonths"].(int)

					_Response, err := svcs.SubscriptionServices.SubscribeToPackage(p.Context, packageID, numberOfMonths, userData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBPaymentToGQLUser(_Response), nil
				},
			),
		},
	}
}

// ExposeTagResolver exposes the tags Reesolver
func ExposeSubscriptionResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    subscriptionsQuery(services),
		Mutation: subscriptionsMutation(services),
	}
}
