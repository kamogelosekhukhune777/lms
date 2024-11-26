package web

import (
	"context"
	"fmt"
	"net/http"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

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

// HandlerFunc sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandlerFunc(method string, group string, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()

		if err := handler(ctx, w, r); err != nil {
			a.log(ctx, "web-respond", "ERROR", err)
			return
		}

	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.ServeMux.HandleFunc(finalPath, h)
}
