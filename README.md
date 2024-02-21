# RouteGroup

Small package or snippets for std lib `http` package to handle route grouping.
It helps to organize routes in a better way.
You can add middlewares to multiple routes at once.
Also, you can add a prefix to multiple routes at once.

Take a look at the [example](example) to see how to use it.

quick example:
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
