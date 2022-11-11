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

// Die folgenden Variablen werden benötigt, um die objektbezogenen Funktionen im Handler aufzurufen
var tv = ka.InitTableView()
var la ka.ListView

// Templates für die Listenansicht sowie die Tabellensansicht
var path, _ = os.Getwd()
var tableTpl, _ = template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")
var listTpl, _ = template.New("liste.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")

// TabellenHandler
// Hier werden all http-Request anfragen geregelt,die im Kontext der Kalenderansicht anfallen
// Ich muss hier noch iwi den Usernamer herausfiltern können
func TabellenHandler(w http.ResponseWriter, r *http.Request) {
	//UserId muss noch iwo geholt werden
	tv.Username = "mik"
	if r.Method == "GET" {
		switch {
		case r.RequestURI == "/tabellenAnsicht?suche=minusMonat":
			tv.JumpMonthFor()
		case r.RequestURI == "/tabellenAnsicht?suche=plusMonat":
			tv.JumpMonthBack()
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?monat="):
			monatStr := r.RequestURI[23:]
			monat, _ := strconv.Atoi(monatStr)
			tv.SelectMonth(time.Month(monat))
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?jahr="):
			summandStr := r.RequestURI[22:]
			summand, _ := strconv.Atoi(summandStr)
			tv.JumpToYear(summand)
		case r.RequestURI == "/tabellenAnsicht?datum=heute":
			tv.JumpToToday()
		}
	}

	if r.Method == "POST" && r.RequestURI == "/tabellenAnsicht?terminErstellen" {
		ka.CreateTermin(r, "mik")
	}
	er := tableTpl.ExecuteTemplate(w, "tbl.html", tv)
	if er != nil {
		log.Fatalln(er)
	}
}

// ListenHandler
// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func ListenHandler(w http.ResponseWriter, r *http.Request) {
	er := listTpl.ExecuteTemplate(w, "liste.html", nil)
	if er != nil {
		log.Fatalln(er)
	}
}
