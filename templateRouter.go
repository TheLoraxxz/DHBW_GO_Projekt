package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		body, err := io.ReadAll(request.Body)
		fmt.Println(body)
		if err != nil {
			log.Fatal("Coudnt export Parsefiles")
		}
	}
	mainRoute, err := template.ParseFiles("./assets/templates/index.html", "./assets/templates/header.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}
