package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func AdminSiteServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isAllowed, user := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host+"/", http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	terminShared, err := terminfindung.GetTerminFromShared(&user, &termin)
	terminForHTML := terminShared.ChangeToCorrectHTML()
	if err != nil {
		log.Fatal("Coudnt find Termin")

	}
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-admin.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, terminForHTML)
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
	isAllowed, _ := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	if len(termin) == 0 {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}

	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, termin)
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
func CreateTest(writer http.ResponseWriter, request *http.Request) {
	user := "admin"
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		user, "test")
	terminID, err := terminfindung.CreateSharedTermin(&termin, &user)
	if err != nil {
		return
	}
	fmt.Printf(terminID)
	_, err = fmt.Fprintf(writer, "Test finished")
	if err != nil {
		return
	}

}
