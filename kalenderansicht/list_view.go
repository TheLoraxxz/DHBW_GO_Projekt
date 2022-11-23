package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"time"
)

type ListView struct {
	SelectedDate        time.Time
	EntriesSinceSelDate []ds.Termin
	EntriesPerPage      int
	CurrentPage         int
}

// InitListView initTableView
// Rückgabewert: Pointer auf ein Objekt ListView
// Dient zur Initialisierung der ListView zum Start des Programms.
// Zu Begin wird diese auf das aktuelle Datum gesetzt, die Seitenanzahl Terminen wird die Seite mehrseitig.
func InitListView(terminCache []ds.Termin) *ListView {
	var lv = new(ListView)
	lv.SelectedDate = time.Now()
	lv.EntriesPerPage = 5
	lv.CurrentPage = 1
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
// Die aktuelle Seite wird wieder auf 1 gesetzt
func (lv *ListView) SelectDate(date time.Time) {
	lv.SelectedDate = date
	lv.CurrentPage = 1
}

// SelectEntriesPerPage
// Parameter: int, gewünschte Anzahl Einträge pro Seite
// setzt die Anzahl Einträge pro Seite auf die vom Benutzer gewählte
// Die aktuelle Seite wird wieder auf 1 gesetzt
func (lv *ListView) SelectEntriesPerPage(amount int) {
	lv.EntriesPerPage = amount
	lv.CurrentPage = 1
}

// JumpPageForward
// springt eine Seite in der Webseite weiter
func (lv *ListView) JumpPageForward() {
	if lv.CurrentPage+1 <= lv.RequiredPages() {
		lv.CurrentPage += 1
	}
}

// JumpPageBack
// springt eine Seite in der Webseite zurück
func (lv *ListView) JumpPageBack() {
	if lv.CurrentPage-1 > 0 {
		lv.CurrentPage -= 1
	}
}

// GetEntries
// Rückgabewert: Ein Slice mit den Terminen, die auf der aktuellen Seite angezeigt werden
func (lv ListView) GetEntries() []ds.Termin {
	entries := make([]ds.Termin, 0, lv.EntriesPerPage)
	sliceStart := lv.EntriesPerPage * (lv.CurrentPage - 1)

	for entryNr := 0; entryNr < lv.EntriesPerPage; entryNr++ {
		if sliceStart+entryNr < len(lv.EntriesSinceSelDate) {
			entries = append(entries, lv.EntriesSinceSelDate[sliceStart+entryNr])
		}
	}
	return entries
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Filtern und Anzeigen der Termine in der Listenansicht dienen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// CreateTerminListEntries
// Parameter: Slice mit allen Terminen des Nutzers
// Die Funktion weist dem Feld EntriesSinceSelDate des Listview Objektes eine Slice mit allen Terminen des Users seit dem
// gewünschten Datum zu.
func (lv *ListView) CreateTerminListEntries(terminCache []ds.Termin) {
	entries := lv.FilterCalendarEntries(terminCache)
	lv.SortEntries(entries)
	lv.EntriesSinceSelDate = entries
}

// FilterCalendarEntries
// Parameter: Slice mit allen Terminen des Nutzers
// Rückgabewert: Ein Slice mit allen Terminen des Users seit dem gewünschten Datum.
// Die Funktion schaut, ob das Enddatum des Termins nach dem auf der Seite gezeigtem Datum liegt.
// Falls es sich bei den Terminwiederholungen um wöchentliche, monatliche oder jährliche Termine handelt,
// wird zusätzlich geprüft, ob der Termin nach dem gezeigten Datum auch nochmal vorkommt, oder ob nur das Enddatum
// weiter hinten liegt. Kommt der Termin nicht nochmal vor, wird er aussortiert.
func (lv *ListView) FilterCalendarEntries(termins []ds.Termin) []ds.Termin {
	startDate := lv.SelectedDate
	entriesSinceSelDate := make([]ds.Termin, 0, len(termins))
	for _, termin := range termins {
		if termin.EndDate.After(startDate) || termin.EndDate.Equal(startDate) {
			if (termin.Recurring == ds.DAILY) || (termin.Recurring == ds.Never) {
				entriesSinceSelDate = append(entriesSinceSelDate, termin)
			} else {
				date := termin.Date
				lastOccuring := date
				for date.Before(termin.EndDate) || date.Equal(termin.EndDate) {
					lastOccuring = date
					switch termin.Recurring {
					case ds.WEEKLY:
						date = date.AddDate(0, 0, 7)
					case ds.MONTHLY:
						date = date.AddDate(0, 1, 0)
					case ds.YEARLY:
						date = date.AddDate(1, 0, 0)
					}
				}
				if lastOccuring.After(startDate) {
					entriesSinceSelDate = append(entriesSinceSelDate, termin)
				}
			}
		}
	}
	return entriesSinceSelDate
}

// SortEntries
// Parameter: Slice mit den Terminen des Nutzers, die ab dem auf der Seite gezeigten Datum vorkommen
// Rückgabewert: ein sortiertes Slice, welches die Termine nach ihrem nächsten Vorkommen zeitlich sortiert
// Sortieralgorithmus: Bubble-Sort
// QUELLE: https://www.linux-magazin.de/ausgaben/2020/02/snapshot-23/
func (lv ListView) SortEntries(entries []ds.Termin) {
	sortedEntries := entries

	for i := range sortedEntries {
		for j := i + 1; j < len(sortedEntries); j++ {
			nextOccurringOld := lv.NextOccurrences(sortedEntries[i])[0]
			nextOccurringNew := lv.NextOccurrences(sortedEntries[j])[0]
			if nextOccurringOld.After(nextOccurringNew) {
				sortedEntries[i], sortedEntries[j] =
					sortedEntries[j], sortedEntries[i]
			}
		}
	}
}

// RequiredPages
// Berechnet je nachdem wie viele Einträge pro Seite gewünscht sind die benötigte Seizenanzahl und weist diese dem entsprechenden
// Feld in dem Objekt ListView zu.
func (lv ListView) RequiredPages() int {
	requiredPages := len(lv.EntriesSinceSelDate) / lv.EntriesPerPage
	if len(lv.EntriesSinceSelDate)%lv.EntriesPerPage != 0 || requiredPages == 0 {
		requiredPages += 1
	}
	return requiredPages
}

// NextOccurrences
// Parameter: ein Termin
// Rückgabewert: drei Instanzen des Typs time.Time
// berechnet je nach Wiederholung des Termins und des gewählten Datums, die nächsten drei Vorkommen des Termins.
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
