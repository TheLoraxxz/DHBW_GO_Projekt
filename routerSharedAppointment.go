package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"html/template"
	"log"
	"net/http"
)

func AdminSiteServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isallowed, _ := checkIfIsAllowed(request)
	if !isallowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
	}
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
func checkIfIsAllowed(request *http.Request) (isAllowed bool, username string) {
	cookie, err := request.Cookie("SessionID-Kalender")
	//if cookie is not existing it returns back to the host
	if err != nil {
		isAllowed = false
		return
	}
	isAllowed, username = authentifizierung.CheckCookie(&cookie.Value)
	return
}
