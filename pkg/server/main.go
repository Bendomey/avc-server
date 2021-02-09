package server

import (
	"fmt"
	"net/http"

	"github.com/Bendomey/avc-server/internal/gql/resolvers"
	"github.com/Bendomey/avc-server/internal/handlers"
	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/getsentry/raven-go"
)

var gqlPgPath, port string

func init() {
	port = utils.MustGet("PORT")
	gqlPgPath = utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH")
}

// Run web server
func Run(services services.Services) {
	r := resolvers.ExposeSchema(services)
	h := handlers.CreateGQLServer(r)

	// // Handlers
	// Simple keep-alive/ping handler
	http.Handle("/", handlers.Ping())

	//serve a grapqhl endpoint at /graphql
	http.Handle(gqlPgPath, handlers.PlaygroundHanlder(h))

	//and serve!
	log.NewLogger().Printf(`[GraphQL] :: Server started successfully on http://localhost:%v%s`, port, gqlPgPath)
	errServer := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if errServer != nil {
		raven.CaptureError(errServer, nil)
		log.Fatalf("Error occured while serving graphql server, %v", errServer)
	}
}
