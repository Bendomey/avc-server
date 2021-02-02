package resolvers

import (
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

// ResolverLoader defines the reolver object
type ResolverLoader struct {
	Query    map[string]*graphql.Field
	Mutation map[string]*graphql.Field
}

var queriesGathering = []map[string]*graphql.Field{
	country.Query,
	admin.Query,
}

var mutationsGathering = []map[string]*graphql.Field{
	country.Mutation,
	admin.Mutation,
}

// QueryType is the main querytype implementation
var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields(utils.GetReolvers(queriesGathering)),
	},
)

// MutationType is the main querytype implementation
var MutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields(utils.GetReolvers(mutationsGathering)),
	},
)
