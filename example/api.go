package main

import (
	"net/http"

	"github.com/hobord/routegroup"
)

func MakeApi(mux *http.ServeMux) http.Handler {
	// create a root group with the provided mux
	// the register callback will be called
	// when a route is registered
	// the register panic callback will be called
	// when a panic occurs while registering a route
	// the recover middleware will be added to the root group
	// that will catch panics and log them
	root := routegroup.NewGroup(
		routegroup.WithMux(mux),
		routegroup.WithRegisterCallback(routegroup.RegisterRouteCallback),
		routegroup.WithRegisterPanicCallback(routegroup.RegisterPanicHandler),
		routegroup.WithMiddlewares(routegroup.Recover),
	)

	// index page handler
	root.HandleFunc("GET /{$}", index)

	// You can add middlewares to the group any time
	// but they are not applied to already registered routes
	// just to the future ones
	// the logger middleware will log all requests (expect the index page)
	// all middlewares are run in the order they are added
	// all middlewares are inherited by sub groups
	root.Use(routegroup.Logger)

	// it couse a panic becouse path is already registered
	// but the registred panic callback will recover the panic
	root.HandleFunc("GET /{$}", index)

	// hello page handlerFunc with a parameter
	// and a wraped with middlewares
	root.HandleFunc("GET /hello/{name}/",
		AfterMiddleware(
			BeforeMiddleware(
				http.HandlerFunc(hello),
			),
		).ServeHTTP,
	)

	// example of a route with panic,
	// the recover middleware will catch the panic and log it
	root.HandleFunc("GET /panic", panicTest)

	// create a route group
	{
		group := root.SubGroup("/group")

		// add middlewares to the group
		group.Use(BeforeMiddleware, AfterMiddleware)

		// group index page handler
		group.HandleFunc("GET /{$}", index)

		// group hello page handler
		group.HandleFunc("GET /{name}/", hello)

		// create a sub group
		{
			subGroup := group.SubGroup("/{subgroup}")

			// you can clean up the inherited middlewares
			subGroup.Middlewares = nil

			// and fresh middlewares to the sub group
			subGroup.Use(routegroup.Recover, routegroup.Logger)

			// sub group index page handler
			subGroup.HandleFunc("GET /{name}", test)
		}

	}

	return root
}
