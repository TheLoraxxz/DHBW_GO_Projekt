package kalenderansicht

import (
	"DHBW_GO_Projekt/terminfindung"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

import ds "DHBW_GO_Projekt/dateisystem"

/*
**************************************************************************************************************
Hilfsfunktionen zum zufälligen generieren von Testdaten.
Diese werden auch in den Tests von TableView und ListView genutzt.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
//Slice, der mit Testterminen gefüllt wird, damit dies für mehrere Tests nur einmal durchgeführt werden muss
var testTermine30 []ds.Termin

// newTerminObj erzeugt einen transitiven Termin
func NewTerminObj(title string, description string, rep ds.Repeat, date time.Time, endDate time.Time, shared bool) ds.Termin {

	t := ds.Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        date,
		EndDate:     endDate,
		Shared:      shared,
		ID:          createDummyID(date, endDate)}

	return t
}

// createID erzeugt neue ID für die Testtermine
func createDummyID(dat time.Time, endDat time.Time) string {

	u := time.Now().String()

	id := dat.String() + endDat.String() + u

	//generiert Hash --> gewährleistet hohe Kollisionsfreiheit bei IDs
	bytes, _ := bcrypt.GenerateFromPassword([]byte(id), 14)
	id = string(bytes)

	//Entfernt problematische Chars aus Hash
	id = strings.Replace(id, "/", "E", 99)
	id = strings.Replace(id, ".", "D", 99)

	return id
}

// create30TestTermins
// Rückgabewert: Slice mit 30 testterminen
// Funktionen wird benötigt um ein Slice mit Test-Terminen zu Begin des Testvorgangs zu erstellen.
// Dieses wird für mehrere Tests benötigt in: list_view_test & in view_manager_test
func create30TestTermins() []ds.Termin {
	testTerminStarts := time.Now()
	testTerminEnds := time.Now().AddDate(1, 0, 0)

	//Slice mit Testterminen erstellen
	testTermine30 := make([]ds.Termin, 0, 30)
	// 5 testTermine30 erstellen
	testTermin1 := NewTerminObj("testTermin1", "test hi", ds.MONTHLY, testTerminStarts, testTerminEnds, false)
	testTermin2 := NewTerminObj("testTermin2", "test hi", ds.YEARLY, testTerminStarts, testTerminEnds, false)
	testTermin3 := NewTerminObj("testTermin3", "test", ds.WEEKLY, testTerminStarts, testTerminEnds, false)

	for i := 0; i < 10; i++ {
		testTermine30 = append(testTermine30, testTermin1)
		testTermine30 = append(testTermine30, testTermin2)
		testTermine30 = append(testTermine30, testTermin3)
	}
	return testTermine30
}

// generateRandomDate
// generiert ein komplett zufälliges Datum (bis zum Jahr 3000)
func generateRandomDate() time.Time {
	month := rand.Intn(13-1) + 1
	var day int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		day = rand.Intn(32-1) + 1
	case 4, 6, 9, 11:
		day = rand.Intn(31-1) + 1
	case 2:
		day = rand.Intn(29-1) + 1
	}
	return time.Date(
		rand.Intn(3000),
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

// generateRandomDateInSpecificMonth
// generiert ein zufälliges Datum eines spezifischen Monats und Jahres
func generateRandomDateInSpecificMonth(year int, month time.Month) time.Time {
	var day int
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		day = rand.Intn(32-1) + 1
	case 4, 6, 9, 11:
		day = rand.Intn(31-1) + 1
	case 2:
		day = rand.Intn(29-1) + 1
	}
	randomDate := time.Date(
		year,
		month,
		day,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	return randomDate
}

// createWeeklyTestTermin
// generiert eines wöchentlichen/jährlichen/monatlichen Termins, damit diese im Kalender anzeigbar sind,
// um die Funktion der Navigation der Webseite zu testen
func createTestTermin(repeat ds.Repeat) *ViewManager {
	vm := InitViewManager("testuser")
	year := time.Now().Year()
	day := time.Now().Day()
	month := time.Now().Month()

	switch repeat {
	case ds.WEEKLY:
		newTermin := NewTerminObj("test Title", "test", ds.WEEKLY, createSpecificDate(year, day, month), createSpecificDate(year+1, day, month), false)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	case ds.YEARLY:
		newTermin := NewTerminObj("test Title", "test", ds.YEARLY, createSpecificDate(year, day, month), createSpecificDate(year+2, day, month), false)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	case ds.MONTHLY:
		newTermin := NewTerminObj("test Title", "test", ds.MONTHLY, createSpecificDate(year-1, day, month), createSpecificDate(year+1, day, month), false)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	}
	return vm
}

// createSpecificDate
// generiert ein spezifisches Datum
func createSpecificDate(year, day int, month time.Month) time.Time {
	testDate := time.Date(
		year,
		month,
		day,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	return testDate
}

/*
**************************************************************************************************************
Hier Folgen die Tests zum Termine erstellen/bearbeiten/löschen und die dafür benötigte Hilfsfunktion filterRepetioition
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testfilterRepetition(t *testing.T) {
	assert.Equal(t, ds.Never, filterRepetition("niemals"))
	assert.Equal(t, ds.DAILY, filterRepetition("täglich"))
	assert.Equal(t, ds.WEEKLY, filterRepetition("wöchentlich"))
	assert.Equal(t, ds.MONTHLY, filterRepetition("monatlich"))
	assert.Equal(t, ds.YEARLY, filterRepetition("jährlich"))
}

func testCreateTermin(t *testing.T) {
	vm := new(ViewManager)
	vm.Username = "testuser"
	//Erstellen der Termininfos, die über die Request gesendet werden
	data := url.Values{}
	data.Add("title", "Test Termin")
	data.Add("description", "Spaßiger Termin")
	data.Add("repeat", "2") //Der Wert 2 entspricht der Wiederholung "täglich"
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2030-11-11")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "?terminErstellen", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	termin := NewTerminObj("Test Termin", "Spaßiger Termin", ds.DAILY, createSpecificDate(2022, 11, 11), createSpecificDate(2030, 11, 11), false)

	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)
	vm.CreateTermin(r, vm.Username)

	//Testen, ob ein Termin dem Cache hinzugefügt worden ist
	assert.Equal(t, oldLen+1, len(vm.TerminCache), "Die Länge sollte um eins erhöht worden sein.")
	//Testen, ob Termin im Cache mit den Infos aus dem erstellten Termin übereinstimmen
	assert.Equal(t, termin.Title, vm.TerminCache[0].Title, "Die Termin-Titel sollten überein stimmen.")
	assert.Equal(t, termin.Description, vm.TerminCache[0].Description, "Die Termin-Beschreibungen sollten überein stimmen.")
	assert.Equal(t, termin.Date, vm.TerminCache[0].Date, "Die Termin-Startdaten sollten überein stimmen.")
	assert.Equal(t, termin.Recurring, vm.TerminCache[0].Recurring, "Die Termin-Wiederholungen sollten überein stimmen.")
	assert.Equal(t, termin.EndDate, vm.TerminCache[0].EndDate, "Die Termin-Enddaten sollten überein stimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

// testCreateTerminLogicCheck
// Wenn der Nutzer einen termin eingibt, dessen Ende zeitlich vor dem Startdatum ist, muss das Enddatum auf das Startdatum
// gesetzt werden
func testCreateTerminLogicCheck(t *testing.T) {
	vm := new(ViewManager)
	vm.Username = "testuser"
	//Erstellen der Termininfos, die über die Request gesendet werden
	data := url.Values{}
	data.Add("title", "Test Termin")
	data.Add("description", "Spaßiger Termin")
	data.Add("repeat", "2") //Der Wert 2 entspricht der Wiederholung "täglich"
	data.Add("date", "2022-11-11")
	//Hier hat der Nutzer einen Endtermin eingegeben, der zeitlich vor dem Starttermin ist
	data.Add("endDate", "2021-11-11")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "?terminErstellen", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	termin := NewTerminObj("Test Termin", "Spaßiger Termin", ds.DAILY, createSpecificDate(2022, 11, 11), createSpecificDate(2022, 11, 11), false)

	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)
	vm.CreateTermin(r, vm.Username)

	//Testen, ob ein Termin dem Cache hinzugefügt worden ist
	assert.Equal(t, oldLen+1, len(vm.TerminCache), "Die Länge sollte um eins erhöht worden sein.")
	//Testen, ob Termin im Cache dem neuen Termin entspricht
	assert.Equal(t, termin.Title, vm.TerminCache[0].Title, "Die Termin-Titel sollten überein stimmen.")
	assert.Equal(t, termin.Description, vm.TerminCache[0].Description, "Die Termin-Beschreibungen sollten überein stimmen.")
	assert.Equal(t, termin.Date, vm.TerminCache[0].Date, "Die Termin-Startdaten sollten überein stimmen.")
	assert.Equal(t, termin.Recurring, vm.TerminCache[0].Recurring, "Die Termin-Wiederholungen sollten überein stimmen.")
	assert.Equal(t, termin.EndDate, vm.TerminCache[0].EndDate, "Die Termin-Enddaten sollten überein stimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

func testGetTerminInfos(t *testing.T) {
	//Erstelle einen wöchentlichen Termin zum Testen
	vm := createTestTermin(ds.WEEKLY)
	termin := vm.TerminCache[0]

	//Erstellen der Request-Werte
	data := url.Values{}
	data.Add("ID", termin.ID)

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "/editor", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assert.Equal(t, termin, vm.GetTerminInfos(r), "Die Termine sollten identisch sein.")
}

func testEditTerminDelete(t *testing.T) {
	vm := new(ViewManager)
	vm.Username = "testuser"
	//Testtermin erstellen und der Request hinzufügen
	termin := NewTerminObj("Test Termin", "Spaßiger Termin", ds.DAILY, createSpecificDate(2022, 11, 11), createSpecificDate(2022, 11, 11), false)
	vm.TerminCache = append(vm.TerminCache, termin)

	//Erstellen der Termin-infos, die über die Request gesendet werden
	data := url.Values{}

	//Erstellen der Lösch-Request, Wert 2 entspricht einer Lösch-Anfrage
	data = url.Values{}
	data.Add("ID", termin.ID)
	data.Add("editing", "2")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "../editor", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)

	//Termin bearbeiten
	vm.EditTermin(r, vm.Username)

	//Testen, ob der Termin aus dem Cache gelöscht worden ist
	assert.Equal(t, oldLen-1, len(vm.TerminCache), "Die Länge sollte um 1 reduziert worden sein.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}
func testEditTerminEdit(t *testing.T) {
	vm := new(ViewManager)
	vm.Username = "testuser"

	//Testtermin erstellen und der Request hinzufügen
	termin := NewTerminObj("Test Termin", "Spaßiger Termin", ds.DAILY, createSpecificDate(2022, 11, 11), createSpecificDate(2022, 11, 11), false)
	vm.TerminCache = append(vm.TerminCache, termin)

	//Länge des TerminCaches vor dem Bearbeiten
	oldLen := len(vm.TerminCache)

	//Erstellen der Termin-infos, die über die Request gesendet werden
	data := url.Values{}
	//Erstellen der Bearbeiten-Request
	data = url.Values{}
	data.Add("editing", "1")
	data.Add("title", "Test Termin Bearbeitet")
	data.Add("description", "Spaßiger bearbeiteter Termin")
	data.Add("rep", "wöchentlich")
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2023-11-01")
	data.Add("ID", termin.ID)

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "?termineBearbeiten", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//bearbeiteter Termin zum Überprüfen
	editetTermin := NewTerminObj("Test Termin Bearbeitet", "Spaßiger bearbeiteter Termin", ds.WEEKLY, createSpecificDate(2022, 11, 11), createSpecificDate(2023, 01, 11), false)

	//Termin bearbeiten
	vm.EditTermin(r, vm.Username)

	//Testen, ob die Terminanzahl in dem Cache gleich geblieben ist
	assert.Equal(t, oldLen, len(vm.TerminCache), "Die Länge sollte gleich geblieben sein.")
	//Testen, ob Termin im Cache dem neuen bearbeiteten Termin entspricht
	assert.Equal(t, editetTermin.Title, vm.TerminCache[0].Title, "Die Termin-Titel sollten überein stimmen.")
	assert.Equal(t, editetTermin.Description, vm.TerminCache[0].Description, "Die Termin-Beschreibungen sollten überein stimmen.")
	assert.Equal(t, editetTermin.Date, vm.TerminCache[0].Date, "Die Termin-Startdaten sollten überein stimmen.")
	assert.Equal(t, editetTermin.Recurring, vm.TerminCache[0].Recurring, "Die Termin-Wiederholungen sollten überein stimmen.")
	assert.Equal(t, editetTermin.EndDate, vm.TerminCache[0].EndDate, "Die Termin-Enddaten sollten überein stimmen.")

	//Löschen des Termins
	vm.TerminCache = ds.DeleteFromCache(vm.TerminCache, vm.TerminCache[0].ID, vm.Username)
}
func testDeleteSharedTermin(t *testing.T) {
	vm := new(ViewManager)
	vm.Username = "testuser"
	testTermin := NewTerminObj("test go", "hui", ds.DAILY, time.Now(), time.Now().AddDate(1, 0, 0), true)

	vm.TerminCache = append(vm.TerminCache, testTermin)
	terminfindung.CreateSharedTermin(&testTermin, &vm.Username)
	assert.Equal(t, testTermin, vm.TerminCache[0], "Der Termin sollte auf dem Cache sein.")

	vm.DeleteSharedTermin(testTermin.ID, vm.Username)
	assert.Equal(t, 0, len(vm.TerminCache), "Der Cache sollte leer sein.")

}

/*
**************************************************************************************************************
Hier Folgen die Tests zum Managen der TableView:
	Es wird getestet, ob die Termine richtig neu gefiltert werden, wenn sich etwas an der darstellenden Zeit ändert.
	Das die Termine sich ändern wurde bereits in den Tests zur TableView sicher gestellt.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testTvJumpMonthFor(t *testing.T) {

	//Erstelle einen wöchentlichen Termin zum Testen
	vm := createTestTermin(ds.WEEKLY)

	//Der in der TableView angezeigte Monat ist heute
	vm.TvJumpToToday()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//der Termin startet heute und findet wöchentlich statt
	//die Schleife wird solange ausgeführt, bis der Monat einmal vorgesprungen wurde
	day := time.Now().Day()
	jumpedMonthFor := false
	for !jumpedMonthFor {
		if day < vm.Tv.getLastDayOfMonth().Day() {
			assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[day-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
			day = day + 7
		}
		//Springe einen Monat vor
		vm.TvJumpMonthFor()
		//setze Tag auf den im Monat
		day = day - vm.Tv.getLastDayOfMonth().Day()
		//JumpedMonthFor Variable auf true setzen
		jumpedMonthFor = true
	}
}
func testTvJumpMonthBack(t *testing.T) {

	//Erstelle einen monatlichen Termin zum Testen: ab heute vor einem Jahr
	vm := createTestTermin(ds.MONTHLY)

	//Der in der TableView angezeigte Monat ist der Monat heute
	vm.TvJumpToToday()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet am 2.11.2021cund findet monatlich statt
	//Monat ist nun November 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe einen Monat vor
	vm.TvJumpMonthBack()

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet ab 02.11.2021 und findet moantlich statt
	//Monat ist nun Oktober 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
}

func testTvJumpYearForOrBack(t *testing.T) {
	//Jährlichen testtermin erstellen: ab heute bis in zwei Jahren
	vm := createTestTermin(ds.YEARLY)

	//Der in der TableView angezeigte Monat ist heute
	vm.TvJumpToToday()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe ein Jahr vor (zu 2023)
	vm.TvJumpYearForOrBack(1)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe ein Jahr zurück (zu 2022)
	vm.TvJumpYearForOrBack(1)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
}

func testTvSelectMonth(t *testing.T) {
	//Monatlichen testtermin erstellen: ab heute -1 Jahr
	vm := createTestTermin(ds.MONTHLY)

	//Der in der TableView angezeigte Monat ist heute
	vm.TvJumpToToday()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe zu einem beliebigen Monat des Jahres 2022
	month := time.Month(rand.Intn(13))
	vm.TvSelectMonth(month)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
}

func testTvJumpToToday(t *testing.T) {
	//Monatlichen testtermin erstellen: ab heute -1 Jahr
	vm := createTestTermin(ds.MONTHLY)

	vm.TvJumpToToday()

	//Der in der TableView angezeigte Monat wird auf einen Tag in der Vergangenheit gesetzt
	vm.Tv.JumpMonthBack()
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe zu heute
	vm.TvJumpToToday()

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[time.Now().Day()-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
}

/*
**************************************************************************************************************
Hier Folgen die Tests zum Managen der ListView
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testLvSelectDate(t *testing.T) {
	vm := new(ViewManager)
	newDate := generateRandomDate()
	vm.LvSelectDate(newDate.Format("2006-01-02"))
	assert.Equal(t, newDate, vm.Lv.SelectedDate, "Die zwei Daten sollten identisch sein.")
}

func testLvSelectEntriesPerPage(t *testing.T) {
	vm := new(ViewManager)
	entriesPerPage := 5
	vm.LvSelectEntriesPerPage(entriesPerPage)
	assert.Equal(t, entriesPerPage, vm.Lv.EntriesPerPage, "Die Anzahl der Einträge pro Seite sollte 5 sein.")
}

func testLvJumpPageForward(t *testing.T) {
	//Test-termin erstellen: ab heute
	vm := createTestTermin(ds.WEEKLY)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.LvJumpPageForward()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1 sein, da es nur einen Eintrag gibt.")

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}
	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < len(testTermine30); i++ {
		vm.TerminCache = ds.AddToCache(testTermine30[i], vm.TerminCache)
	}

	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.LvJumpPageForward()
	assert.Equal(t, 2, vm.Lv.CurrentPage, "Die Seite sollte 2 sein.")

}

func testLvJumpPageBack(t *testing.T) {
	//Test-termin erstellen: ab heute
	vm := createTestTermin(ds.WEEKLY)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.Lv.JumpPageForward()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1 sein, da die Seitennummer nicht kleiner als 1 sein kann..")

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < len(testTermine30); i++ {
		vm.TerminCache = ds.AddToCache(testTermine30[i], vm.TerminCache)
	}
	vm.Lv.CreateTerminListEntries(vm.TerminCache)

	//Aktuelle Seite auf Seite 2 setzten
	vm.Lv.CurrentPage = 2
	vm.LvJumpPageBack()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1sein.")
}

/*
**************************************************************************************************************
Hier Folgen die Tests zum Managen der FilterView
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testFvSelectEntriesPerPage(t *testing.T) {
	vm := new(ViewManager)
	entriesPerPage := 5
	vm.FvSelectEntriesPerPage(entriesPerPage)
	assert.Equal(t, entriesPerPage, vm.Fv.EntriesPerPage, "Die Anzahl der Einträge pro Seite sollte 5 sein.")
}

func testFvJumpPageForward(t *testing.T) {
	//Test-termin erstellen: ab heute
	vm := createTestTermin(ds.WEEKLY)
	vm.Fv.CreateTerminFilterEntries(vm.TerminCache)
	vm.FvJumpPageForward()
	assert.Equal(t, 1, vm.Fv.CurrentPage, "Die Seite sollte 1 sein, da es nur einen Eintrag gibt.")

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}
	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < len(testTermine30); i++ {
		vm.TerminCache = ds.AddToCache(testTermine30[i], vm.TerminCache)
	}

	vm.Fv.CreateTerminFilterEntries(vm.TerminCache)
	vm.FvJumpPageForward()
	assert.Equal(t, 2, vm.Fv.CurrentPage, "Die Seite sollte 2 sein.")
}
func testFvJumpPageBack(t *testing.T) {
	//Test-termin erstellen: ab heute
	vm := createTestTermin(ds.WEEKLY)
	vm.Fv.CreateTerminFilterEntries(vm.TerminCache)
	vm.Fv.JumpPageForward()
	assert.Equal(t, 1, vm.Fv.CurrentPage, "Die Seite sollte 1 sein, da die Seitennummer nicht kleiner als 1 sein kann..")

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < len(testTermine30); i++ {
		vm.TerminCache = ds.AddToCache(testTermine30[i], vm.TerminCache)
	}
	vm.Lv.CreateTerminListEntries(vm.TerminCache)

	//Aktuelle Seite auf Seite 2 setzten
	vm.Lv.CurrentPage = 2
	vm.LvJumpPageBack()
	assert.Equal(t, 1, vm.Fv.CurrentPage, "Die Seite sollte 1sein.")
}
func testFvFilter(t *testing.T) {
	//Test-termin erstellen: ab heute
	vm := new(ViewManager)

	todayYear := time.Now().Year()
	todayMonth := time.Now().Month()
	todayDay := time.Now().Day()
	today := time.Date(todayYear, todayMonth, todayDay, 0, 0, 0, 0, time.UTC)

	//Daten für Testtermine erstellen

	//Slice mit Testterminen erstellen, jeder Wiederholungstyp dabei
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = NewTerminObj("test go", "ich", ds.DAILY, today, today.AddDate(1, 0, 0), false)
	testTermine[1] = NewTerminObj("test ist", "lala", ds.WEEKLY, today, today.AddDate(1, 0, 0), false)
	testTermine[2] = NewTerminObj("test eine", "ich bin toll", ds.YEARLY, today, today.AddDate(1, 0, 0), false)
	//dem Cache die Termine hinzufügen
	for i := 0; i < len(testTermine); i++ {
		vm.TerminCache = ds.AddToCache(testTermine[i], vm.TerminCache)
	}
	vm.Fv.CreateTerminFilterEntries(vm.TerminCache)

	// Filter-Request erstellen
	data := url.Values{}
	data.Add("title", "test")
	data.Add("description", "ich")

	r, _ := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	vm.FvFilter(r)
	assert.Equal(t, 2, len(vm.Fv.FilteredTermins), "Es sollten 2 Termine herausgefiltert worden sein.")
}

/*
**************************************************************************************************************
Aufrufen aller Tests
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func TestViewManager(t *testing.T) {
	//createWeeklyTestTermin()
	//Tests zum Erstellens/Bearbeitens/Löschens eines Termins
	t.Run("testRuns filterRepetition", testfilterRepetition)
	t.Run("testRuns CreateTermin", testCreateTermin)
	t.Run("testRuns CreateTerminLogicCheck", testCreateTerminLogicCheck)
	t.Run("testRuns GetTerminInfos", testGetTerminInfos)
	t.Run("testRuns EditTermin-Delete", testEditTerminDelete)
	t.Run("testRuns EditTermin-Edit", testEditTerminEdit)
	t.Run("testRuns DeleteShared-Termin", testDeleteSharedTermin)

	//Tests zum Managen der TableView
	t.Run("testRuns TvJumpMonthFor", testTvJumpMonthFor)
	t.Run("testRuns TvJumpMonthBack", testTvJumpMonthBack)
	t.Run("testRuns TvJumpYearForOrBack", testTvJumpYearForOrBack)
	t.Run("testRuns TvSelectMonth", testTvSelectMonth)
	t.Run("testRuns TvJumpToToday", testTvJumpToToday)

	//Tests zum Managen der ListView
	t.Run("testRuns LvSelectDate", testLvSelectDate)
	t.Run("testRuns LvSelectEntriesPerPage", testLvSelectEntriesPerPage)
	t.Run("testRuns LvJumpPageForward", testLvJumpPageForward)
	t.Run("testRuns LvJumpPageBack", testLvJumpPageBack)

	//Tests zum Managen der FilerView
	t.Run("testRuns FvSelectEntriesPerPage", testFvSelectEntriesPerPage)
	t.Run("testRuns FvJumpPageForward", testFvJumpPageForward)
	t.Run("testRuns FvJumpPageBack", testFvJumpPageBack)
	t.Run("testRuns FvFilter", testFvFilter)
}
