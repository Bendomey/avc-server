package resolvers

import (
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/graphql-go/graphql"
)

var newsletterQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{}
}

var newsletterMutation = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"subscribeToNewsletter": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Subscribe To Newsletter",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)

				_Response, err := svcs.NewsletterServices.SubscribeToNewsletter(p.Context, email)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
	}
}

// ExposeNewsletterResolver exposes the newsletter resolver
func ExposeNewsletterResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    newsletterQuery(services),
		Mutation: newsletterMutation(services),
	}
}
