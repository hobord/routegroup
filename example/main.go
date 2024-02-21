package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	root := MakeApi(mux)

	http.ListenAndServe(":8080", root)
}
