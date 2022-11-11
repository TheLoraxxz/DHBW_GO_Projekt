package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"net/http"
	"time"
)

type CalendarEntriesGenerator interface {
	FilterCalendarEntries(username string) string
}

// getMaxDays
// Parameter: Monat und Jahr eines Datums
// Rückgabewert: Anzahl der Tage des Monats
// -> Schaltjahre werden berücksichtigt
func getMaxDays(month, year int) int {
	var days int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		days = 31
	case 4, 6, 9, 11:
		days = 30
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

// CreateTermin
// Parameter: eine Post-Request mit Informationen über einen Termin und den Usernamen des Nutzers, der diesen anlegen möchte
// CreateTermin ruft die Funktion zum Erstellen des Termins auf
func CreateTermin(r *http.Request, username string) ds.Termin {

	//Filtern der Termininfos
	title := r.FormValue("title")
	description := r.FormValue("description")
	repStr := r.FormValue("repeat")

	//Filter das Wiederholungsintervall aus der Antwort
	var rep ds.Repeat
	switch repStr {
	case "täglich":
		rep = 0
	case "wöchentlich":
		rep = 1
	case "jährlich":
		rep = 2
	case "niemals":
		rep = 3
	}

	//Daten in das richtige Format überführen mithilfe eines Layouts
	layout := "2006-01-02"
	date, _ := time.Parse(layout, r.FormValue("date"))
	endDate, _ := time.Parse(layout, r.FormValue("endDate"))

	//End Date Logik-check
	if repStr == "niemals" || endDate.Before(date) {
		endDate = date
	}

	//erstelle neuen Termin
	return ds.CreateNewTermin(title, description, rep, date, endDate, username)
}
