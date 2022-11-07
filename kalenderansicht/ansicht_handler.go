package kalenderansicht

import (
	"html/template"
	"net/http"
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
var tabellenTpl, _ = template.New("tbl.html").ParseFiles("../assets/templates/tbl.html", "../assets/templates/header.html", "../assets/templates/footer.html")
var listenTpl, _ = template.New("listenAnsicht.html").ParseFiles("../assets/templates/lise.html", "../assets/templates/header.html", "../assets/templates/footer.html")

// Hier wereden all http-Request anfragen geregelt,die im Kontext der Kalenderansicht anfallen
func tabellenHandler(w http.ResponseWriter, r *http.Request) {
	//r.RequestURI == /monat -> einträge erstellen
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
	tabellenTpl.ExecuteTemplate(w, "tbl.html", ta)
}

// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func listenHandler(w http.ResponseWriter, r *http.Request) {
	listenTpl.Execute(w, nil)
}
