package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// getMaxDays
// Parameter: Monat und Jahr eines Datums
// Rückgabewert: Anzahl der Tage des Monats
// Funktion, die öfters zur hilfe aufgerufen wird
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

// filterRepetition
// Parameter: string, eine nummer als string
// Rückgabewert: ein Typ von ds.Repeat, der der entsprechenden Nummer entspricht
func filterRepetition(repStr string) ds.Repeat {
	var rep ds.Repeat
	switch repStr {
	case "1":
		rep = ds.Never
	case "2":
		rep = ds.DAILY
	case "3":
		rep = ds.WEEKLY
	case "4":
		rep = ds.MONTHLY
	case "5":
		rep = ds.YEARLY
	}
	return rep
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
	rep := filterRepetition(repStr)

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

// EditTermin
// Parameter: Post-Request mit den Werten die geändert/gelöscht werden sollen
// Die Funktion ermittelt aus der Post-request die Termine an dem Tag, der bearbeitet wird.
// Anschließend wird für jeden Termin an dem tag ermittelt, ob dieser bearbeitet oder gelöscht werden soll:
//
//	-> Berbeiten: Dei neuen Termininfos werden ermittelt und ein neuer Termin erstellt. Der alte Termin wird gelöscht.
//	-> Löschen: Der Termin wird gelöscht.
func EditTermin(r *http.Request, username string, monthEntries []dayInfos) {
	//Ermitteln des Index des Tages, dessen Termine bearbeitet werden
	editingDayIndex, _ := strconv.Atoi(r.FormValue("editingIndex"))
	editedDay := monthEntries[editingDayIndex].Dayentries

	//Schleife, die alle bearbeiteten Termine des Tages durchgeht
	for i := 1; i <= len(editedDay); i++ {
		//index ist die variable i als String, benötigt um mach den richtigen Values zu suchen
		index := strconv.FormatInt((int64(i)), 10)

		//Filtern der übergebenen neuen Werte
		mode := r.FormValue("editing" + index)

		//Wenn der Filtermodus nicht beearbeiten iste, müssen die übergebenen Values nicht gelesen werden
		if strings.Contains(mode, "Bearbeiten") {
			title := r.FormValue("title" + index)
			description := r.FormValue("description" + index)
			repStr := r.FormValue("repeat" + index)

			//Filter das Wiederholungsintervall aus der Antwort
			rep := filterRepetition(repStr)

			//Daten in das richtige Format überführen mithilfe eines Layouts
			layout := "2006-01-02"
			date, _ := time.Parse(layout, r.FormValue("date"+index))
			endDate, _ := time.Parse(layout, r.FormValue("endDate"+index))

			//EditetTermin kann gelöscht werden, dient nur der Kontrolle!!!!
			editetTermin := ds.Termin{
				title,
				description,
				rep,
				date,
				endDate,
			}
			//Lösche Print anweisung
			fmt.Println(editetTermin)
			//Ein neuer Termin mit den geänderten Werten wird erstellt
			ds.CreateNewTermin(title, description, rep, date, endDate, username)
		}
		//Der alte Termin wird gelöscht -> muss noch implementiert werden
		//deleteTermin(editedDay[i-1])
	}
}
