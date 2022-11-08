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
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	http.Handle("/", &root)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
