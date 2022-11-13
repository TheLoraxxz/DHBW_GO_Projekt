package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
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
// Parameter: spezifisches Datum
// setzt das Datum der Listenansicht auf das vom Benutzer gewählte
func (lv *ListView) SelectDate(date time.Time) {
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
		entriesForThisMonth = append(entriesForThisMonth, ds.Termin{Title: "Test1", Description: "boa", Recurring: ds.Repeat((i % 5)), Date: startDate})
	}
	return entriesForThisMonth
}

func (lv ListView) requiredPages() int {
	return len(lv.EntriesSinceSelDate) / lv.EntriesPerPage
}
