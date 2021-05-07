package resolvers

import (
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

// ResolverLoader defines the reolver object
type ResolverLoader struct {
	Query    map[string]*graphql.Field
	Mutation map[string]*graphql.Field
}

//ExposeSchema sends querytype and mutation type
func ExposeSchema(services services.Services) []*graphql.Object {

	var queriesGathering = []map[string]*graphql.Field{
		ExposeCountryResolver(services).Query,
		ExposeAdminResolver(services).Query,
		ExposeUserResolver(services).Query,
		ExposeLegalAreaResolver(services).Query,
		ExposeNewsletterResolver(services).Query,
		ExposeTagResolver(services).Query,
		ExposeFaqResolver(services).Query,
		ExposeBlogPostResolver(services).Query,
		ExposePackageResolver(services).Query,
	}

	var mutationsGathering = []map[string]*graphql.Field{
		ExposeCountryResolver(services).Mutation,
		ExposeAdminResolver(services).Mutation,
		ExposeUserResolver(services).Mutation,
		ExposeLegalAreaResolver(services).Mutation,
		ExposeNewsletterResolver(services).Mutation,
		ExposeTagResolver(services).Mutation,
		ExposeFaqResolver(services).Mutation,
		ExposeBlogPostResolver(services).Mutation,
		ExposePackageResolver(services).Mutation,
	}

	// QueryType is the main querytype implementation
	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields(utils.GetResolvers(queriesGathering)),
		},
	)

	// MutationType is the main querytype implementation
	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields(utils.GetResolvers(mutationsGathering)),
		},
	)

	return []*graphql.Object{
		queryType, mutationType,
	}
}
