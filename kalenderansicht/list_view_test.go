package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
**************************************************************************************************************
Tests für Custom-Settings innerhalb der Webseite
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testSelectDate(t *testing.T) {
	var lv = new(ListView)
	newDate := time.Date(2030, 11, 11, 0, 0, 0, 0, time.UTC)
	lv.SelectDate(newDate)
	assert.Equal(t, newDate, lv.SelectedDate, "Die zwei Daten sollten identisch sein.")
}

func testSelectEntriesPerPage(t *testing.T) {
	var lv = new(ListView)
	//Die Values der Webseite gehen von 1 bis 3, wobei 1 für 5 Einträge steht und 3 für 15
	entriesAmountValue := 5
	lv.SelectEntriesPerPage(entriesAmountValue)
	assert.Equal(t, entriesAmountValue, lv.EntriesPerPage, "Die Nummern sollten identisch sein.")
}

/*
**************************************************************************************************************
Test zur Navigation innerhalb der Listenansicht
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testJumPageFor(t *testing.T) {
	//lv := InitListView()

}
func testJumpPageback(t *testing.T) {
	//lv := InitListView()
}

/*
*********************************************************************************************************************
Tests zum Filtern und Infos für die Darstellung der Termine in der Listenansicht
*********************************************************************************************************************
*/
func testFilterCalendarEntries(t *testing.T) {
	var lv ListView
	//Hier wird der 1.11.2022 als Startdatum gesetzt
	lv.SelectedDate = createSpecificDate(2022, 1, 11)

	//Datum für Testtermine erstellen
	testTermin1Starts := createSpecificDate(2022, 9, 11)
	testTermin1Ends := createSpecificDate(2022, 29, 11)

	testTermin2Starts := createSpecificDate(2022, 27, 10)
	testTermin2Ends := createSpecificDate(2023, 27, 10)

	testTermin3Starts := createSpecificDate(2021, 9, 11)
	testTermin3Ends := createSpecificDate(2023, 9, 11)

	testTermin4Starts := createSpecificDate(2022, 10, 10)
	testTermin4Ends := createSpecificDate(2023, 10, 10)

	testTermin5Starts := createSpecificDate(2021, 10, 8)
	testTermin5Ends := createSpecificDate(2022, 10, 9)

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
	assert.Equal(t, 4, len(filteredSlice), "4 Termine sollten sich in der Slice befinden")

}
func testNextOccurrences(t *testing.T) {
	var lv ListView
	//Hier wird der 1.11.2022 als Startdatum gesetzt
	lv.SelectedDate = createSpecificDate(2022, 1, 11)

	//Slice mit den zu erwartenden Ergebnissen
	expectedResults := make([][]time.Time, 5)
	//Datum für Testtermine erstellen

	//täglicher Termin
	testTermin1Starts := createSpecificDate(2022, 9, 11)
	testTermin1Ends := createSpecificDate(2022, 29, 11)
	expectedResults[0] = []time.Time{createSpecificDate(2022, 9, 11), createSpecificDate(2022, 10, 11), createSpecificDate(2022, 11, 11)}

	//Wöchentlicher Termin
	testTermin2Starts := createSpecificDate(2022, 27, 10)
	testTermin2Ends := createSpecificDate(2023, 27, 10)
	expectedResults[1] = []time.Time{createSpecificDate(2022, 3, 11), createSpecificDate(2022, 10, 11), createSpecificDate(2022, 17, 11)}

	//Jährlicher Termin
	testTermin3Starts := createSpecificDate(2021, 9, 11)
	testTermin3Ends := createSpecificDate(2023, 9, 11)
	expectedResults[2] = []time.Time{createSpecificDate(2022, 9, 11), createSpecificDate(2023, 9, 11)} //Nur noch zwei vorkommende Termine

	//Monatlicher Termin
	testTermin4Starts := createSpecificDate(2022, 10, 10)
	testTermin4Ends := createSpecificDate(2023, 10, 10)
	expectedResults[3] = []time.Time{createSpecificDate(2022, 10, 11), createSpecificDate(2022, 10, 12), createSpecificDate(2023, 10, 1)}

	//sich nie wiederholender Termin
	testTermin5Starts := createSpecificDate(2022, 5, 11)
	testTermin5Ends := createSpecificDate(2022, 5, 11)
	expectedResults[4] = []time.Time{createSpecificDate(2022, 5, 11)} //einmaliger Termin

	//Slice mit Testterminen erstellen, jeder Wiederholungstyp dabei
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.Repeat(ds.DAILY), testTermin1Starts, testTermin1Ends)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.WEEKLY), testTermin2Starts, testTermin2Ends)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.Repeat(ds.YEARLY), testTermin3Starts, testTermin3Ends)
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.Repeat(ds.MONTHLY), testTermin4Starts, testTermin4Ends)
	testTermine[4] = ds.NewTerminObj("testTermin5", "test", ds.Repeat(ds.Never), testTermin5Starts, testTermin5Ends)

	assert.Equal(t, expectedResults[0], lv.NextOccurrences(testTermine[0]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[1], lv.NextOccurrences(testTermine[1]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[2], lv.NextOccurrences(testTermine[2]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[3], lv.NextOccurrences(testTermine[3]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[4], lv.NextOccurrences(testTermine[4]), "Daten-in den Slices sollten identisch sein")

}

/*
*********************************************************************************************************************
Führe alle Tests aus
*********************************************************************************************************************
*/
func TestListView(t *testing.T) {
	//Tests für Custom-Settings innerhalb der Webseite
	t.Run("testRuns SelectDate", testSelectDate)
	t.Run("testRuns SelectEntriesPerPage", testSelectEntriesPerPage)

	//Test zur Navigation innerhalb der Listenansicht
	t.Run("testRuns JumPageFor", testJumPageFor)
	t.Run("testRuns JumpPageback", testJumpPageback)

	//Tests zum Filtern und tests für Infos für die Darstellung der Termine in der Listenansicht
	t.Run("testRuns FilterCalendarEntries", testFilterCalendarEntries)
	t.Run("testRuns  NextOccurrences", testNextOccurrences)
}
