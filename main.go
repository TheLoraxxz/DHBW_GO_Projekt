package main

import (
	"DHBW_GO_Projekt/assets/templates"
	"log"
	"net/http"
)

type RootHandler struct {
}

var Server = http.Server{
	Addr: ":80",
}

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	templates.TemplTest.ExecuteTemplate(writer, "page", nil)
}

func main() {
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	http.Handle("/", &root)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
