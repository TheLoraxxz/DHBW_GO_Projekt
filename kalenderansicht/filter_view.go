package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"time"
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
	fv.FilteredTermins = make([]ds.Termin, len(terminCache))
	copy(fv.FilteredTermins, terminCache)
	fv.SortEntries(fv.FilteredTermins)
	fv.CurrentPage = 1
	fv.EntriesPerPage = 5
	return fv
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die den Benutzer Custom-Settings & Navigation innerhalb der Webseite ermöglichen.
(Bsp.: Seitenanzahl festlegen, Seite weiter navigieren...)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// SelectEntriesPerPage
// Parameter: int, gewünschte Anzahl Einträge pro Seite
// setzt die Anzahl Einträge pro Seite auf die vom Benutzer gewählte
// Die aktuelle Seite wird wieder auf 1 gesetzt
func (fv *FilterView) SelectEntriesPerPage(amount int) {
	fv.EntriesPerPage = amount
	fv.CurrentPage = 1
}

// JumpPageForward
// springt eine Seite in der Webseite weiter
func (fv *FilterView) JumpPageForward() {
	if fv.CurrentPage+1 <= fv.RequiredPages() {
		fv.CurrentPage += 1
	}
}

// JumpPageBack
// springt eine Seite in der Webseite zurück
func (fv *FilterView) JumpPageBack() {
	if fv.CurrentPage-1 > 0 {
		fv.CurrentPage -= 1
	}
}

// GetEntries
// Rückgabewert: Ein Slice mit den der entsprechenden Anzahl an Terminen, die auf der aktuellen Seite angezeigt werden
// Funktion wird im template aufgerufen, um Termine anzuzeigen
func (fv FilterView) GetEntries() []ds.Termin {
	entries := make([]ds.Termin, 0, fv.EntriesPerPage)
	sliceStart := fv.EntriesPerPage * (fv.CurrentPage - 1)

	for entryNr := 0; entryNr < fv.EntriesPerPage; entryNr++ {
		if sliceStart+entryNr < len(fv.FilteredTermins) {
			entries = append(entries, fv.FilteredTermins[sliceStart+entryNr])
		}
	}
	return entries
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Filtern und Sortieren der Termine in der Filteransicht dienen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// RequiredPages
// Berechnet je nachdem wie viele Einträge pro Seite gewünscht sind die benötigte Seitenanzahl und weist diese dem entsprechenden
// Feld in dem Objekt ListView zu.
func (fv FilterView) RequiredPages() int {
	requiredPages := len(fv.FilteredTermins) / fv.EntriesPerPage
	if len(fv.FilteredTermins)%fv.EntriesPerPage != 0 || requiredPages == 0 {
		requiredPages += 1
	}
	return requiredPages
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

// NextOccurrences
// Parameter: ein Termin
// Rückgabewert: drei Instanzen des Typs time.Time
// berechnet die nächsten drei Vorkommen des Termins ab heute
func (fv FilterView) NextOccurrences(termin ds.Termin) []time.Time {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	nextOccurrences := make([]time.Time, 0, 3)

	occur := termin.Date
	noMoreOccur := false
	//solange nicht die drei nächsten Termine gefiltert worden sind
	//und das letzte Vorkommen des Termins noch nicht erreicht worden ist, füge weitere Termine der Liste hinzu
	//Wenn der Termin nur einmal vorkommt, sorgt die Variable noMoreOccur für einen Abbruch,
	//so wird nicht 3 Mal derselbe Termin hinzugefügt.
	for len(nextOccurrences) < 3 && (!occur.After(termin.EndDate)) && noMoreOccur == false {
		if occur.After(today) || occur.Equal(today) {
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

// CreateTerminFilterEntries
// Parameter: Slice mit allen Terminen des Nutzers.
// Die Funktion überreicht der FilterAnsicht die aktuellste Liste der Termine.
func (fv *FilterView) CreateTerminFilterEntries(terminCache []ds.Termin) {
	fv.FilteredTermins = make([]ds.Termin, len(terminCache))
	copy(fv.FilteredTermins, terminCache)
	fv.SortEntries(fv.FilteredTermins)
}

func (fv *FilterView) FilterTermins(filterTitle, filterDescription string, allTermins []ds.Termin) {
	fv.FilteredTermins = allTermins
	if filterTitle != "" {
		fv.FilteredTermins = ds.FilterByTitle(fv.FilteredTermins, filterTitle)
	}
	if filterDescription != "" {
		fv.FilteredTermins = ds.FilterByDescription(fv.FilteredTermins, filterDescription)
	}
	fv.SortEntries(fv.FilteredTermins)
}
