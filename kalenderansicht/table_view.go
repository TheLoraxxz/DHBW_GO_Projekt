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
func InitTableView() TableView {
	var tv = new(TableView)
	tv.ShownDate = tv.getFirstDayOfMonth(time.Now())
	tv.CreateTerminTable()
	return *tv
}

// getFirstDayOfMonth
// liefert den ersten Tag des (auf der Webseite) betrachteten Monats
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

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die der Tabellenansicht dienen.
Diese werden im html-template angesprochen.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// ShowYear
// gibt das Jahr des auf der Webseite zu sehenden Monats zurück
func (tv TableView) ShowYear() int {
	return tv.ShownDate.Year()
}

// ShowMonth
// gibt das Jahr des auf der Webseite zu sehenden Monats zurück
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
func (tv *TableView) CreateTerminTable() {
	termins := ds.GetTermine(tv.Username)
	tv.MonthEntries = tv.FilterCalendarEntries(termins)
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
