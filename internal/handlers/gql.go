package handlers

import (
	"net/http"

	"github.com/Bendomey/avc-server/internal/gql/resolvers"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
	handler "github.com/graphql-go/graphql-go-handler"
)

var schema graphql.Schema
var pgEnabled bool

func init() {
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: resolvers.QueryType,
		},
	)
	pgEnabled = utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED")
}

//CreateGQLServer creates a grapqhl playground
func CreateGQLServer() *handler.Handler {
	//create a graphql-go http handler with our defined schema
	// and also set to return a pretty JSON
	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: pgEnabled,
	})
}

//PlaygroundHanlder registers a route for plaground access
func PlaygroundHanlder(handler *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
