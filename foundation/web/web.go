// Package web contains a small web framework extension.
package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Encoder defines behavior that can encode a data model and provide
// the content type for that encoding.
type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

// HandlerFunc represents a function that handles a http request within our own
// little mini framework.
type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, args ...any)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	log Logger
	*http.ServeMux
	mw      []MidFunc
	origins []string
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log Logger, mw ...MidFunc) *App {

	mux := http.NewServeMux()

	return &App{
		log:      log,
		ServeMux: mux,
		mw:       mw,
	}
}

// EnableCORS enables CORS preflight requests to work in the middleware. It
// prevents the MethodNotAllowedHandler from being called. This must be enabled
// for the CORS middleware to work.
func (a *App) EnableCORS(origins []string) {
	a.origins = origins

	handler := func(ctx context.Context, r *http.Request) Encoder {
		return nil
	}
	handler = wrapMiddleware([]MidFunc{a.corsHandler}, handler)

	a.HandlerFuncNoMid(http.MethodOptions, "", "/", handler)
}

func (a *App) corsHandler(webHandler HandlerFunc) HandlerFunc {
	h := func(ctx context.Context, r *http.Request) Encoder {
		w := GetWriter(ctx)

		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		return webHandler(ctx, r)
	}

	return h
}

// HandlerFunc sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandlerFunc(method string, group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	handlerFunc = wrapMiddleware(mw, handlerFunc)
	handlerFunc = wrapMiddleware(a.mw, handlerFunc)

	if a.origins != nil {
		handlerFunc = wrapMiddleware([]MidFunc{a.corsHandler}, handlerFunc)
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
		}
		ctx := setValues(r.Context(), &v)
		ctx = setWriter(ctx, w)

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
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

// HandlerFunc sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandlerFuncNoMid(method string, group string, path string, handlerFunc HandlerFunc) {
	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
		}
		ctx := setValues(r.Context(), &v)
		ctx = setWriter(ctx, w)

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
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
