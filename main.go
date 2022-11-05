package main

import (
	"fmt"
	"log"
	"net/http"
)

type RootHandler struct {
}

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello!")
}

func main() {
	root := RootHandler{}
	server := http.Server{
		Addr: ":80",
	}
	http.Handle("/", &root)
	if err := server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}
}
