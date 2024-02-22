package routegroup

import "net/http"

type Option func(*group)

// WithMux sets the mux for the group.  (default: http.NewServeMux()
func WithMux(mux *http.ServeMux) Option {
	return func(g *group) {
		g.ServeMux = mux
	}
}

// WithMiddlewares adds the middlewares for the group.
func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) Option {
	return func(g *group) {
		g.Middlewares = append(g.Middlewares, middlewares...)
	}
}

// WithPrefix sets the prefix for the group from root level.
func WithPrefix(prefix string) Option {
	return func(g *group) {
		g.Prefix = prefix
	}
}

// WithSubPrefix sets the sub prefix for the group from parent level.
func WithSubPrefix(prefix string) Option {
	return func(g *group) {
		g.Prefix = g.Prefix + prefix
	}
}

// WithRegisterCallback sets the callback that is called when a route is registered.
func WithRegisterCallback(callback func(pattern string)) Option {
	return func(g *group) {
		g.RegisterCallback = callback
	}
}

// WithRegisterPanicCallback sets the callback that is called when a panic occurs while registering a route.
func WithRegisterPanicCallback(callback func(pattern string, err interface{})) Option {
	return func(g *group) {
		g.RegisterPanicCallback = callback
	}
}
