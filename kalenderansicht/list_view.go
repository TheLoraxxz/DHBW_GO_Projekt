package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"math/rand"
	"net/http"
	"time"
)

type ListView struct {
	SelectedDate        time.Time
	Username            string
	EntriesSinceSelDate []ds.Termin
	EntriesPerPage      int
	PagesAmount         int
}

// initTableView
// Rückgabewert: Pointer auf ein Objekt ListView
// Dient zur Initialisierung der ListView zum Start des Programms.
// Zu Begin wird diese auf das aktuelle Datum gesetzt, die Seitenanzahl Terminen wird die Seite mehrseitig.
func InitListView() *ListView {
	var lv = new(ListView)
	lv.SelectedDate = time.Now()
	lv.EntriesPerPage = 5
	lv.PagesAmount = lv.requiredPages()
	return lv
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die den Benutzer Custom-Settings & Navigation innerhalb der Webseite ermöglichen.
(Bsp.: Seitenanzahl festlegen, Seite weiter navigieren...)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// SelectDate
// Parameter: Post Request mit einem spezifischem Datum
// setzt das Datum der Listenansicht auf das vom Benutzer gewählte
func (lv *ListView) SelectDate(r *http.Request) {

	//Datum Filtern und in das richtige Format überführen mithilfe eines Layouts
	layout := "2006-01-02"
	date, _ := time.Parse(layout, r.FormValue("selDate"))
	lv.SelectedDate = date
}

// SelectEntriesPerPage
// Parameter: int, gewünschte Anzahl Einträge pro Seite
// setzt die Anzahl Einträge pro Seite auf die vom Benutzer gewählte
func (lv *ListView) SelectEntriesPerPage(amount int) {
	lv.EntriesPerPage = amount
}

// JumpPageForward
// springt eine Seite in der Webseite weiter
func (lv ListView) JumpPageForward() {
}

// JumpPageBack
// springt eine Seite in der Webseite zurück
func (lv ListView) JumpPageBack() {

}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Filtern und Anzeigen der Termine in der Listenansicht dienen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// CreateTerminList
func (lv *ListView) CreateTerminList() {
	termins := ds.GetTermine(lv.Username)
	lv.EntriesSinceSelDate = lv.FilterCalendarEntries(termins)
}

func (lv ListView) FilterCalendarEntries(termins []ds.Termin) []ds.Termin {
	startDate := lv.SelectedDate

	//Die Termine für diesen Monat werden in ein Slice gefiltert
	//für jeden Tag des Monats befindet sich ein Objekt DayInfos in der Slice
	//Der Index des Tages ist in diesem Falle die Tagesnummer im Monat -1
	//Der 1.01.2022 wäre dementsprechend beim Index 0

	entriesForThisMonth := make([]ds.Termin, 0, len(termins))
	for _, termin := range termins {
		if termin.EndDate.After(startDate) || termin.EndDate.Equal(startDate) {
			switch termin.Recurring {
			case ds.Never, ds.YEARLY, ds.MONTHLY:
				entriesForThisMonth = append(entriesForThisMonth, termin)
				// Vom Start des Termins wird je eine Woche dazu addiert und geprüft, ob dieses neue Datum in den betrachteten Monat fällt
				// Fällt der Termin in den gewählten Zeitraum, wird der termin in die Slice hinzugefügt
			case ds.WEEKLY:
				startDateOfTermin := termin.Date
				folgeTermin := startDateOfTermin
				for folgeTermin.Before(termin.EndDate) {
					if folgeTermin.After(startDate) || folgeTermin.Equal(startDate) {
						entriesForThisMonth = append(entriesForThisMonth, termin)
					}
					folgeTermin = folgeTermin.AddDate(0, 0, 7)
				}
			}
		}
	}
	//Hier werden Termine zum testen hinzugeügt: LÖSCHEN SPÄTER
	for i := 0; i < 200; i++ {
		startDate = startDate.AddDate(0, 0, 1)
		//Delete: is just for testing
		entriesForThisMonth = append(entriesForThisMonth, ds.Termin{Title: "Test1", Description: "boa", Recurring: ds.Repeat((i % 5)), Date: startDate, EndDate: startDate.AddDate(rand.Int(), rand.Int(), rand.Int())})
	}
	return entriesForThisMonth
}

func (lv ListView) requiredPages() int {
	return len(lv.EntriesSinceSelDate) / lv.EntriesPerPage
}
func (lv ListView) NextOccurrences(termin ds.Termin) []time.Time {
	selDate := lv.SelectedDate
	nextOccurrences := make([]time.Time, 0, 3)

	occur := termin.Date
	noMoreOccur := false
	//solange nicht die drei nächsten Termine gefiltert worden sind
	//und das letzte Vorkommen des Termins noch nicht erreicht worden ist, füge weitere Termine der Liste hinzu
	//Wenn der Termin nur einmal vorkommt, sorgt die Variable noMoreOccur für einen Abbruch,
	//so wird nicht 3 Mal derselbe Termin hinzugefügt.
	for len(nextOccurrences) < 3 && (!occur.After(termin.EndDate)) && noMoreOccur == false {
		if occur.After(selDate) || occur.Equal(selDate) {
			nextOccurrences = append(nextOccurrences, occur)
		}
		switch termin.Recurring {
		case ds.YEARLY:
			occur = occur.AddDate(1, 0, 0)
		case ds.MONTHLY:
			occur = occur.AddDate(0, 1, 0)
		case ds.WEEKLY:
			occur = occur.AddDate(0, 0, 7)
		case ds.DAILY:
			occur = occur.AddDate(0, 0, 1)
		case ds.Never:
			noMoreOccur = true
		}
	}
	return nextOccurrences
}
