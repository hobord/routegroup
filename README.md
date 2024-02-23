# RouteGroup

[![GoDoc](https://godoc.org/github.com/hobord/routegroup?status.svg)](https://godoc.org/github.com/hobord/routegroup)
[![Go Report Card](https://goreportcard.com/badge/github.com/hobord/routegroup)](https://goreportcard.com/report/github.com/hobord/routegroup)

Small package or snippets for extend http.ServeMux. It adds route grouping feature.
It helps to organize routes in a better way.
You can add middlewares to multiple routes at once.
Also, you can add a prefix to multiple routes at once.
You add callback function that runs when new route added to the group (build api documentation, loging, etc.), and you can add calback function to catch panics on handler registration too. (You have registred route with pattern already...)

Take a look at the [detailed example](example) to see how to use it.


```go
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
```

## Options

```go
// WithMux sets the mux for the group. (default: http.NewServeMux()
func WithMux(mux *http.ServeMux) Option 

// WithMiddlewares adds the middlewares for the group.
func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) Option

// WithPrefix sets the prefix for the group from root level.
func WithPrefix(prefix string) Option

// WithSubPrefix sets the sub prefix for the group from parent level.
func WithSubPrefix(prefix string) Option

// WithRegisterCallback sets the callback that is called when a route is registered.
func WithRegisterCallback(callback func(pattern string)) Option

// WithRegisterPanicCallback sets the callback that is called when a panic occurs while registering a route.
func WithRegisterPanicCallback(callback func(pattern string, err interface{})) Option
```

## Short example:
```go
package main

import (
	"net/http"

	"github.com/hobord/routegroup"
)

func main() {
    mux := http.NewServeMux()

    root := routegroup.NewGroup(
		routegroup.WithMux(mux),
		routegroup.WithMiddlewares(routegroup.Recover),
	)

	// index page handler
	root.HandleFunc("GET /$", index)

    subgroup := root.SubGroup("/group")

    // add middlewares to the sub group
    subgroup.Use(MiddlewareOne, OtherMiddleware)

    // subgroup index page handler -> /group/
    subgroup.HandleFunc("GET /", groupindex)

    http.ListenAndServe(":8080", root)
}
```
