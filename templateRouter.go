package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	ka "DHBW_GO_Projekt/kalenderansicht"
	"html/template"
	"log"
	"net/http"
	"net/url"
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
			http.Redirect(writer, request, "https://"+request.Host+"/user/view/table", http.StatusFound)
			return
		} else {
			// wenn nicht authentifiziert ist wird weiter geleitet oder bei problemen gibt es ein 500 status
			if len(cookieText) == 0 {
				writer.WriteHeader(500)
			} else {
				request.Method = "GET"
				urls := "https://" + request.Host + "/error?type=wrongAuthentication&link=" + url.QueryEscape("/")
				http.Redirect(writer, request, urls, http.StatusContinue)
				return
			}
		}
	}
	mainRoute, err := template.ParseFiles("./assets/sites/index.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}

func (createUser CreatUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// get cookie
	isAllowed, _ := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	//if it is post it should process the data
	if request.Method == "POST" {
		// if the parseform isnt correct it should return
		err := request.ParseForm()
		if err != nil {
			http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		}
		//get the user from the form request
		user := request.Form.Get("newUsername")
		password := request.Form.Get("newPassword")
		err = authentifizierung.CreateUser(&user, &password)
		if err != nil {
			http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)

		}
		//if successfull on post it should return back to the user
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)

	}
	mainRoute, err := template.ParseFiles("./assets/sites/user-create.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
		return
	}
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
		return
	}

}

func (changeUser ChangeUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isAllowed, user := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	if request.Method == "POST" {
		//if post request it actually parses the form and trys to change the password and create a new cookie
		err := request.ParseForm()
		if err != nil {
			return
		}
		//change the user to new user
		password := request.Form.Get("oldPassword")
		newPassword := request.Form.Get("newPassword")
		cookies, err := authentifizierung.ChangeUser(&user, &password, &newPassword)
		if err != nil {
			http.Redirect(writer, request, "https://"+request.Host+"/user", http.StatusContinue)
			return
		}
		// set cookie so it automatically updates and it doesnt throw one back to the login site
		cookie := &http.Cookie{
			Name:     "SessionID-Kalender",
			Value:    cookies,
			Path:     "/",
			MaxAge:   3600,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		//set new cookie and redirect
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, "https://"+request.Host+"/user", http.StatusContinue)
		return

	}
	//execute own template from userchange and put in footer and header
	mainRoute, err := template.ParseFiles("./assets/sites/user-change.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
}

func (user UserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isAllowed, username := checkIfIsAllowed(request)
	if !isAllowed {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	mainRoute, err := template.ParseFiles("./assets/sites/user.html", "./assets/templates/footer.html", "./assets/templates/header.html")
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return
	}
	err = mainRoute.Execute(writer, username)
	if err != nil {
		log.Fatal("Coudnt Execute Parsefiles")
	}
}

// ServeHTTP
// Hier werden all http-Request anfragen geregelt, die im Kontext der Terminansichten anfallen.
// Zunächst wird der Cookie geprüft und ggf. die Termine/Infos des Users geladen.
// Nach erfolgreicher Prüfung, wird die Anfrage an entweder den ListViewHandler oder den TableViewHandler weitergeleitet.
func (v *ViewManagerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//cookie-Check
	isAllowed, username := checkIfIsAllowed(r)

	//Falls kein Berechtigter-User: Errormeldung + Redirect
	if !isAllowed {
		urls := "https://" + r.Host + "/error?type=wrongAuthentication&link=" + url.QueryEscape("/")
		http.Redirect(w, r, urls, http.StatusContinue)
		return
	}

	//Falls vm noch nicht initialisiert
	if v.vm == nil {
		v.vm = ka.InitViewManager(username)
		v.vm.Username = username
	}

	//Termin bearbeiten/erstellen/löschen ist überall identisch
	edit := r.FormValue("edit")
	create := r.FormValue("create")
	deleteShared := r.FormValue("deleteSharedTermin")
	switch r.Method {
	case "GET":
		if edit == "true" {
			terminToEdit, err := v.vm.GetTerminInfos(r)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=" + err.Error() + "&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
			err = v.viewManagerTpl.ExecuteTemplate(w, "editor.html", terminToEdit)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=internal&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
			return
		}
		if deleteShared != "" {
			terminToDeleteID := r.FormValue("deleteSharedTermin")
			err := v.vm.DeleteSharedTermin(terminToDeleteID, v.vm.Username)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=" + err.Error() + "&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
		}
	case "POST":
		if edit == "true" {
			err := v.vm.EditTermin(r, v.vm.Username)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=" + err.Error() + "&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
		}
		if create == "true" {
			err := v.vm.CreateTermin(r, v.vm.Username)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=" + err.Error() + "&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
		}
	}

	// Anfrage entsprechend weiterleiten (Listen- Tabellen- oder Filteransicht)
	switch {
	case strings.Contains(r.URL.String(), "/user/view/table"):
		v.handleTableView(w, r)
	case strings.Contains(r.URL.String(), "/user/view/list"):
		v.handleListView(w, r)
	case strings.Contains(r.URL.String(), "/user/view/filterTermins"):
		v.handleFilterView(w, r)
	}

}

// handleTableView
// Hier werden all http-Request anfragen geregelt, die im Kontext der TableView anfallen
func (v ViewManagerHandler) handleTableView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch {
		case r.URL.String() == "/user/view/table?suche=minusMonat":
			v.vm.TvJumpMonthBack()
		case r.URL.String() == "/user/view/table?suche=plusMonat":
			v.vm.TvJumpMonthFor()
		case strings.Contains(r.URL.String(), "/user/view/table?monat="):
			monatStr := r.FormValue("monat")
			monat, err := strconv.Atoi(monatStr)
			if err != nil || monat < 1 || monat > 12 {
				urls := "https://" + r.Host + "/error?type=NowValidMonth&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
			v.vm.TvSelectMonth(time.Month(monat))
		case r.URL.String() == "/user/view/table?jahr=Zurueck":
			v.vm.TvJumpYearForOrBack(-1)
		case r.URL.String() == "/user/view/table?jahr=Vor":
			v.vm.TvJumpYearForOrBack(1)
		case r.URL.String() == "/user/view/table?datum=heute":
			v.vm.TvJumpToToday()
		}
	}

	err := v.viewManagerTpl.ExecuteTemplate(w, "tbl.html", v.vm.Tv)
	if err != nil {
		urls := "https://" + r.Host + "/error?type=internal&link=" + r.URL.Path + url.QueryEscape("/")
		http.Redirect(w, r, urls, http.StatusContinue)
		return
	}
}

// ListHandler
// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func (v *ViewManagerHandler) handleListView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch {
		case strings.Contains(r.URL.String(), "/user/view/list?selDate="):
			dateStr := r.FormValue("selDate")
			err := v.vm.LvSelectDate(dateStr)
			if err != nil {
				urls := "https://" + r.Host + "/error?type=" + err.Error() + "&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
		case strings.Contains(r.URL.String(), "/user/view/list?Eintraege"):
			amountStr := r.FormValue("Eintraege")
			amount, err := strconv.Atoi(amountStr)
			if err != nil || !(amount == 5 || amount == 10 || amount == 15) {
				urls := "https://" + r.Host + "/error?type=Unvalid_Entries_Per_Page&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
			v.vm.LvSelectEntriesPerPage(amount)
		case r.URL.String() == "/user/view/list?Seite=Vor":
			v.vm.LvJumpPageForward()
		case r.URL.String() == "/user/view/list?Seite=Zurueck":
			v.vm.LvJumpPageBack()
		}
	}

	err := v.viewManagerTpl.ExecuteTemplate(w, "liste.html", v.vm.Lv)
	if err != nil {
		urls := "https://" + r.Host + "/error?type=internal&link=" + r.URL.Path + url.QueryEscape("/")
		http.Redirect(w, r, urls, http.StatusContinue)
		return
	}
}

// filterTerminsHandler
// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func (v *ViewManagerHandler) handleFilterView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch {
		case strings.Contains(r.URL.String(), "/user/view/filterTermins?Eintraege"):
			amountStr := r.FormValue("Eintraege")
			amount, err := strconv.Atoi(amountStr)
			if err != nil || !(amount == 5 || amount == 10 || amount == 15) {
				urls := "https://" + r.Host + "/error?type=Unvalid_Entries_Per_Page&link=" + r.URL.Path + url.QueryEscape("/")
				http.Redirect(w, r, urls, http.StatusContinue)
				return
			}
			v.vm.FvSelectEntriesPerPage(amount)
		case r.URL.String() == "/user/view/filterTermins?Seite=Vor":
			v.vm.FvJumpPageForward()
		case r.URL.String() == "/user/view/filterTermins?Seite=Zurueck":
			v.vm.FvJumpPageBack()
		case strings.Contains(r.URL.String(), "/user/view/filterTermins?title="):
			v.vm.FvFilter(r)
		}
	}

	err := v.viewManagerTpl.ExecuteTemplate(w, "filterTermins.html", v.vm.Fv)
	if err != nil {
		urls := "https://" + r.Host + "/error?type=internal&link=" + r.URL.Path + url.QueryEscape("/")
		http.Redirect(w, r, urls, http.StatusContinue)
		return
	}
}

func (l LogoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("SessionID-Kalender")
	if err != nil {
		http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
		return

	}

	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(writer, cookie)
	http.Redirect(writer, request, "https://"+request.Host, http.StatusContinue)
}
