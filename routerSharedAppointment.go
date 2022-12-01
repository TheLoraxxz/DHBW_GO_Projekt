package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/terminfindung"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
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
		fmt.Println("Coudnt find Termin")
		return
	}
	if len(selectedDay) != 0 {
		err := terminfindung.SelectDate(&selectedDay, &termin, &user)
		if err != nil {
			http.Redirect(writer, request, "https://"+request.Host+"/error?type=shared_admin_WrongSelected&link="+url.QueryEscape("/shared?terminID="+termin), http.StatusContinue)
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
	//because go automatically returns it as unescaped query we need to redo it
	var apikey string
	if request.Method == http.MethodPost {
		err := request.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		apikey = url.QueryEscape(request.Form.Get("apiKey"))
		dateKey := request.Form.Get("dateKey")
		voted := request.Form.Get("voted")
		votedBool := false
		if strings.Compare(voted, "on") == 0 {
			votedBool = true
		}
		termin, user, err := terminfindung.GetTerminViaApiKey(&apikey)
		if err != nil {
			return
		}
		err = terminfindung.VoteForDay(&termin.Info.ID, &termin.User, &user, &dateKey, votedBool)
		if err != nil {
			return
		}
	} else {
		apikey = url.QueryEscape(request.URL.Query().Get("apiKey"))
	}
	termin, user, err := terminfindung.GetTerminViaApiKey(&apikey)
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	htmlInput := termin.ConvertUserSiteToRightHTML(&user, &apikey)
	linkRoute, err := template.ParseFiles("./assets/sites/terminfindung/termin-public.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = linkRoute.Execute(writer, htmlInput)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
	return
}

func ErrorSite_ServeHttp(writer http.ResponseWriter, request *http.Request) {
	type errorConfig struct {
		Text string
		Link string
	}
	var config errorConfig
	errorRoute, err := template.ParseFiles("./assets/sites/error.html", "./assets/templates/footer.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	typeErr := request.URL.Query().Get("type")
	link := request.URL.Query().Get("link")
	if val, ok := errorconfigs[typeErr]; ok {
		config = errorConfig{
			Text: val,
			Link: "https://" + request.Host + link,
		}
	} else {
		config = errorConfig{
			Text: errorconfigs["emptyError"],
			Link: "https://" + request.Host,
		}
	}
	errorRoute.Execute(writer, config)
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
