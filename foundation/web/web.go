package web

import (
	"context"
	"net/http"
)

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, args ...any)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	log Logger
	*http.ServeMux
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log Logger) *App {
	mux := http.NewServeMux()

	return &App{
		log:      log,
		ServeMux: mux,
	}
}
