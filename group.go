package routegroup

import (
	"fmt"
	"net/http"
	"strings"
)

// Group is a group of routes.
// It can be used to add middlewares to the group of routes at once.
// It can help to manage the prefixes of the routes easily.
type Group interface {
	// Handle registers a new handler with the given pattern.
	Handle(pattern string, handler http.Handler)

	// HandleFunc registers a new handler function with the given pattern.
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	// Use adds a middlewares to the group.
	Use(middlewares ...func(http.Handler) http.Handler) Group

	// Clone creates a new group that is a copy of the current group.
	Clone(options ...Option) Group

	// SubGroup creates a new group that is a sub group of the current group.
	// It takes a prefix and a list of options that can be used to configure the new group.
	// The prefix is used to create a new sub pattern for the new group.
	// The middlewares of the current group are inherited to the new group.
	SubGroup(prefix string, options ...Option) Group
}

type group struct {
	// ServeMux is the http.ServeMux that is used to register the routes.
	*http.ServeMux

	// Prefix is the prefix of the group.
	// It is used to create a new pattern for the group.
	Prefix string

	// Middlewares is a list of middlewares that are applied to the group.
	Middlewares []func(http.Handler) http.Handler

	RegisterCallback      func(pattern string)
	RegisterPanicCallback func(pattern string, err interface{})
}

// NewGroup creates a new group.
// It takes a list of options that can be used to configure the group.
func NewGroup(options ...Option) *group {
	g := &group{
		Prefix:      "",
		Middlewares: []func(http.Handler) http.Handler{},
	}
	for _, option := range options {
		option(g)
	}

	if g.ServeMux == nil {
		g.ServeMux = http.NewServeMux()
	}

	return g
}

// Clone creates a new group that is a copy of the current group.
// It takes a list of options that can be used to configure the new group.
// The middlewares of the current group are inherited by the new group.
// The prefix of the current group is inherited by the new group.
// The ServeMux of the current group is inherited by the new group.
// The RegisterCallback of the current group is inherited by the new group.
// The RegisterPanicCallback of the current group is inherited by the new group.
func (g *group) Clone(options ...Option) *group {
	ng := NewGroup(
		WithMux(g.ServeMux),
		WithPrefix(g.Prefix),
		WithMiddlewares(g.Middlewares...),
		WithRegisterCallback(g.RegisterCallback),
		WithRegisterPanicCallback(g.RegisterPanicCallback),
	)

	for _, option := range options {
		option(ng)
	}

	return ng
}

// SubGroup creates a new group that is a sub group of the current group.
// It takes a prefix and a list of options that can be used to configure the new group.
// The prefix is used to create a new sub pattern for the new group.
// The middlewares of the current group are inherited by the new group.
func (g *group) SubGroup(prefix string, options ...Option) *group {
	options = append(options, WithSubPrefix(prefix))
	return g.Clone(options...)
}

// Use adds a middlewares to the group.
func (g *group) Use(middlewares ...func(http.Handler) http.Handler) *group {
	g.Middlewares = append(g.Middlewares, middlewares...)
	return g
}

// Handle registers a new handler with the given pattern.
func (g *group) Handle(pattern string, handler http.Handler) {
	defer func() {
		if err := recover(); err != nil {
			if g.RegisterPanicCallback == nil {
				panic(err)
			}

			g.RegisterPanicCallback(pattern, err)

			return
		}

		if g.RegisterCallback != nil {
			g.RegisterCallback(pattern)
		}
	}()

	for _, middleware := range g.Middlewares {
		handler = middleware(handler)
	}

	if g.Prefix != "" {
		sp := strings.Split(pattern, " ")

		if len(sp) == 1 {
			pattern = g.Prefix + pattern
		} else {
			pattern = fmt.Sprintf("%s %s%s", sp[0], g.Prefix, sp[1])
		}
	}

	g.ServeMux.Handle(pattern, handler)
}

// HandleFunc registers a new handler function with the given pattern.
func (g *group) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	g.Handle(pattern, handlerFn)
}
