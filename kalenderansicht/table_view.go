package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"time"
)

type TableView struct {
	ShownDate    time.Time
	Username     string
	MonthEntries []dayInfos
}

// initTableView
// Rückgabewert: Pointer auf ein Objekt TableView
// Dient zur Initialisierung der TableView zu Start des Programms.
// Zu Begin wird diese auf den ersten Tag des aktuellen Monats gesetzt.
func InitTableView() *TableView {
	var tv = new(TableView)
	tv.ShownDate = tv.getFirstDayOfMonth(time.Now())
	return tv
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die der Tabellenansicht dienen.
Diese werden im html-template angesprochen.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// ShowYear
// gibt das Jahr des betrachten Monats zurück, um dies auf der Webseite anzuzeigen
func (tv TableView) ShowYear() int {
	return tv.ShownDate.Year()
}

// ShowMonth
// gibt den betrachten Monat zurück, um ihn auf der Webseite anzuzeigen
func (tv TableView) ShowMonth() time.Month {
	return tv.ShownDate.Month()
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Navigieren innerhalb der Tabellenansicht dienen.
Alle Funktionen, die mit dem Navigieren zu tun haben, sind an das Objekt Tabellenansicht.
gebunden und greifen somit auf dasselbe Datum zu.
Die Funktionen werden mit Hilfe von JavaSkript aufgerufen, wenn das entsprechende Objekt auf der Webseite angeklickt wird.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// JumpMonthBack
// Springt einen Monat in der Webseiten Ansicht vor
func (tv *TableView) JumpMonthBack() {
	tv.ShownDate = tv.ShownDate.AddDate(0, 1, 0)
	tv.CreateTerminTable()
}

// JumpMonthFor
// Springt einen Monat in der Webseiten Ansicht zurück
func (tv *TableView) JumpMonthFor() {
	tv.ShownDate = tv.ShownDate.AddDate(0, -1, 0)
	tv.CreateTerminTable()
}
func (tv *TableView) JumpToYear(summand int) {
	tv.ShownDate = tv.ShownDate.AddDate(summand, 0, 0)
	tv.CreateTerminTable()
}

// SelectMonth
// Parameter: vom benutzer auf der Webseite gewählter Monat
// Setzt den Monat auf den gewünschten Monat
func (tv *TableView) SelectMonth(monat time.Month) {
	jahr := tv.ShownDate.Year()
	tv.ShownDate = time.Date(
		jahr,
		monat,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	tv.CreateTerminTable()
}

// JumpToToday
// Springt in der Webseiten Ansicht auf den heutigen Monat
func (tv *TableView) JumpToToday() {
	tv.ShownDate = tv.getFirstDayOfMonth(time.Now())
	tv.CreateTerminTable()
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Filtern und Anzeigen der Termine in der Tabellenansicht dienen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// dayInfos
// Enthält das Datum eines spezifischen tag und ein Slice mit den Terminen an diesem Tag
type dayInfos struct {
	Day        time.Time //kann raus, da über Index abfragbar
	Dayentries []ds.Termin
	NeedsBreak func(day time.Time) bool
	IsToday    func(day time.Time) bool
}

// CreateTerminTable
// Lädt alle Termine des Benutzers.
// Ruft die Funktion zum Filtern der Termine auf, die in den betrachteten Monat fallen.
// Weißt diese Termine dem Feld MonthEntries von TableView zu.
func (tv *TableView) CreateTerminTable() {
	termins := ds.GetTermine(tv.Username)
	tv.MonthEntries = tv.FilterCalendarEntries(termins)
}

// FilterCalendarEntries
// Parameter: Slice, die alle Termine des Users enthält
// Rückgabewert: Gefilterte Slice, deren Länge der Monatslänge entspricht und für jeden Tag die Termininfos enthält
// Der Termintag entspricht dabei dem Index -1 in dem Slice
func (tv TableView) FilterCalendarEntries(termins []ds.Termin) []dayInfos {

	monthStartDate := tv.ShownDate
	monthEndDate := tv.getLastDayOfMonth()

	entriesForThisMonth := make([]dayInfos, getMaxDays(int(tv.ShownDate.Month()), tv.ShownDate.Year()))
	for _, termin := range termins {
		if (termin.Date.Before(monthEndDate) || termin.Date.Equal(monthEndDate)) && (termin.EndDate.After(monthStartDate) || termin.EndDate.Equal(monthStartDate)) {
			switch termin.Recurring {
			case ds.Never, ds.YEARLY, ds.MONTHLY:
				monthDay := termin.Date.Day()
				entriesForThisMonth[monthDay-1].Dayentries = append(entriesForThisMonth[monthDay-1].Dayentries, termin)
				// Vom Start des Termins wird je eine Woche dazu addiert und geprüft, ob dieses neue Datum in den betrachteten Monat fällt
				// Fällt der Termin in den gewählten Zeitraum, wird der termin in die Slice hinzugefügt
			case ds.DAILY, ds.WEEKLY:
				folgeTermin := termin.Date
				for folgeTermin.Before(termin.EndDate) || folgeTermin.Equal(termin.EndDate) {
					if (folgeTermin.Before(monthEndDate) || folgeTermin.Equal(monthEndDate)) && (folgeTermin.After(monthStartDate) || folgeTermin.Equal(monthStartDate)) {
						monthDay := folgeTermin.Day()
						entriesForThisMonth[monthDay-1].Dayentries = append(entriesForThisMonth[monthDay-1].Dayentries, termin)
					}
					if termin.Recurring == ds.WEEKLY {
						folgeTermin = folgeTermin.AddDate(0, 0, 7)
					} else {
						folgeTermin = folgeTermin.AddDate(0, 0, 1)
					}
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
		//Delete: is just for testing Webview
		// entriesForThisMonth[i].Dayentries = append(entriesForThisMonth[i].Dayentries, ds.Termin{Title: "Test1", Description: "boa", Recurring: ds.Repeat((i % 5)), Date: monthStartDate})
		// entriesForThisMonth[i].Dayentries = append(entriesForThisMonth[i].Dayentries, ds.Termin{Title: "Test2", Description: "boa", Recurring: ds.Repeat((i % 5)), Date: monthStartDate})

	}
	return entriesForThisMonth
}

// getFirstDayOfMonth
// Hilfsfunktion, die den ersten Tag des Monats liefert
func (tv TableView) getFirstDayOfMonth(specificDate time.Time) time.Time {
	return time.Date(
		specificDate.Year(),
		specificDate.Month(),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

// getLastDayOfMonth
// liefert den letzten Tag des (auf der Webseite) betrachteten Monats
// Rückgabe: letzter Tag des Monats, welcher in der Webansicht zu sehen ist.
func (tv TableView) getLastDayOfMonth() time.Time {
	maxDays := getMaxDays(int(tv.ShownDate.Month()), tv.ShownDate.Year())
	return time.Date(
		tv.ShownDate.Year(),
		tv.ShownDate.Month(),
		maxDays,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

// MonthStarts
// dient dazu, dass die Tabelle beim richtigen Wochentag beginnt.
// Hierzu wird eine Slice erstellt, deren Länge so groß ist wie die Anzahl der zu überspringenden Tage.
// über die Slice wird iteriert wird und für jedes Objekt wir ein leeres Tabellenfeld erstellt
func (tv TableView) MonthStarts() []int {
	firstDayOfMonth := tv.ShownDate
	sliceSize := firstDayOfMonth.Weekday() - 1
	if firstDayOfMonth.Weekday()-1 < 0 {
		sliceSize = 6
	}
	return make([]int, sliceSize)
}

// NeedsBreak
// Parameter: Datum
// Rückgabewert: bool, handelt es sich bei dem Tag um einen Sonntag?
// Ist der Tag ein Sonntag, wird ein Tabellenumbruch in der Html-Datei benötigt.
func NeedsBreak(day time.Time) bool {
	return day.Weekday() == time.Sunday
}

// IsToday
// Parameter: Datum
// Rückgabewert: bool, handelt es sich bei dem datum um heute?
// Handelt es sich um heute, muss das Tabellenfeld gekennzeichnet werden.
func IsToday(day time.Time) bool {
	dayYear, dayMonth, dayNr := day.Date()
	todayYear, todayMonth, todayNr := time.Now().Date()
	return (dayYear == todayYear) && (dayMonth == todayMonth) && (dayNr == todayNr)
}
