package kalenderansicht

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
// generateRandomDate
//generiert ein komplett zufälliges Datum (bis zum Jahr 3000)
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
	fmt.Println(randomDate)
	return randomDate
}

// createWeeklyTestTermin
// generiert eines wöchentlichen/jährlichen/monatlichen Termins, damit diese im Kalender anzeigbar sind,
// um die Funktion der Navigation der Webseite zu testen
func createTestTermin(repeat ds.Repeat) *ViewManager {
	vm := InitViewManager("testuser")

	switch repeat {
	case ds.WEEKLY:
		newTermin := ds.CreateNewTermin("test Title", "test", ds.WEEKLY, createSpecificDate(2022, 2, 11), createSpecificDate(2023, 2, 11), vm.Username)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	case ds.YEARLY:
		newTermin := ds.CreateNewTermin("test Title", "test", ds.YEARLY, createSpecificDate(2020, 2, 11), createSpecificDate(2024, 2, 11), vm.Username)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	case ds.MONTHLY:
		newTermin := ds.CreateNewTermin("test Title", "test", ds.MONTHLY, createSpecificDate(2021, 2, 11), createSpecificDate(time.Now().Year(), 30, 12), vm.Username)
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
	assert.Equal(t, ds.Never, filterRepetition("1"))
	assert.Equal(t, ds.DAILY, filterRepetition("2"))
	assert.Equal(t, ds.WEEKLY, filterRepetition("3"))
	assert.Equal(t, ds.MONTHLY, filterRepetition("4"))
	assert.Equal(t, ds.YEARLY, filterRepetition("5"))
}

func testCreateTermin(t *testing.T) {
	vm := new(ViewManager)
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

	termin := ds.Termin{
		Title:       "Test Termin",
		Description: "Spaßiger Termin",
		Recurring:   ds.DAILY,
		Date:        createSpecificDate(2022, 11, 11),
		EndDate:     createSpecificDate(2030, 11, 11),
	}
	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)
	vm.CreateTermin(r, vm.Username)

	//Testen, ob ein Termin dem Cache hinzugefügt worden ist
	assert.Equal(t, oldLen+1, len(vm.TerminCache), "Die Länge sollte um eins erhöht worden sein.")
	//Testen, ob Termin im Cache dem neuen Termin entspricht
	assert.Equal(t, termin, vm.TerminCache[0], "Die Termine sollten überein stimmen.")

}

// testCreateTerminLogicCheck
// Wenn der Nutzer einen termin eingibt, dessen Ende zeitlich vor dem Startdatum ist, muss das Enddatum auf das Startdatum
// gesetzt werden
func testCreateTerminLogicCheck(t *testing.T) {
	vm := new(ViewManager)
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

	termin := ds.Termin{
		Title:       "Test Termin",
		Description: "Spaßiger Termin",
		Recurring:   ds.DAILY,
		Date:        createSpecificDate(2022, 11, 11),
		EndDate:     createSpecificDate(2022, 11, 11),
	}
	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)
	vm.CreateTermin(r, vm.Username)

	//Testen, ob ein Termin dem Cache hinzugefügt worden ist
	assert.Equal(t, oldLen+1, len(vm.TerminCache), "Die Länge sollte um eins erhöht worden sein.")
	//Testen, ob Termin im Cache dem neuen Termin entspricht
	assert.Equal(t, termin, vm.TerminCache[0], "Die Termine sollten überein stimmen.")

}
func testEditTerminDelete(t *testing.T) {
	vm := new(ViewManager)
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

	//Termin dem TerminCaches hinzufügen
	vm.CreateTermin(r, vm.Username)

	//Erstellen der Lösch-Request
	data = url.Values{}
	data.Add("oldTitle", "Test Termin")
	data.Add("editing", "Löschen: Termin")
	data.Add("title", "Test Termin")
	data.Add("description", "Spaßiger Termin")
	data.Add("repeat", "4") //Der Wert 1 entspricht der Wiederholung "niemals"
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2030-11-11")

	//Erstellen der Request
	r, _ = http.NewRequest("POST", "?termineBearbeiten", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)

	//Termin bearbeiten
	vm.EditTermin(r, vm.Username)

	//Testen, ob der Termin aus dem Cache gelöscht worden ist
	assert.Equal(t, oldLen-1, len(vm.TerminCache), "Die Länge sollte um 1 reduziert worden sein.")
}
func testEditTerminEdit(t *testing.T) {
	vm := new(ViewManager)
	//Erstellen der Termininfos, die über die Request gesendet werden
	data := url.Values{}
	data.Add("title", "Test Termin")
	data.Add("description", "Spaßiger Termin")
	data.Add("repeat", "täglich")
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2030-11-11")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "?terminErstellen", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Termin dem TerminCaches hinzufügen
	vm.CreateTermin(r, vm.Username)

	//Erstellen der Lösch-Request
	data = url.Values{}
	data.Add("oldTitle", "Test Termin")
	data.Add("editing", "1")
	data.Add("title", "Bearbeiteter Test Termin")
	data.Add("description", "Spaßiger bearbeiteter Termin")
	data.Add("repeat", "1") //Der Wert 1 entspricht der Wiederholung "niemals"
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2032-11-11")

	//Erstellen der Request
	r, _ = http.NewRequest("POST", "?termineBearbeiten", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//bearbeiteter Termin zum Überprüfen
	editetTermin := ds.Termin{
		Title:       "Bearbeiteter Test Termin",
		Description: "Spaßiger bearbeiteter Termin",
		Recurring:   ds.Never,
		Date:        createSpecificDate(2022, 11, 11),
		EndDate:     createSpecificDate(2032, 11, 11),
	}
	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)

	//Termin bearbeiten
	vm.EditTermin(r, vm.Username)

	//Testen, ob ein Termin dem Cache hinzugefügt worden ist
	assert.Equal(t, oldLen, len(vm.TerminCache), "Die Länge sollte gleich geblieben sein.")
	//Testen, ob Termin im Cache dem neuen Termin entspricht
	assert.Equal(t, editetTermin, vm.TerminCache[0], "Die Termine sollten überein stimmen.")
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

	//Der in der TableView angezeigte Monat ist November 2022
	vm.Tv.ShownDate = createSpecificDate(2022, 01, 11)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet am 2.11.2022 und findet wöchentlich statt
	//Monat ist nun November 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[9-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[16-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[23-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[30-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe einen Monat vor
	vm.TvJumpMonthFor()

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet am 2.11.2022 und findet wöchentlich statt
	//Monat ist nun Dezember 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[7-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[14-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[21-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[28-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}
func testTvJumpMonthBack(t *testing.T) {

	//Erstelle einen monatlichen Termin zum Testen: ab 02.11.2021
	vm := createTestTermin(ds.MONTHLY)

	//Der in der TableView angezeigte Monat ist November 2022
	vm.Tv.ShownDate = createSpecificDate(2022, 01, 11)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet am 2.11.2021cund findet monatlich statt
	//Monat ist nun November 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe einen Monat vor
	vm.TvJumpMonthBack()

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	//Der Termin startet ab 02.11.2021 und findet moantlich statt
	//Monat ist nun Oktober 2022
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")
	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

func testTvJumpYearForOrBack(t *testing.T) {
	//Jährlichen testtermin erstellen 2.11.2020-  2.11.2024
	vm := createTestTermin(ds.YEARLY)

	//Der in der TableView angezeigte Monat ist November 2022
	vm.Tv.ShownDate = createSpecificDate(2022, 01, 11)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe ein Jahr vor (zu 2023)
	vm.TvJumpYearForOrBack(1)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe ein Jahr zurück (zu 2022)
	vm.TvJumpYearForOrBack(1)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

func testTvSelectMonth(t *testing.T) {
	//Monatlichen testtermin erstellen 2.11.2021 - 30.12 des heutigen Jahres
	vm := createTestTermin(ds.MONTHLY)

	//Der in der TableView angezeigte Monat ist November 2022
	vm.Tv.ShownDate = createSpecificDate(2022, 01, 11)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe zu einem beliebigen Monat des Jahres 2022
	month := time.Month(rand.Intn(13))
	vm.TvSelectMonth(month)

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

func testTvJumpToToday(t *testing.T) {
	//Monatlichen testtermin erstellen: 2.11.2021 - 30.12. des heutigen Jahres
	vm := createTestTermin(ds.MONTHLY)

	//Der in der TableView angezeigte Monat wird auf einen Tag in der Vergangenheit gesetzt (-> 1.1.2022)
	vm.Tv.ShownDate = createSpecificDate(2022, 01, 01)
	vm.Tv.CreateTerminTableEntries(vm.TerminCache)

	//Kontrollieren, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Springe zu heute
	vm.TvJumpToToday()

	//Testen, ob der Termin in den richtigen stellen in dem Slice MonthEntries von der tableView hinzugefügt worden ist
	assert.Equal(t, vm.TerminCache[0], vm.Tv.MonthEntries[2-1].Dayentries[0], "Die Termine sollten übereinstimmen.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
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
	entriesPerPage := 1
	vm.LvSelectEntriesPerPage(entriesPerPage)
	assert.Equal(t, entriesPerPage*5, vm.Lv.EntriesPerPage, "Die Anzahl der Einträge pro Seite sollte 5 sein.")
}

func testLvJumpPageForward(t *testing.T) {
	vm := createTestTermin(ds.WEEKLY)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.LvJumpPageForward()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1 sein, da es nur einen Eintrag gibt.")

	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < 30; i++ {
		newTermin := ds.CreateNewTermin("test Title"+fmt.Sprint(i), "test", ds.WEEKLY, createSpecificDate(2022, 2, 11), createSpecificDate(2023, 2, 11), vm.Username)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	}

	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.LvJumpPageForward()
	assert.Equal(t, 2, vm.Lv.CurrentPage, "Die Seite sollte 2 sein.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
}

func testLvJumpPageBack(t *testing.T) {
	vm := createTestTermin(ds.WEEKLY)
	vm.Lv.CreateTerminListEntries(vm.TerminCache)
	vm.Lv.JumpPageForward()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1 sein, da die Seitennummer nicht kleiner als 1 sein kann..")

	//dem Cache mehrere Termine hinzufügen
	for i := 0; i < 30; i++ {
		newTermin := ds.CreateNewTermin("test Title"+fmt.Sprint(i), "test", ds.WEEKLY, createSpecificDate(2022, 2, 11), createSpecificDate(2023, 2, 11), vm.Username)
		vm.TerminCache = ds.AddToCache(newTermin, vm.TerminCache)
	}
	vm.Lv.CreateTerminListEntries(vm.TerminCache)

	//Aktuelle Seite auf Seite 2 setzten
	vm.Lv.CurrentPage = 2
	vm.LvJumpPageBack()
	assert.Equal(t, 1, vm.Lv.CurrentPage, "Die Seite sollte 1sein.")

	//Löschen der erstellten Testdaten
	vm.TerminCache = ds.DeleteAll(vm.TerminCache, vm.Username)
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
	t.Run("testRuns EditTermin-Delete", testEditTerminDelete)
	t.Run("testRuns EditTermin-Edit", testEditTerminEdit)

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
}
