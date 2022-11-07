package kalenderansicht

import (
	"html/template"
	"net/http"
	"time"
)

type kalenderEinträgeGenerator interface {
	ErstelleKalenderEintraege(userId string)
}

// Hier werden die Objekte für die Tabellen sowie die Listenasnicht erstellt, die für das Template nötig sind
var tk = TabellenAnsicht{
	Datumsanzeige: time.Now(),
}
var ck = new(ListenAnsicht)

// Templates für die Listenasnischt sowie die Tabellensansicht
var tabellenTpl, _ = template.New("tabellenAnsicht.html").ParseFiles("./sources/tabellenAnsicht.html")
var listenTpl, _ = template.New("listenAnsicht.html").ParseFiles("./sources/tabellenAnsicht.html")

// Hier wereden all http-Request anfragen geregelt,die im Kontext der Kalenderansicht anfallen
func tabellenHandler(w http.ResponseWriter, r *http.Request) {
	//r.RequestURI == /monat -> einträge erstellen
	switch r.RequestURI {
	case "/tabellenAnsicht?suche=minusMonat":
		tk.SpringMonatZurueck()
	case "/tabellenAnsicht?suche=plusMonat":
		tk.SpringMonatVor()
	}
	tabellenTpl.Execute(w, tk)
}

// Hier werden all http-Request-Anfragen geregelt, die im Kontext der Listenansicht anfallen
func listenHandler(w http.ResponseWriter, r *http.Request) {
	listenTpl.Execute(w, nil)
}
