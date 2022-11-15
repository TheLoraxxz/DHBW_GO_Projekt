package main

import (
	ka "DHBW_GO_Projekt/kalenderansicht"
	"fmt"
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
		fmt.Println(request.FormValue("user"))
		writer.WriteHeader(404)
	}
	mainRoute, err := template.ParseFiles("./assets/sites/index.html", "./assets/templates/footer.html")
	err = mainRoute.Execute(writer, nil)
	if err != nil {
		log.Fatal("Coudnt export Parsefiles")
	}

}

// Das Objekt ViewManager wird benötigt, um den Username, die ListView, die TableView sowie den Cache mit den Terminen des Nutzers zu verwalten
var vm ka.ViewManager

// Templates für die Tabellensansicht sowie die Listenansicht
var path, _ = os.Getwd()
var tableTpl, _ = template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")
var listTpl, _ = template.New("liste.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")

// TableHandler
// Hier werden all http-Request anfragen geregelt,die im Kontext der Kalenderansicht anfallen
// Ich muss hier noch iwi den Usernamer herausfiltern können
func TableHandler(w http.ResponseWriter, r *http.Request) {
	//UserId muss noch iwo geholt werden
	vm.InitViewManager("mik")
	if r.Method == "GET" {
		switch {
		case r.RequestURI == "/tabellenAnsicht?suche=minusMonat":
			vm.Tv.JumpMonthFor()
			vm.Tv.CreateTerminTableEntries(vm.TerminCache)
		case r.RequestURI == "/tabellenAnsicht?suche=plusMonat":
			vm.Tv.JumpMonthBack()
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?monat="):
			monatStr := r.RequestURI[23:]
			monat, _ := strconv.Atoi(monatStr)
			vm.Tv.SelectMonth(time.Month(monat))
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?jahr="):
			summandStr := r.RequestURI[22:]
			summand, _ := strconv.Atoi(summandStr)
			vm.Tv.JumpToYear(summand)
		case r.RequestURI == "/tabellenAnsicht?datum=heute":
			vm.Tv.JumpToToday()
		}
	}

	if r.Method == "POST" {
		switch {
		case r.RequestURI == "/tabellenAnsicht?terminErstellen":
			vm.CreateTermin(r, vm.Username)
		case r.RequestURI == "/tabellenAnsicht?termineBearbeiten":
			vm.EditTermin(r, vm.Username)
		}
	}
	er := tableTpl.ExecuteTemplate(w, "tbl.html", vm.Tv)
	if er != nil {
		log.Fatalln(er)
	}
}

// ListHandler
// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func ListHandler(w http.ResponseWriter, r *http.Request) {
	vm.InitViewManager("mik")
	if r.Method == "GET" {
		switch {
		case strings.Contains(r.RequestURI, "/listenAnsicht?Eintraege="):
			amountStr := r.RequestURI[25:]
			amount, _ := strconv.Atoi(amountStr)
			vm.Lv.SelectEntriesPerPage(amount)
		}
	}

	if r.Method == "POST" {
		switch {
		case r.RequestURI == "/listenAnsicht?selDatum":
			vm.Lv.SelectDate(r)
		case r.RequestURI == "/listenAnsicht?termineBearbeiten":
			vm.EditTermin(r, vm.Username)
		}
	}

	er := listTpl.ExecuteTemplate(w, "liste.html", vm.Lv)
	if er != nil {
		log.Fatalln(er)
	}
}
