package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

/*
**************************************************************************************************************
Funktionen zum zufälligen generieren von Testdaten
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func ganerateRandomDate() time.Time {
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

/*
**************************************************************************************************************
Test zur Datenbeschaffung für die Tabellenansicht (um Jahr und Monat anzuzeigen)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testShowYear(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	jahr := ta.ShownDate.Year()
	assert.Equal(t, jahr, ta.ShowYear(), "Die Jahre sollten identisch sein.")
}

func testShowMonth(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	monat := ta.ShownDate.Month()
	assert.Equal(t, monat, ta.ShowMonth(), "Die Monate sollten identisch sein.")
}

/*
**************************************************************************************************************
Test zur Navigation innerhalb der Tabellenansicht
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testJumpMonthBack(t *testing.T) {
	i := 100
	var ta TableView
	for i > 0 {
		i -= 1
		ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
		aktuellerMonat := (int(ta.ShownDate.Month()))
		ta.JumpMonthBack()
		neuerMonat := (int(ta.ShownDate.Month()))
		if aktuellerMonat < 12 {
			assert.Equal(t, aktuellerMonat+1, neuerMonat)
		} else {
			assert.Equal(t, 1, neuerMonat)
		}
	}
}
func testJumpMonthFor(t *testing.T) {
	i := 100
	var ta TableView
	for i > 0 {
		i -= 1
		ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
		aktuellerMonat := (int(ta.ShownDate.Month()))
		ta.JumpMonthFor()
		neuerMonat := (int(ta.ShownDate.Month()))
		if aktuellerMonat == 1 {
			assert.Equal(t, 12, neuerMonat, "Die Monate sollten identisch sein.")
		} else {
			assert.Equal(t, aktuellerMonat-1, neuerMonat, "Die Monate sollten identisch sein.")
		}
	}
}
func testJumpToYear(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	aktuellesJahr := ta.ShownDate.Year()
	ta.JumpToYear(1)
	assert.Equal(t, aktuellesJahr+1, ta.ShownDate.Year(), "Die Jahre sollten identisch sein.")
	ta.JumpToYear(-1)
	assert.Equal(t, aktuellesJahr, ta.ShownDate.Year(), "Die Jahre sollten identisch sein.")
}

func testSelectMonth(t *testing.T) {
	i := 100
	var ta TableView
	for i > 0 {
		i -= 1
		ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
		gewaehlterMonat := time.Month(rand.Intn(12-1) + 1)
		ta.SelectMonth(gewaehlterMonat)
		assert.Equal(t, ta.ShownDate.Month(), gewaehlterMonat, "Die Monate sollten identisch sein.")
	}
}

func testJumpToToday(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	ta.JumpToToday()
	assert.Equal(t, ta.ShownDate, ta.getFirstDayOfMonth(time.Now()), "Das Datum sollte das heutige sein.")
}

/*
**************************************************************************************************************
Tests zum Filtern der Termine
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

// testFilterCalendarEntriesCorrectAppending
// Testen ob Daten sortiert am richtigen Tag eingefügt werden
func testFilterCalendarEntriesCorrectAppending(t *testing.T) {
	var ta TableView
	//Hier wird August als Monat von Interesse gesetzt
	ta.ShownDate = createTestDate(2022, 1, 8)

	//Datum für Testtermine erstellen
	testTermin1 := createTestDate(2022, 15, 8)
	testTermin2 := createTestDate(2022, 10, 8)
	testTermin3 := createTestDate(2022, 20, 8)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 3)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.Repeat(ds.Never), testTermin1, testTermin1)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.Never), testTermin2, testTermin2)
	testTermine[2] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.Never), testTermin3, testTermin3)

	//Testen ob Daten sortiert am richtigen Tag eingefügt werden
	filteresSlice := ta.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteresSlice[15-1].Dayentries[0])
	assert.Equal(t, testTermine[1], filteresSlice[10-1].Dayentries[0])
	assert.Equal(t, testTermine[2], filteresSlice[20-1].Dayentries[0])
}

// testFilterCalendarEntriesCorrectRecurring
// Testen ob wöchentlich/Jährlich/moatlich auftetende Termine erkannt und richtig eingefügt werden
// Testen von Terminen, die am selben Tag sind
func testFilterCalendarEntriesCorrectRecurring(t *testing.T) {
	var ta TableView
	//Hier wird November 2022 als Monat von Interesse gesetzt
	ta.ShownDate = createTestDate(2022, 1, 11)

	//Datum für Testtermine erstellen
	testTermin1Starts := createTestDate(2022, 10, 10)
	testTermin1Ends := createTestDate(2023, 10, 10)

	testTermin2Starts := createTestDate(2022, 27, 10)
	testTermin2Ends := createTestDate(2023, 27, 10)

	testTermin3Starts := createTestDate(2021, 9, 11)
	testTermin3Ends := createTestDate(2023, 9, 11)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 3)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.Repeat(ds.MONTHLY), testTermin1Starts, testTermin1Ends)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.WEEKLY), testTermin2Starts, testTermin2Ends)
	testTermine[2] = ds.NewTerminObj("testTermin2", "test", ds.Repeat(ds.YEARLY), testTermin3Starts, testTermin3Ends)

	//Testen ob Daten sortiert am richtigen Tag eingefügt werden
	filteredSlice := ta.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteredSlice[10-1].Dayentries[0])
	assert.Equal(t, testTermine[1], filteredSlice[10-1].Dayentries[1])
	assert.Equal(t, testTermine[1], filteredSlice[17-1].Dayentries[0])
	assert.Equal(t, testTermine[1], filteredSlice[24-1].Dayentries[0])
	assert.Equal(t, testTermine[2], filteredSlice[9-1].Dayentries[0])
}
func TestTabellenAnsicht(t *testing.T) {
	//Tests zur Kontrolle der richtigen Wertanzeige in der Webansicht
	t.Run("testRuns", testShowYear)
	t.Run("testRuns", testShowMonth)

	//Tests zum Navigieren in der Tabellenansicht
	t.Run("testRuns", testJumpMonthBack)
	t.Run("testRuns", testJumpMonthFor)
	t.Run("testRuns", testJumpToYear)
	t.Run("testRuns", testSelectMonth)
	t.Run("testRuns", testJumpToToday)

	//Tests zum Termine filtern
	t.Run("testRuns", testFilterCalendarEntriesCorrectAppending)
	t.Run("testRuns", testFilterCalendarEntriesCorrectRecurring)
}
