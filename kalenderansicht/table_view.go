package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"time"
)

type TableView struct {
	ShownDate time.Time
	Username  string
}

// initTableView
// Rückgabewert: Pointer auf ein Objekt TableView
// Dient zur Initialisierung der TableView zu Start des Programms.
// Zu Begin wird diese auf den ersten Tag des aktuellen Monats gesetzt.
func InitTableView() TableView {
	var tv = new(TableView)
	tv.ShownDate = tv.getFirstDayOfMonth(time.Now())
	return *tv
}

// getFirstDayOfMonth
// liefert den ersten Tag des (auf der Webseite) betrachteten Monats
func (c TableView) getFirstDayOfMonth(specificDate time.Time) time.Time {
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
func (c TableView) ShowYear() int {
	return c.ShownDate.Year()
}

// ShowMonth
// gibt das Jahr des auf der Webseite zu sehenden Monats zurück
func (c TableView) ShowMonth() time.Month {
	return c.ShownDate.Month()
}

/**********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Navigieren innerhalb der Tabellenansicht dienen.
Alle Funktionen, die mit dem Navigieren zu tun haben, sind an das Objekt Tabellenansicht.
gebunden und greifen somit auf dasselbe Datum zu.
Die Funktionen werden mit Hilfe von JavaSkript aufgerufen, wenn das entsprechende Objekt auf der Webseite angeklickt wird.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// JumpMonthBack
// Springt einen Monat in der Webseiten Ansicht vor
func (c *TableView) JumpMonthBack() {
	c.ShownDate = c.ShownDate.AddDate(0, 1, 0)
}

// JumpMonthFor
// Springt einen Monat in der Webseiten Ansicht zurück
func (c *TableView) JumpMonthFor() {
	c.ShownDate = c.ShownDate.AddDate(0, -1, 0)
}
func (c *TableView) JumpToYear(summand int) {
	c.ShownDate = c.ShownDate.AddDate(summand, 0, 0)
}

// SelectMonth
// Parameter: vom benutzer auf der Webseite gewählter Monat
// Setzt den Monat auf den gewünschten Monat
func (c *TableView) SelectMonth(monat time.Month) {
	jahr := c.ShownDate.Year()
	c.ShownDate = time.Date(
		jahr,
		monat,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

// JumpToToday
// Springt in der Webseiten Ansicht auf den heutigen Monat
func (c *TableView) JumpToToday() {
	c.ShownDate = c.getFirstDayOfMonth(time.Now())
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
func (c TableView) CreateTerminTable() []dayInfos {
	termins := ds.GetTermine(c.Username)
	return c.FilterCalendarEntries(termins)
}

// getLastDayOfMonth
// liefert den letzten Tag des (auf der Webseite) betrachteten Monats
// Rückgabe: letzter Tag des Monats, welcher in der Webansicht zu sehen ist.
func (c TableView) getLastDayOfMonth() time.Time {
	maxDays := getMaxDays(int(c.ShownDate.Month()), c.ShownDate.Year())
	return time.Date(
		c.ShownDate.Year(),
		c.ShownDate.Month(),
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
func (c TableView) MonthStarts() []int {
	firstDayOfMonth := c.ShownDate
	sliceSize := firstDayOfMonth.Weekday() - 1
	if firstDayOfMonth.Weekday()-1 < 0 {
		sliceSize = 6
	}
	return make([]int, sliceSize)
}

// NeedsBreak
// Parameter: Datum
// Rückgabewert: bool, handelt es sich bei dem tag um einen Sonntag?
// Ist der Tag ein Sonntag wird ein Tabellenumbruch in der Html-Datei benötigt.
func NeedsBreak(day time.Time) bool {
	return day.Weekday() == time.Sunday
}

// IsToday
// Parameter: Datum
// Rückgabewert: bool, handelt es sich bei dem datum um heute?
// Handelt es sich um heute muss das Tabellenfeld gekennzeichnet werden.
func IsToday(day time.Time) bool {
	dayYear, dayMonth, dayNr := day.Date()
	todayYear, todayMonth, todayNr := time.Now().Date()
	return (dayYear == todayYear) && (dayMonth == todayMonth) && (dayNr == todayNr)
}
