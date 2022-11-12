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

// FilterCalendarEntries
// Parameter: Slice, die alle Termine des Users enthält
// Rückgabewert: Gefilterte Slice, deren Länge der Monatslänge entspricht und für jeden Tag die Termininfos enthält
// Der Termintag entspricht dabei dem Index -1 in dem Slice
func (c TableView) FilterCalendarEntries(termins []ds.Termin) []dayInfos {

	monthStartDate := c.ShownDate
	monthEndDate := c.getLastDayOfMonth()
	//Die Termine für diesen Monat werden in ein Slice gefiltert
	//für jeden Tag des Monats befindet sich ein Objekt DayInfos, welches alle termine enthält in dem Slice
	//Der Index des Tages ist in diesem Falle die Tagesnummer im Monat -1
	//Der 1.01.2022 wäre dementsprechend beim Index 0
	entriesForThisMonth := make([]dayInfos, getMaxDays(int(c.ShownDate.Month()), c.ShownDate.Year()))
	for _, termin := range termins {
		if (termin.Date.Before(monthEndDate) || termin.Date.Equal(monthEndDate)) && (termin.EndDate.After(monthStartDate) || termin.EndDate.Equal(monthStartDate)) {
			switch termin.Recurring {
			case ds.Never, ds.YEARLY, ds.MONTHLY:
				monthDay := termin.Date.Day()
				entriesForThisMonth[monthDay-1].Dayentries = append(entriesForThisMonth[monthDay-1].Dayentries, termin)
				// Vom Start des Termins wird je ein Woche dazu addiert und geprüft, ob diese in den betrachteten Monat fallen
				// Fallen sie in den Zeitraum werden diese der Slice Liste hinzugefügt
			case ds.WEEKLY:
				startDateOfTermin := termin.Date
				folgeTermin := startDateOfTermin
				for folgeTermin.Before(termin.EndDate) {
					if (folgeTermin.Before(monthEndDate) || folgeTermin.Equal(monthEndDate)) && (folgeTermin.After(monthStartDate) || folgeTermin.Equal(monthStartDate)) {
						monthDay := folgeTermin.Day()
						entriesForThisMonth[monthDay-1].Dayentries = append(entriesForThisMonth[monthDay-1].Dayentries, termin)
					}
					folgeTermin = folgeTermin.AddDate(0, 0, 7)
				}
			}
		}
	}
	//Hier werden für jeden Tag die restlichen Informationen sowie Funktionen hinzugefügt
	for i := 0; i < len(entriesForThisMonth); i++ {
		entriesForThisMonth[i].Day = monthStartDate
		monthStartDate = monthStartDate.AddDate(0, 0, 1)
		entriesForThisMonth[i].NeedsBreak = NeedsBreak
		entriesForThisMonth[i].IsToday = IsToday
		//Delete: is just for testing
		entriesForThisMonth[i].Dayentries = append(entriesForThisMonth[i].Dayentries, ds.Termin{Title: "Test1", Description: "boa", Date: monthStartDate})
		entriesForThisMonth[i].Dayentries = append(entriesForThisMonth[i].Dayentries, ds.Termin{Title: "Test2", Description: "boa", Date: monthStartDate})

	}
	return entriesForThisMonth
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

// Die Variable wird gesetzt, wenn der Benutzer auf "bearbeiten" auf der Webseite klickt
// Die Variable enthält alle Termine, die im moment bearbeitet werden
var terminsForEditing []ds.Termin

// SetEditableTermins
// Parameter: ein Objekt DayInfos, das die zu bearbeitenden Termine enthält
// Die Funktion setzt die Termine, die momentan auf der Webseite bearbeitet werden
func SetEditableTermins(editableTermins dayInfos) {
	terminsForEditing = editableTermins.Dayentries
}

// EditTermin
// Parameter:
// Die Funktion ...
func EditTermin(r *http.Request, username string) {
	//Filtern der Termininfos
	for i := 1; i <= len(terminsForEditing); i++ {
		index := strconv.FormatInt((int64(i)), 10)
		mode := r.FormValue("editing" + index)
		//editedTerminIndex := mode[19:]
		if strings.Contains(mode, "Bearbeiten") {
			title := r.FormValue("title" + index)
			description := r.FormValue("description" + index)
			repStr := r.FormValue("repeat" + index)

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
			date, _ := time.Parse(layout, r.FormValue("date"+index))
			endDate, _ := time.Parse(layout, r.FormValue("endDate"+index))

			//EditetTermin kann gelöscht werden, dient nur der Kontrolle
			editetTermin := ds.Termin{
				title,
				description,
				rep,
				date,
				endDate,
			}
			fmt.Println(editetTermin)
			//editedTerminIndex, _ := strconv.Atoi(editedTerminIndex)
			ds.CreateNewTermin(title, description, rep, date, endDate, username)
		}
		//Hier implemntieren, dass gelöscht werden muss
		//deleteTermin(title, description, rep, date, endDate, username)
	}
}
