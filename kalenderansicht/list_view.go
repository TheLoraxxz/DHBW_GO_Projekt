package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
	"net/http"
	"time"
)

type ListView struct {
	SelectedDate        time.Time
	EntriesSinceSelDate []ds.Termin
	EntriesPerPage      int
	PagesAmount         int
}

// initTableView
// Rückgabewert: Pointer auf ein Objekt ListView
// Dient zur Initialisierung der ListView zum Start des Programms.
// Zu Begin wird diese auf das aktuelle Datum gesetzt, die Seitenanzahl Terminen wird die Seite mehrseitig.
func InitListView(terminCache []ds.Termin) *ListView {
	var lv = new(ListView)
	lv.SelectedDate = time.Now()
	lv.EntriesPerPage = 5
	lv.PagesAmount = lv.requiredPages()
	lv.CreateTerminListEntries(terminCache)
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
	lv.EntriesPerPage = amount * 5
	fmt.Println(lv.EntriesPerPage)
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

// CreateTerminListEntries
// Parameter: Slice mit allen Terminen des Nutzers
// Die Funktion weist dem Feld EntriesSinceSelDate des Listview Objektes eine Slice mit allen Terminen des Users seit dem
// gewünschten Datum zu.
func (lv *ListView) CreateTerminListEntries(terminCache []ds.Termin) {
	lv.EntriesSinceSelDate = lv.FilterCalendarEntries(terminCache)
}

// FilterCalendarEntries
// Parameter: Slice mit allen Terminen des Nutzers
// Rückgabewert: Ein Slice mit allen Terminen des Users seit dem gewünschten Datum
func (lv ListView) FilterCalendarEntries(termins []ds.Termin) []ds.Termin {
	startDate := lv.SelectedDate

	entriesSinnceSelDate := make([]ds.Termin, 0, len(termins))
	for _, termin := range termins {
		if termin.EndDate.After(startDate) || termin.EndDate.Equal(startDate) {
			entriesSinnceSelDate = append(entriesSinnceSelDate, termin)
		}
	}
	//Hier werden Termine zum testen hinzugeügt: LÖSCHEN SPÄTER
	for i := 0; i < 200; i++ {
		startDate = startDate.AddDate(0, 0, 1)
		//Delete: is just for testing
		//entriesSinnceSelDate = append(entriesSinnceSelDate, ds.Termin{Title: "Test1", Description: "boa", Recurring: ds.Repeat((i % 5)), Date: startDate, EndDate: startDate.AddDate(3, 4, 2)})
	}
	return entriesSinnceSelDate
}

// requiredPages
// Berechnet je nachdem wie viele Einträge pro Seite gewünscht sind die benötigte Seizenanzahl und weist diese dem entsprechenden
// Feld in dem Objekt ListView zu.
func (lv ListView) requiredPages() int {
	return len(lv.EntriesSinceSelDate) / lv.EntriesPerPage
}

// NextOccurrences
// Parameter: ein Termin
// Rückgabewert: drei Instanzen des Typs time.Time
// Berechnet je nach Wiederholung des Termins und des gewählten Datums, die nächsten drei Vorkommen des Termins.
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
