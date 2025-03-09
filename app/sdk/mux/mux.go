// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) http.Handler {

	app := http.NewServeMux()

	return app
}
