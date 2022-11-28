package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

func AdminSiteServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isAllowed, user := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host+"/", http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	selectedDay := request.URL.Query().Get("selected")
	terminShared, err := terminfindung.GetTerminFromShared(&user, &termin)
	if err != nil {
		log.Fatal("Coudnt find Termin")

	}
	if len(selectedDay) != 0 {
		err := terminfindung.SelectDate(&selectedDay, &termin, &user)
		if err != nil {
			return
		}
		terminShared, err = terminfindung.GetTerminFromShared(&user, &termin)
		if err != nil {
			return
		}
	}
	terminForHTML := terminShared.ConvertAdminToHTML()

	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-admin.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal(err)
	}
	err = mainRoute.Execute(writer, terminForHTML)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateLinkServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isAllowed, user := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host+"/", http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			return
		}
		name := request.Form.Get("name")
		linkForPerson, err := terminfindung.CreatePerson(&name, &termin, &user)
		if err != nil {
			return
		}
		type links struct {
			LinkForUser string
			LinkBack    string
		}
		linksFortemp := links{
			LinkBack:    "/shared?terminID=" + termin,
			LinkForUser: "https://" + request.Host + "/shared/public?" + linkForPerson,
		}
		postRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-showlink.html", "./assets/templates/footer.html", "./assets/templates/header.html")
		if err != nil {
			log.Fatal("Coudnt export Parsefiles")
		}
		err = postRoute.Execute(writer, linksFortemp)
		if err != nil {
			log.Fatal("Coudnt Execute Parsefiles")
		}
		return
	}
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-create-link.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, termin)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
}

func ServeHTTPSharedAppCreateDate(writer http.ResponseWriter, request *http.Request) {
	isAllowed, user := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	if request.Method == http.MethodPost {
		//if is submited the form it pases the request
		err := request.ParseForm()
		if err != nil {
			return
		}
		// get start date and enddate and parse it to time format
		startDate := request.Form.Get("startdate")
		endDate := request.Form.Get("enddate")
		startDateFormated, err := time.Parse("2006-01-02", startDate)
		enddateFormated, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return
		}
		//create a new proposed date and redirect to the main website
		err = terminfindung.CreateNewProposedDate(startDateFormated, enddateFormated, &user, &termin, false)
		if err != nil {
			return
		}
		http.Redirect(writer, request, "https://"+request.Host+"/shared?terminID="+termin, http.StatusContinue)
		return
	}
	if len(termin) == 0 {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	mainRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-create-app.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, termin)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}

}

func ShowAllLinksServeHttp(writer http.ResponseWriter, request *http.Request) {
	isAllowed, userAdmin := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	termin := request.URL.Query().Get("terminID")
	links, err := terminfindung.GetAllLinks(&userAdmin, &termin)
	if err != nil {
		return
	}
	for key, user := range links {
		links[key].Url = "https://" + request.Host + "/shared/public?terminID=" + url.QueryEscape(termin) + "&name=" + user.Name + "&user=" + userAdmin + "&apiKey=" + user.Url

	}
	//setup struct for html template
	type shared struct {
		Users     []terminfindung.UserTermin
		Routeback string
	}
	forTemplate := shared{
		Users:     links,
		Routeback: termin,
	}
	linkRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-admin-showAll.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = linkRoute.Execute(writer, forTemplate)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
	return
}

func PublicSharedWebsite(writer http.ResponseWriter, request *http.Request) {
	linkRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-public.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = linkRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
	return
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
