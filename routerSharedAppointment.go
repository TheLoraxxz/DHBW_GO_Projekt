package main

import (
	"html/template"
	"log"
	"net/http"
)

func AdminSiteServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-admin.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
}

func CreateLinkServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-create-link.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
}

func ServeHTTPSharedAppCreateDate(writer http.ResponseWriter, request *http.Request) {
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-create-app.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}

}
