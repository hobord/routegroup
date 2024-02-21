package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf("index of %s", r.RequestURI)

	w.Write([]byte(content))
}

func hello(w http.ResponseWriter, r *http.Request) {
	nameParam := r.PathValue("name")

	content := fmt.Sprintf("Hello %s", nameParam)

	w.Write([]byte(content))
}

func test(w http.ResponseWriter, r *http.Request) {
	groupParam := r.PathValue("subgroup")
	nameParam := r.PathValue("name")

	content := fmt.Sprintf("Hello, group:%s,  name:%s", groupParam, nameParam)

	w.Write([]byte(content))
}

func BeforeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" before middleware ran, "))
		h.ServeHTTP(w, r)
	})
}

func AfterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		w.Write([]byte(", after middleware ran "))
	})
}



func panicTest(w http.ResponseWriter, r *http.Request) {
	panic("v any")
}
