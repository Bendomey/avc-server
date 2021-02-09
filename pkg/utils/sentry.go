package utils

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

var dsn, environment string

func init() {
	dsn = MustGet("SENTRY_DSN")
	environment = MustGet("SENTRY_ENVIRONMENT")
}

//Sentry starts
func Sentry() {
	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: dsn,
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		Environment: environment,
		Release:     "avc@1.0.0",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
}
