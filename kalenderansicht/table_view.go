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
}

// CreateTerminTable
func (c TableView) CreateTerminTable() []dayInfos {
	termins := ds.GetTermine(c.Username)
	return c.FilterCalendarEntries(termins)
}

// getLastDayOfMonth
// liefert den letzten Tag des (auf der Webseite) betrachteten Monats
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
	return entriesForThisMonth
}

// MonthStarts
// dient dazu, dass die Tabelle beim richtigen Wochentag beginn
// Hierzu wird eine Slice erstellt.
// über die Slice wird iteriert wird und für jedes Objekt wir ein leeres Tabellenfeld erstellt
func (c TableView) MonthStarts() []int {
	firstDayOfMonth := c.ShownDate
	return make([]int, firstDayOfMonth.Weekday()-1)
}

// MonthStarts
// dient dazu, dass die Tabelle beim richtigen Wochentag beginn
// Hierzu wird eine Slice erstellt.
// über die Slice wird iteriert wird und für jedes Objekt wir ein leeres Tabellenfeld erstellt
func (c TableView) NeedsBreak(indexDay int) int {
	firstDayOfMonth := c.ShownDate.Weekday()
	var needsBreak int
	if indexDay%(7-int(firstDayOfMonth)) == 0 {
		needsBreak = 1
	} else {
		needsBreak = 0
	}
	return needsBreak
}

/*
// tableStructWrapper
// Dient dazu, die Termine in eine passende tabellenstruktur zu bringen
// Hierzu wird kontrolliert, dass die Tabellenansicht des Monats am entsprechenden Wochentag beginnt
// In diesem Kontext werden wird mehrfach createDayEntry aufgerufen um die einzelnen Tabellenfelder zu erstellen
func (c TableView) tableStructWrapper(entriesForThisMonth []dayInfos) string {
	var monthHtmlStr string
	var weekDay int
	switch entriesForThisMonth[0].day.Weekday() {
	case time.Monday:
		weekDay = 1
		monthHtmlStr += "<tr>"
	case time.Tuesday:
		weekDay = 2
		monthHtmlStr += "<tr>" +
			"<th></th>"
	case time.Wednesday:
		weekDay = 3
		monthHtmlStr += "<tr>" +
			"<th></th>" +
			"<th></th>"
	case time.Thursday:
		weekDay = 4
		monthHtmlStr += "<tr>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>"
	case time.Friday:
		weekDay = 5
		monthHtmlStr += "<tr>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>"
	case time.Saturday:
		weekDay = 6
		monthHtmlStr += "<tr>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>"
	case time.Sunday:
		weekDay = 7
		monthHtmlStr += "<tr>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>" +
			"<th></th>"
	}

	//Hier wird über jeden Tag iteriert und ein entsprechendes tabellenfeld erstellt
	//Nach jedem Sonntag(bzw. 7. Tag) wird eine neue Tabellenreihe angefangen, indem ein </tr>-Element hinzugefügt wird
	for _, dayNr := range entriesForThisMonth {
		if weekDay > 7 {
			weekDay = 1
			monthHtmlStr += "</tr><tr>"
		}
		weekDay += 1
		entriesNr := len(dayNr.dayentries)
		monthHtmlStr += c.createDayEntry(dayNr.day, entriesNr, dayNr.dayentries)
	}
	if weekDay < 7 {
		monthHtmlStr += "</tr>"
	}
	return monthHtmlStr
}

// createDayEntry
// erstellt für jeden Tag ein Tabellenfeld mit einem Button, aud den man klicken kann.
// Bei einem Klick werden die termine in einem Popup angezeigt
// ruft die Funktion createTerminPopup auf, die die einzelnen Html-Popup-strings erzeugt
func (c *TableView) createDayEntry(day time.Time, entriesNr int, entries []ds.Termin) string {
	daynr := day.Day()

	//wenn es sich um den heutigen Tag handelt soll der Button Rot sein
	todayYear, todayMonth, todayDay := time.Now().Date()
	var buttonClr string
	if todayYear == day.Year() && todayMonth == day.Month() && todayDay == day.Day() {
		buttonClr = "class=\"btn btn-danger\""
	} else {
		buttonClr = "class=\"btn btn-primary\""
	}

	//Der Button wird erstellt, auf diesem steht die Tagesnummer sowie Terminanzahl
	dayBtn := "" +
		"<th>" +
		"<button type=\"button\"" +
		buttonClr +
		"data-bs-toggle=\"modal\" " +
		"data-bs-target=" + "#" + strconv.Itoa(daynr) + ">" +
		strconv.Itoa(daynr) + ". Termine: " + strconv.Itoa(entriesNr) +
		"</button>" +
		"</th>"

	//Das zum Button gehörige Popup wird erstellt
	modal := c.createTerminPopup(day, entries)
	//return Html-String für
	return dayBtn + modal
}*/

/*
// createTerminPopup
// Erstellt den Html-String für ein Popup, der die Infos und Termine für einen spezifischen Tag enthält
// Das Popup wird aufgerufen, wenn der User auf den Button im entsprechenden Feld (=Tag) der Kalendertabelle klickt
func (c *TableView) createTerminPopup(day time.Time, dayentries []ds.Termin) string {
	var modalEntries string
	var entrieStr string
	var recurringStr string
	if dayentries != nil {
		for _, entrie := range dayentries {
			recurring := entrie.Recurring
			switch recurring {
			case ds.Never:
				recurringStr = "Niemals"
			case ds.YEARLY:
				recurringStr = "Jährlich"
			case ds.MONTHLY:
				recurringStr = "Monatlich"
			case ds.WEEKLY:
				recurringStr = "Wöchentlich"
			}
			entrieStr = "" +
				"<h6>" + entrie.Title + "</h6>" +
				"<b>Beschreibung: </b>" + entrie.Description +
				"<b>Wiederholung: </b>" + recurringStr
			modalEntries += entrieStr
		}
	} else {
		modalEntries = "<h6>Keine Termine an diesem Tag.</h6>"
	}

	modal :=
		"<div class=\"modal fade\"" +
			"id=" + strconv.Itoa(day.Day()) +
			"tabindex=\"-1\" " +
			"aria-labelledby=\"exampleModalLabel\" " +
			"aria-hidden=\"true\">" +
			"<div class=\"modal-dialog\">" +
			"<div class=\"modal-content\">" +
			"<div class=\"modal-header\">" +
			"<h5 class=\"modal-title\" id=\"exampleModalLabel\">Termine am" + day.String() + "</h5>" +
			"<button type=\"button\" class=\"btn-close\" data-bs-dismiss=\"modal\" aria-label=\"Close\"></button>" +
			"</div>" +
			"<div class=\"modal-body\">" + modalEntries +
			"</div>" +
			"<div class=\"modal-footer\">" +
			"<button type=\"button\" class=\"btn btn-secondary\" data-bs-dismiss=\"modal\">Close</button>" +
			"</div>" +
			"</div>" +
			"</div>" +
			"</div>"
	return modal
}*/
