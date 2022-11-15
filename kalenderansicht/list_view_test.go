package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

/*
**************************************************************************************************************
Tests für Custom-Settings innerhalb der Webseite
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testSelectDate(t *testing.T) {
	lv := InitListView()
	newDate := time.Date(2030, 11, 11, 0, 0, 0, 0, time.UTC)

	//Erstellen des Datums als POST-Value
	data := url.Values{}
	data.Add("selDate", "2030-11-11")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "/listenAnsicht?selDatum=Datum\"", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	lv.SelectDate(r)
	assert.Equal(t, newDate, lv.SelectedDate, "Die zwei Daten sollten identisch sein.")
}

func testSelectEntriesPerPage(t *testing.T) {
	lv := InitListView()
	//Die Values der Webseite gehen von 1 bis 3, wobei 1 für 5 Einträge steht und 3 für 15
	entriesAmountValue := 2
	lv.SelectEntriesPerPage(2)
	assert.Equal(t, entriesAmountValue*5, lv.EntriesPerPage, "Die Nummern sollten identisch sein.")
}

func testFilterCalendarEntries(t *testing.T) {
	var lv ListView
	//Hier wird der 1.11.2022 als Startdatum gesetzt
	lv.SelectedDate = createTestDate(2022, 1, 11)

	//Datum für Testtermine erstellen
	testTermin1Starts := createTestDate(2022, 9, 11)
	testTermin1Ends := createTestDate(2022, 29, 11)

	testTermin2Starts := createTestDate(2022, 27, 10)
	testTermin2Ends := createTestDate(2023, 27, 10)

	testTermin3Starts := createTestDate(2021, 9, 11)
	testTermin3Ends := createTestDate(2023, 9, 11)

	testTermin4Starts := createTestDate(2022, 10, 10)
	testTermin4Ends := createTestDate(2023, 10, 10)

	testTermin5Starts := createTestDate(2021, 10, 8)
	testTermin5Ends := createTestDate(2022, 10, 9)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.Repeat(ds.DAILY), testTermin1Starts, testTermin1Ends)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.WEEKLY), testTermin2Starts, testTermin2Ends)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.Repeat(ds.YEARLY), testTermin3Starts, testTermin3Ends)
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.Repeat(ds.MONTHLY), testTermin4Starts, testTermin4Ends)
	//Termin sollte nicht hinzugefügt werden
	testTermine[4] = ds.NewTerminObj("testTermin5", "test", ds.Repeat(ds.MONTHLY), testTermin5Starts, testTermin5Ends)

	//Testen ob Daten in Slice eingefügt wurden
	filteredSlice := lv.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteredSlice[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteredSlice[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[2], filteredSlice[2], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[3], filteredSlice[3], "Termine an diesen Positionen sollten identisch sein")
	assert.True(t, len(filteredSlice) == 4, "4 Termine sollten sich in der Slice befinden")

}

/*
**************************************************************************************************************
Test zur Navigation innerhalb der Listenansicht
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testSpringSeiteWeiter(t *testing.T) {
	//lv := InitListView()

}
func testSpringSeiteZurueck(t *testing.T) {
	//lv := InitListView()
}

func TestListView(t *testing.T) {
	t.Run("testRuns", testSelectDate)
	t.Run("testRuns", testSelectEntriesPerPage)
	t.Run("testRuns", testSpringSeiteWeiter)
	t.Run("testRuns", testSpringSeiteZurueck)
	t.Run("testRuns", testFilterCalendarEntries)
}
