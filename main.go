package main

import (
	"html/template"
	"log"
	"net/http"
)

type RootHandler struct {
}

var Server = http.Server{
	Addr: ":80",
}

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mainRoute, err := template.ParseFiles("./assets/templates/index.html", "./assets/templates/header.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}

func main() {
	//hier weitere handler hinzufügen in ähnlicher fashion für die verschiedenen Templates
	root := RootHandler{}
	http.Handle("/", &root)
	if err := Server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatal(err)
	}

}
