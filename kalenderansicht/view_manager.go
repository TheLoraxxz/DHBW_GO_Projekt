package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ViewManager struct {
	Tv          TableView
	Lv          ListView
	Username    string
	TerminCache []ds.Termin
}

func (vm *ViewManager) InitViewManager(username string) {
	vm.Tv = *InitTableView(vm.TerminCache)
	vm.Lv = *InitListView(vm.TerminCache)
	vm.Username = username
	vm.TerminCache = ds.GetTermine(vm.Username)
}

// getMaxDays
// Parameter: Monat und Jahr eines Datums
// Rückgabewert: Anzahl der Tage des Monats
// Funktion, die öfters zur hilfe aufgerufen wird
// -> Schaltjahre werden berücksichtigt
func getMaxDays(month, year int) int {
	var days int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		days = 31
	case 4, 6, 9, 11:
		days = 30
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

// filterRepetition
// Parameter: string, eine nummer als string
// Rückgabewert: ein Typ von ds.Repeat, der der entsprechenden Nummer entspricht
func filterRepetition(repStr string) ds.Repeat {
	var rep ds.Repeat
	switch repStr {
	case "1":
		rep = ds.Never
	case "2":
		rep = ds.DAILY
	case "3":
		rep = ds.WEEKLY
	case "4":
		rep = ds.MONTHLY
	case "5":
		rep = ds.YEARLY
	}
	return rep
}

// CreateTermin
// Parameter: eine Post-Request mit Informationen über einen Termin und den Usernamen des Nutzers, der diesen anlegen möchte
// CreateTermin ruft die Funktion zum Erstellen des Termins auf
func (vm *ViewManager) CreateTermin(r *http.Request, username string) {

	//Filtern der Termininfos
	title := r.FormValue("title")
	description := r.FormValue("description")
	repStr := r.FormValue("repeat")

	//Filter das Wiederholungsintervall aus der Antwort
	rep := filterRepetition(repStr)

	//Daten in das richtige Format überführen mithilfe eines Layouts
	layout := "2006-01-02"
	date, _ := time.Parse(layout, r.FormValue("date"))
	endDate, _ := time.Parse(layout, r.FormValue("endDate"))

	//End Date Logik-check
	if repStr == "niemals" || endDate.Before(date) {
		endDate = date
	}

	//Erstelle neuen Termin und füge diesen dem Cache hinzu
	newTermin := ds.CreateNewTermin(title, description, rep, date, endDate, username)
	vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
}

// EditTermin
// Parameter: Post-Request mit den Werten die geändert/gelöscht werden sollen
// Die Funktion ermittelt aus der Post-request die Termine an dem Tag, der bearbeitet wird.
// Anschließend wird für jeden Termin an dem tag ermittelt, ob dieser bearbeitet oder gelöscht werden soll:
//
//	-> Berbeiten: Dei neuen Termininfos werden ermittelt und ein neuer Termin erstellt. Der alte Termin wird gelöscht.
//	-> Löschen: Der Termin wird gelöscht.
func (vm *ViewManager) EditTermin(r *http.Request, username string) {

	//Den alten Titel zum Löschen ermitteln
	oldTitle := r.FormValue("oldTitle")

	//Filtern des gewünschten Modus: bearbeiten oder Löschen
	mode := r.FormValue("editing")

	//Wenn der Modus Bearbeiten ist, müssen die übergebenen Values gelesen werden und der Termin
	//Wenn der Modus Löschen ist, kann der Termin gleich gelöscht werden (ohne Werteauslesen)
	if strings.Contains(mode, "Bearbeiten") {
		title := r.FormValue("title")
		description := r.FormValue("description")
		repStr := r.FormValue("repeat")

		//Filter das Wiederholungsintervall aus der Antwort
		rep := filterRepetition(repStr)

		//Daten in das richtige Format überführen mithilfe eines Layouts
		layout := "2006-01-02"
		date, _ := time.Parse(layout, r.FormValue("date"))
		endDate, _ := time.Parse(layout, r.FormValue("endDate"))

		editetTermin := ds.Termin{
			title,
			description,
			rep,
			date,
			endDate,
		}

		//Lösche Print anweisung
		fmt.Println(editetTermin)

		//Erstelle neuen Termin mit den geänderten und füge diesen dem Cache hinzu
		newTermin := ds.CreateNewTermin(title, description, rep, date, endDate, username)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)

	}
	//Löschen des alten Termins
	vm.TerminCache = ds.DeleteFromCache(vm.TerminCache, oldTitle, vm.Username)
}

/**********************************************************************************************************************
Hier Folgen Funktionen, die dem Handeln der Tabellenansicht/TableView dienen.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// JumpMonthBack
// Springt einen Monat in der Webseiten Ansicht vor
func (vm *ViewManager) TvJumpMonthBack() {
	vm.Tv.JumpMonthBack()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvJumpMonthFor
// Springt einen Monat in der Webseiten Ansicht zurück
func (vm *ViewManager) TvJumpMonthFor() {
	vm.Tv.JumpMonthFor()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvJumpToYear
// Springt zu einem bestimmten Jahr
func (vm *ViewManager) TvJumpToYear(summand int) {
	vm.Tv.JumpToYear(summand)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvSelectMonth
// Parameter: vom benutzer auf der Webseite gewählter Monat
// Setzt den Monat auf den gewünschten Monat
func (vm *ViewManager) TvSelectMonth(monat time.Month) {
	vm.Tv.SelectMonth(monat)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvJumpToToday
// Springt in der Webseiten Ansicht auf den heutigen Monat
func (vm *ViewManager) TvJumpToToday() {
	vm.Tv.getFirstDayOfMonth(time.Now())
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

/**********************************************************************************************************************
Hier Folgen Funktionen, die dem Handeln der Listenansicht/ListView dienen.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
