package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
)

type FilterView struct {
	FilteredTermins []ds.Termin
	CurrentPage     int
	EntriesPerPage  int
}

// InitFilterView
// Rückgabewert: Pointer auf ein Objekt FilterView
// Dient zur Initialisierung der FilterView zum Start des Programms.
func InitFilterView(terminCache []ds.Termin) *FilterView {
	var fv = new(FilterView)
	copy(fv.FilteredTermins, terminCache)
	fv.SortEntries(fv.FilteredTermins)
	fv.CurrentPage = 1
	fv.EntriesPerPage = 5
	return fv
}

// GetEntries
// liefert die gefilterten Informationen
func (fv FilterView) GetEntries() []ds.Termin {
	return fv.FilteredTermins
}

// SelectEntriesPerPage
// Parameter: int, gewünschte Anzahl Einträge pro Seite
// setzt die Anzahl Einträge pro Seite auf die vom Benutzer gewählte
// Die aktuelle Seite wird wieder auf 1 gesetzt
func (fv *FilterView) SelectEntriesPerPage(amount int) {
	fv.EntriesPerPage = amount
	fv.CurrentPage = 1
}

// SortEntries
// Parameter: Slice mit den Terminen des Nutzers, die den Filteroptionen entsprechen
// Rückgabewert: ein sortiertes Slice, welches die Termine nach ihrem Startdatum sortiert
// Sortieralgorithmus: Bubble-Sort
// QUELLE: https://www.linux-magazin.de/ausgaben/2020/02/snapshot-23/
func (fv FilterView) SortEntries(entries []ds.Termin) {
	sortedEntries := entries

	for i := range sortedEntries {
		for j := i + 1; j < len(sortedEntries); j++ {
			dateOld := sortedEntries[i].Date
			dateNew := sortedEntries[j].Date
			if dateOld.After(dateNew) {
				sortedEntries[i], sortedEntries[j] =
					sortedEntries[j], sortedEntries[i]
			}
		}
	}
}

// RequiredPages
// Berechnet je nachdem wie viele Einträge pro Seite gewünscht sind die benötigte Seizenanzahl und weist diese dem entsprechenden
// Feld in dem Objekt ListView zu.
func (fv FilterView) RequiredPages() int {
	requiredPages := len(fv.FilteredTermins) / fv.EntriesPerPage
	if len(fv.FilteredTermins)%fv.EntriesPerPage != 0 || requiredPages == 0 {
		requiredPages += 1
	}
	return requiredPages
}
