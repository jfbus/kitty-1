package kitty

import (
	"net/http"
)

// Router is an interface to wrap different router implementations.
type Router interface {
	// Handle registers a handler to the router.
	// Multiple methods and paths can be configured.
	Handle(method []string, path []string, handler http.Handler)
	// SetNotFoundHandler will sets the NotFound handler.
	SetNotFoundHandler(handler http.Handler)
	// ServeHTTP dispatches the handler registered in the matched route.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// RouterOption sets optional Router overrides.
type RouterOption func(Router) Router

// CustomRouter allows users to inject an alternate Router implementation.
func (s *Server) Router(r Router, opts ...RouterOption) *Server {
	for _, opt := range opts {
		r = opt(r)
	}
	s.mux = r
	return s
}

// StdlibRouter returns a Router based on stdlib.
func StdlibRouter() Router {
	return &stdlibRouter{mux: http.NewServeMux()}
}

// NotFoundHandler will set the not found handler of the router.
func NotFoundHandler(h http.Handler) RouterOption {
	return func(r Router) Router {
		r.SetNotFoundHandler(h)
		return r
	}
}

var _ Router = &stdlibRouter{}

// StdlibRouter is a Router implementation for the Stdlib's `http.ServeMux`.
type stdlibRouter struct {
	mux *http.ServeMux
}

// Handle registers a handler to the router.
func (g *stdlibRouter) Handle(methods, paths []string, h http.Handler) {
	mc := map[string]struct{}{}
	for _, method := range methods {
		mc[method] = struct{}{}
	}
	for _, path := range paths {
		g.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if _, found := mc[r.Method]; found {
				h.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
		})
	}
}

// SetNotFoundHandler will do nothing as we cannot override the Not Found handler from the stdlib.
func (g *stdlibRouter) SetNotFoundHandler(h http.Handler) {
}

// ServeHTTP dispatches the handler registered in the matched route.
func (g *stdlibRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}
