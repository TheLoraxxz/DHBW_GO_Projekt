package main

import (
	"DHBW_GO_Projekt/kalenderansicht"
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

var username = "Testuser" //Nur zum testen, schauen woher ich diesen bekomme
// Diese Variablen werden benötigt um die objektbezogenen Funktionen im Handler aufzurufen

var ta = kalenderansicht.TabellenAnsicht{
	Datumsanzeige: time.Now(),
}
var la = new(kalenderansicht.ListenAnsicht)

// Templates für die Listenansicht sowie die Tabellensansicht
var path, _ = os.Getwd()
var tabellenTpl, _ = template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")
var listenTpl, _ = template.New("listenAnsicht.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html")

// Hier wereden all http-Request anfragen geregelt,die im Kontext der Kalenderansicht anfallen
func TabellenHandler(w http.ResponseWriter, r *http.Request) {
	//r.RequestURI == /monat -> einträge erstellen
	if r.Method == "GET" {
		switch {
		case r.RequestURI == "/tabellenAnsicht?suche=minusMonat":
			ta.SpringMonatZurueck()
		case r.RequestURI == "/tabellenAnsicht?suche=plusMonat":
			ta.SpringMonatVor()
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?monat="):
			monatStr := r.RequestURI[23:]
			monat, _ := strconv.Atoi(monatStr)
			ta.WaehleMonat(time.Month(monat))
		case strings.Contains(r.RequestURI, "/tabellenAnsicht?jahr="):
			summandStr := r.RequestURI[22:]
			summand, _ := strconv.Atoi(summandStr)
			ta.SpringeJahr(summand)
		case r.RequestURI == "/tabellenAnsicht?datum=heute":
			ta.SpringZuHeute()
		}
	}

	if r.Method == "POST" && r.RequestURI == "/tabellenAnsicht?terminErstellen" {
		kalenderansicht.CreateTermin(r, username)
	}
	er := tabellenTpl.ExecuteTemplate(w, "tbl.html", ta)
	if er != nil {
		log.Fatalln(er)
	}
}

// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func ListenHandler(w http.ResponseWriter, r *http.Request) {
	listenTpl.Execute(w, nil)
}
