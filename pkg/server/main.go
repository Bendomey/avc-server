package server

import (
	"net/http"

	"github.com/Bendomey/avc-server/internal/handlers"
	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/pkg/utils"
)

var gqlPgPath, port string

func init() {
	port = utils.MustGet("PORT")
	gqlPgPath = utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH")
}

// Run web server
func Run() {
	h := handlers.CreateGQLServer()

	// // Handlers
	// Simple keep-alive/ping handler
	http.Handle("/ping", handlers.Ping())

	//serve a grapqhl endpoint at /graphql
	http.Handle(gqlPgPath, handlers.PlaygroundHanlder(h))

	//and serve!
	port := ":5000"
	log.NewLogger().Printf(`GraphQL server starting up on http://localhost%v`, port)
	errServer := http.ListenAndServe(port, nil)
	if errServer != nil {
		log.Fatalf("Error occured while serving graphql server, %v", errServer)
	}
}
