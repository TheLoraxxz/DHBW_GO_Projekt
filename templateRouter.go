package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	ka "DHBW_GO_Projekt/kalenderansicht"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (h RootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			return
		}
		username := request.Form.Get("user")
		password := request.Form.Get("password")
		// cookie authentifizieren checken
		isUser, cookieText := authentifizierung.AuthenticateUser(&username, &password)
		if isUser == true {
			// wenn user authentifiziert ist dann wird cookie erstellt und
			cookie := &http.Cookie{
				Name:     "SessionID-Kalender",
				Value:    cookieText,
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(writer, cookie)
			//redirect to new site
			http.Redirect(writer, request, "https://"+request.Host+"/user/Create", http.StatusFound)
			return
		} else {
			// wenn nicht authentifiziert ist wird weiter geleitet oder bei problemen gibt es ein 500 status
			if len(cookieText) == 0 {
				writer.WriteHeader(500)
			} else {
				http.Redirect(writer, request, "/", http.StatusContinue)
			}
		}
	}
	mainRoute, err := template.ParseFiles("./assets/sites/index.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}

func (c CreatUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// get cookie
	cookie, err := request.Cookie("SessionID-Kalender")
	//if cookie is not existing it returns back to the host
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	//if it is not allowed then continue with normal website else redirect to root
	isAllowed, _ := authentifizierung.CheckCookie(&cookie.Value)
	if isAllowed {
		mainRoute, err := template.ParseFiles("./assets/sites/create-User.html", "./assets/templates/footer.html", "./assets/templates/header.html")
		if err != nil {
			log.Fatal("Coudnt export Parsefiles")
			return
		}
		err = mainRoute.Execute(writer, nil)
		if err != nil {
			log.Fatal("Coudnt Execute Parsefiles")
			return
		}
	} else {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
	}

}

// Templates für die Tabellensansicht sowie die Listenansicht
var path, _ = os.Getwd()
var tableTpl, _ = template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")
var listTpl, _ = template.New("liste.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")

// ServeHTTP
// Hier werden all http-Request anfragen geregelt, die im Kontext der Terminasnichten anfallen
// Zunächst wird der Cookie geprüft und ggf. die Termine/Infos des Users geladen
// Nach erfolgreicher Prüfung, wird die Anfrage an entweder den ListViewHandler oder den TableViewHandler weitergeleitet.
func (v *ViewmanagerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID-Kalender")
	//wenn kein cookie gesetzt, zurück zur Startwebseite
	if err != nil {
		http.Redirect(w, r, "https://"+r.Host, http.StatusContinue)
		return
	}
	cookieVal := cookie.Value

	//checkt cookies: wenn false redirect zur Startwebseite
	if cookie.Value != v.cookie {
		isAllowed, username := authentifizierung.CheckCookie(&cookieVal)
		if isAllowed {
			v.vm = ka.InitViewManager(username)
			v.cookie = cookie.Value
		} else {
			http.Redirect(w, r, "https://"+r.Host, http.StatusContinue)
			return
		}
	}
	//leitet Anfrage entsprechend weiter

	switch {
	case strings.Contains(r.RequestURI, "/user/view/table"):
		v.handleTableView(w, r)
	case strings.Contains(r.RequestURI, "/user/view/list"):
		v.handleListView(w, r)
	}
}

// handleTableView
// Hier werden all http-Request anfragen geregelt, die im Kontext der TableView anfallen
func (v *ViewmanagerHandler) handleTableView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch {
		case r.RequestURI == "/user/view/table?suche=minusMonat":
			v.vm.TvJumpMonthFor()
		case r.RequestURI == "/user/view/table?suche=plusMonat":
			v.vm.TvJumpMonthBack()
		case strings.Contains(r.RequestURI, "/user/view/table?monat="):
			monatStr := r.RequestURI[23:]
			monat, _ := strconv.Atoi(monatStr)
			v.vm.TvSelectMonth(time.Month(monat))
		case strings.Contains(r.RequestURI, "/user/view/table?jahr="):
			summandStr := r.RequestURI[22:]
			summand, _ := strconv.Atoi(summandStr)
			v.vm.TvJumpYearForOrBack(summand)
		case r.RequestURI == "/user/view/table?datum=heute":
			v.vm.TvJumpToToday()
		}
	}

	if r.Method == "POST" {
		switch {
		case r.RequestURI == "/user/view/table?terminErstellen":
			v.vm.CreateTermin(r, v.vm.Username)
		case r.RequestURI == "/user/view/table?termineBearbeiten":
			v.vm.EditTermin(r, v.vm.Username)
		}
	}
	er := tableTpl.ExecuteTemplate(w, "tbl.html", v.vm.Tv)
	if er != nil {
		log.Fatalln(er)
	}
}

// ListHandler
// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func (v *ViewmanagerHandler) handleListView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch {
		case strings.Contains(r.RequestURI, "/user/view/list?Eintraege="):
			amountStr := r.RequestURI[26:]
			amount, _ := strconv.Atoi(amountStr)
			v.vm.LvSelectEntriesPerPage(amount)
		case r.RequestURI == "/user/view/list?Seite=+1":
			v.vm.LvJumpPageForward()
		case r.RequestURI == "/user/view/list?Seite=-1":
			v.vm.LvJumpPageBack()
		}
	}

	if r.Method == "POST" {
		switch {
		case r.RequestURI == "/user/view/list?selDatum":
			v.vm.LvSelectDate(r)
		case r.RequestURI == "/user/view/list?termineBearbeiten":
			v.vm.EditTermin(r, v.vm.Username)
		case r.RequestURI == "/user/view/list?terminErstellen":
			v.vm.CreateTermin(r, v.vm.Username)
		}
	}

	er := listTpl.ExecuteTemplate(w, "liste.html", v.vm.Lv)
	if er != nil {
		log.Fatalln(er)
	}
}
