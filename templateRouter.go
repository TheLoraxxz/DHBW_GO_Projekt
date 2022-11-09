package main

import (
	"html/template"
	"log"
	"net/http"
)

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

	}
	mainRoute, err := template.ParseFiles("./assets/sites/index.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}
