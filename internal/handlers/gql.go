package handlers

import (
	"context"
	"net/http"

	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
	handler "github.com/graphql-go/graphql-go-handler"
)

var pgEnabled bool

func init() {

	pgEnabled = utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED")
}

//CreateGQLServer creates a grapqhl playground
func CreateGQLServer(r []*graphql.Object) *handler.Handler {
	schema, _ := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    r[0],
			Mutation: r[1],
		},
	)
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
		setupResponse(&w, r)
		ctx := context.WithValue(r.Context(), utils.GetPrincipalID(), r.Header.Get("Authorization"))
		handler.ContextHandler(ctx, w, r)
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
