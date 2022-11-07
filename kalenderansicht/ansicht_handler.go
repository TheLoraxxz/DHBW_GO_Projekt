package kalenderansicht

import (
	"html/template"
	"net/http"
	"time"
)

type calenderEntriesGenerator interface {
	createCalendarEntries(userId string)
}

var tk = TabellenKalendar{
	Datumsanzeige: time.Now(),
}
var ck = new(ListenAnsicht)

var tabellenTpl, _ = template.New("tabellenAnsicht.html").ParseFiles("./sources/tabellenAnsicht.html")
var listenTpl, _ = template.New("tabellenAnsicht.html").ParseFiles("./sources/tabellenAnsicht.html")

func tabellenHandler(w http.ResponseWriter, r *http.Request) {
	//r.RequestURI == /monat -> einträge erstellen
	switch r.RequestURI {
	case "/tabellenAnsicht?suche=minusMonat":
		tk.springMonatZurueck()
	case "/tabellenAnsicht?suche=plusMonat":
		tk.springMonatVor()
	}
	tabellenTpl.Execute(w, tk)
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	listenTpl.Execute(w, nil)
}
