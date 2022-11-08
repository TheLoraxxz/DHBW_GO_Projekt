package kalenderansicht

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type kalenderEinträgeGenerator interface {
	ErstelleKalenderEintraege(userId string)
}

// Hier werden die Objekte für die Tabellen sowie die Listenasnicht erstellt, die für das Template nötig sind
var ta = TabellenAnsicht{
	Datumsanzeige: time.Now(),
}
var la = new(ListenAnsicht)

// Templates für die Listenansicht sowie die Tabellensansicht
var path, _ = os.Getwd()
var tabellenTpl, _ = template.New("tbl.html").ParseFiles(path+"/assets/templates/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html")
var listenTpl, _ = template.New("listenAnsicht.html").ParseFiles(path+"/assets/templates/lise.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html")

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
	if r.Method == "POST" {
		fmt.Println(r.FormValue("titel"))
		w.WriteHeader(404)
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
