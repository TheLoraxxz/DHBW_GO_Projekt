package main

import (
	"fmt"
	"log"
	"net/http"
)

type RootHandler struct {
}

var Server = http.Server{
	Addr: ":80",
}

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello!")
}

func main() {
	root := RootHandler{}
	http.Handle("/", &root)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
