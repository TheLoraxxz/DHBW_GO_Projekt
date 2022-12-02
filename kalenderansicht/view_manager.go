package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"errors"
	"net/http"
	"time"
)

type ViewManager struct {
	Tv          TableView
	Lv          ListView
	Fv          FilterView
	Username    string
	TerminCache []ds.Termin
}

func InitViewManager(username string) *ViewManager {
	vm := new(ViewManager)
	vm.Username = username
	vm.TerminCache = ds.GetTermine(vm.Username)
	vm.Tv = *InitTableView(vm.TerminCache)
	vm.Lv = *InitListView(vm.TerminCache)
	vm.Fv = *InitFilterView(vm.TerminCache)
	return vm
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
func filterRepetition(repStr string) (rep ds.Repeat, err error) {
	switch repStr {
	case "täglich", "0":
		rep = ds.DAILY
	case "wöchentlich", "1":
		rep = ds.WEEKLY
	case "monatlich", "2":
		rep = ds.MONTHLY
	case "jährlich", "3":
		rep = ds.YEARLY
	case "niemals", "4":
		rep = ds.Never
	default:
		err = errors.New("No_valid_repetition")
	}
	return rep, err
}

// CreateTermin
// Parameter: eine Post-Request mit Informationen über einen Termin und den Usernamen des Nutzers, der diesen anlegen möchte
// CreateTermin ruft die Funktion zum Erstellen des Termins auf
func (vm *ViewManager) CreateTermin(r *http.Request, username string) (err error) {

	//Filtern der Termininfos
	sharedStr := r.FormValue("shared")
	var shared bool
	if sharedStr == "true" {
		shared = true
	} else {
		shared = false
	}
	title := r.FormValue("title")
	if title == "" {
		err = errors.New("Missing_title")
		return
	}

	description := r.FormValue("description")
	if description == "" {
		err = errors.New("Missing_description")
		return
	}
	repStr := r.FormValue("rep")

	//Filter das Wiederholungsintervall aus der Antwort
	rep, err := filterRepetition(repStr)
	if err != nil {
		return err
	}

	//Daten in das richtige Format überführen mithilfe eines Layouts
	layout := "2006-01-02"
	date, err := time.Parse(layout, r.FormValue("date"))
	if err != nil {
		return errors.New("wrong_date_format")
	}

	endDate, err := time.Parse(layout, r.FormValue("endDate"))
	if err != nil {
		return errors.New("wrong_date_format")
	}

	//End Date Logik-check
	if repStr == "niemals" || endDate.Before(date) {
		endDate = date
	}

	//Erstelle neuen Termin und füge diesen dem Cache hinzu
	newTermin := ds.CreateNewTermin(title, description, rep, date, endDate, shared, username)
	vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)

	//Falls es sich um einen Terminvorschlag handelt, muss dieser noch den Terminvorschlägen hinzugefügt werden
	terminfindung.CreateSharedTermin(&newTermin, &username)

	//Anzuzeigende Einträge in den Ansichten aktualisieren
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	return
}

// GetTerminInfos
// Parameter: Request mit Termininfos
// Rückgabewert: Termin, aus dem Cahce mit entsprechender ID
// Die Funktion wird genutzt, um den Termin zu erhalten, der bearbeitet/gelöscht werden soll
func (vm *ViewManager) GetTerminInfos(r *http.Request) (termin ds.Termin, err error) {
	//Filtern der Termin-Id und des zu bearbeitenden Termins aus dem Cache
	id := r.FormValue("ID")
	termin = ds.FindInCacheByID(vm.TerminCache, id)
	if termin.ID == "" {
		err = errors.New("shared_wrong_terminId")
	}

	return termin, err
}

// EditTermin
// Parameter: Post-Request mit den Werten die geändert/gelöscht werden sollen
// Die Funktion ermittelt aus der Post-request die Termine an dem Tag, der bearbeitet wird.
// Anschließend wird für jeden Termin an dem tag ermittelt, ob dieser bearbeitet oder gelöscht werden soll:
//
//	-> Berbeiten: Dei neuen Termininfos werden ermittelt und ein neuer Termin erstellt. Der alte Termin wird gelöscht.
//	-> Löschen: Der Termin wird gelöscht.
func (vm *ViewManager) EditTermin(r *http.Request, username string) (err error) {

	//Die ID zum Löschen ermitteln
	id := r.FormValue("ID")
	if id == "" {
		return errors.New("shared_wrong_terminId")
	}

	//Filtern des gewünschten Modus: bearbeiten oder Löschen
	mode := r.FormValue("editing")
	if !(mode == "editing" || mode == "delete") {
		return errors.New("wrong_editing_mode")
	}

	//egal ob löschen oder bearbeiten, der Termin muss zunächst gelöscht werden
	vm.TerminCache = ds.DeleteFromCache(vm.TerminCache, id, vm.Username)

	//Wenn der Modus editing ist, muss der aktualisierte Termin noch erstellt werden
	if mode == "editing" {
		err = vm.CreateTermin(r, username)
		if err != nil {
			return
		}
	} else {
		//Anzuzeigende Einträge in den Ansichten aktualisieren (dies geschieht auch in vm.CreateTermin(r, username))
		//entfällt deshalb im Bearbeitungsmodus, da sonst doppelter Funktionsaufruf
		vm.Tv.CreateTerminTableEntries(vm.TerminCache)
		vm.Lv.CreateTerminListEntries(vm.TerminCache)
		vm.Fv.CreateTerminFilterEntries(vm.TerminCache)
	}
	return
}

// DeleteSharedTermin
// Parameter: id des terminvorschlags, der gelöscht werden soll; username
// Die Funktion löscht den Termin aus den Teminvorschlägen und aus dem Cache
func (vm *ViewManager) DeleteSharedTermin(id, username string) (err error) {
	//Termin vom Cache löschen
	vm.TerminCache = ds.DeleteFromCache(vm.TerminCache, id, username)
	//Termin aus den Vorschlägen entfernen
	err = terminfindung.DeleteSharedTermin(&id, &username)
	if err != nil {
		return
	}
	//Anzuzeigende Einträge in den Ansichten aktualisieren
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.Fv.CreateTerminFilterEntries(vm.TerminCache)
	return
}

/**********************************************************************************************************************
Hier Folgen Funktionen, die dem Handeln der Tabellenansicht/TableView dienen.
Nach jedem ändern der Ansicht der TableView, müssen die Einträge
des Users dem neu angezeigten Monat entsprechend gefiltert werden.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// TvJumpMonthBack
// Springt einen Monat in der Webseiten Ansicht zurück
func (vm *ViewManager) TvJumpMonthBack() {
	vm.Tv.JumpMonthBack()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvJumpMonthFor
// Springt einen Monat in der Webseiten Ansicht vor
func (vm *ViewManager) TvJumpMonthFor() {
	vm.Tv.JumpMonthFor()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

// TvJumpYearForOrBack
// Parameter: +1 oder -1
// Springt zu ein Jahr vor oder zurück, je nachdem ob der Parameter +1 oder -1 ist
func (vm *ViewManager) TvJumpYearForOrBack(summand int) {
	vm.Tv.JumpYearForOrBack(summand)
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
	vm.Tv.JumpToToday()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)
}

/**********************************************************************************************************************
Hier Folgen Funktionen, die dem Handeln der Listenansicht/ListView dienen.
Nach jedem ändern der Ansicht der ListView, müssen die Einträge
des Users ab dem neu angezeigten Datum entsprechend gefiltert werden.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// LvSelectDate
// Parameter: Datum in string-form
// Funktion setzt das angezeigte Datum der Terminansicht
func (vm *ViewManager) LvSelectDate(dateStr string) (err error) {
	//Datum Filtern und in das richtige Format überführen mithilfe eines Layouts
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return errors.New("wrong_date_format")
	}
	vm.Lv.SelectDate(date)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	return
}

// LvSelectEntriesPerPage
// Parameter: int Wert
// Funktion setzt die angezeigte Einträge-Anzahl in der Terminansicht auf übergebenen int Wert
func (vm *ViewManager) LvSelectEntriesPerPage(amount int) {
	vm.Lv.SelectEntriesPerPage(amount)
	return
}

// LvJumpPageForward
// Funktion springt eine Seite vor in der Listenansicht
func (vm *ViewManager) LvJumpPageForward() {
	vm.Lv.JumpPageForward()
}

// LvJumpPageBack
// Funktion springt eine Seite zurück in der Listenansicht
func (vm *ViewManager) LvJumpPageBack() {
	vm.Lv.JumpPageBack()
}

/**********************************************************************************************************************
Hier Folgen Funktionen, die dem Handeln der Filteransicht/FilterView dienen.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

// FvSelectEntriesPerPage
// Parameter: int Wert
// Funktion setzt die angezeigte Einträge-Anzahl in der Filteransicht auf übergebenen int Wert
func (vm *ViewManager) FvSelectEntriesPerPage(amount int) {
	vm.Fv.SelectEntriesPerPage(amount)
	return
}

// FvJumpPageForward
// Funktion springt eine Seite vor in der Filteransicht
func (vm *ViewManager) FvJumpPageForward() {
	vm.Fv.JumpPageForward()
}

// FvJumpPageBack
// Funktion springt eine Seite zurück in der Filteransicht
func (vm *ViewManager) FvJumpPageBack() {
	vm.Fv.JumpPageBack()
}

// FvFilter
// Parameter: Request mit Strings des Termin-Titels/der Termin-Beschreibung nach der, gefiltert werden soll
func (vm *ViewManager) FvFilter(r *http.Request) {
	filterTitle := r.FormValue("title")
	filterDescription := r.FormValue("description")
	vm.Fv.FilterTermins(filterTitle, filterDescription, vm.TerminCache)
}
