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
Test zur Datenbeschaffung für die Tabellenansicht (um Jahr und Monat anzuzeigen)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testShowYear(t *testing.T) {
	var tv TableView
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
	jahr := tv.ShownDate.Year()
	assert.Equal(t, jahr, tv.ShowYear(), "Die Jahre sollten identisch sein.")
}

func testShowMonth(t *testing.T) {
	var tv TableView
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
	monat := tv.ShownDate.Month()
	assert.Equal(t, monat, tv.ShowMonth(), "Die Monate sollten identisch sein.")
}

/*
**************************************************************************************************************
Test zur Navigation innerhalb der Tabellenansicht
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testJumpMonthFor(t *testing.T) {
	i := 100
	var tv TableView
	for i > 0 {
		i -= 1
		tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
		aktuellerMonat := (int(tv.ShownDate.Month()))
		tv.JumpMonthFor()
		neuerMonat := (int(tv.ShownDate.Month()))
		if aktuellerMonat < 12 {
			assert.Equal(t, aktuellerMonat+1, neuerMonat, "Die Monate müssen identisch sein")
		} else {
			assert.Equal(t, 1, neuerMonat, "Die Monate müssen identisch sein")
		}
	}
}
func testJumpMonthBack(t *testing.T) {
	i := 100
	var tv TableView
	for i > 0 {
		i -= 1
		tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
		aktuellerMonat := (int(tv.ShownDate.Month()))
		tv.JumpMonthBack()
		neuerMonat := (int(tv.ShownDate.Month()))
		if aktuellerMonat == 1 {
			assert.Equal(t, 12, neuerMonat, "Die Monate sollten identisch sein.")
		} else {
			assert.Equal(t, aktuellerMonat-1, neuerMonat, "Die Monate sollten identisch sein.")
		}
	}
}
func testJumpToYear(t *testing.T) {
	var tv TableView
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
	aktuellesJahr := tv.ShownDate.Year()
	tv.JumpYearForOrBack(1)
	assert.Equal(t, aktuellesJahr+1, tv.ShownDate.Year(), "Die Jahre sollten identisch sein.")
	tv.JumpYearForOrBack(-1)
	assert.Equal(t, aktuellesJahr, tv.ShownDate.Year(), "Die Jahre sollten identisch sein.")
}

func testSelectMonth(t *testing.T) {
	i := 100
	var tv TableView
	for i > 0 {
		i -= 1
		tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDate())
		gewaehlterMonat := time.Month(rand.Intn(12-1) + 1)
		tv.SelectMonth(gewaehlterMonat)
		assert.Equal(t, tv.ShownDate.Month(), gewaehlterMonat, "Die Monate sollten identisch sein.")
	}
}

func testJumpToToday(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(generateRandomDate())
	ta.JumpToToday()
	assert.Equal(t, ta.ShownDate, ta.getFirstDayOfMonth(time.Now()), "Das Datum sollte das heutige sein.")
}

/*
**************************************************************************************************************
Tests zum Filtern der Termine und Hilfsfunktionen zur Anzeige auf der Webseite
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testGetFirstDayOfMonth(t *testing.T) {
	var ta TableView
	ta.ShownDate = generateRandomDate()
	firstDay := ta.getFirstDayOfMonth(ta.ShownDate)
	assert.Equal(t, 1, firstDay.Day(), "Der Tag sollte der Erste sein")
}

func testGetLastDayOfMonth(t *testing.T) {
	var ta TableView

	//Test Monate mit 31 Tagen
	months31 := []int{1, 3, 5, 7, 8, 10, 12}
	ta.ShownDate = generateRandomDateInSpecificMonth(2022, time.Month(months31[rand.Intn(6)]))
	lastDay := ta.getLastDayOfMonth()
	assert.Equal(t, 31, lastDay.Day(), "Der Tag sollte der 31. sein")

	//Test Monate mit 30 Tagen
	months30 := []int{4, 6, 9, 11}
	ta.ShownDate = generateRandomDateInSpecificMonth(2022, time.Month(months30[rand.Intn(3)]))
	lastDay = ta.getLastDayOfMonth()
	assert.Equal(t, 30, lastDay.Day(), "Der Tag sollte der 30. sein")

	//Test Februar
	ta.ShownDate = generateRandomDateInSpecificMonth(2022, 2)
	lastDay = ta.getLastDayOfMonth()
	assert.Equal(t, 28, lastDay.Day(), "Der Tag sollte der 28. sein")

	//Test Februar für Schaltjahr
	ta.ShownDate = generateRandomDateInSpecificMonth(2020, 2)
	lastDay = ta.getLastDayOfMonth()
	assert.Equal(t, 29, lastDay.Day(), "Der Tag sollte der 29. sein")
}
func testMonthStarts(t *testing.T) {
	var tv TableView
	//1. Test Monat: November 2022
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDateInSpecificMonth(2022, 11))
	assert.Equal(t, 1, len(tv.MonthStarts()), "Die Slice sollte die Länge 1 haben, da der 11.2022 an einem Dienstag startet.")

	//2. Test Monat: Dezember 2022
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDateInSpecificMonth(2022, 12))
	assert.Equal(t, 3, len(tv.MonthStarts()), "Die Slice sollte die Länge 3 haben, da der 12.2022 an einem Donnerstag startet.")

	//3. Test Monat: Januar 2023
	tv.ShownDate = tv.getFirstDayOfMonth(generateRandomDateInSpecificMonth(2023, 1))
	assert.Equal(t, 6, len(tv.MonthStarts()), "Die Slice sollte die Länge 6 haben, da der 1.2023 an einem Sonntag startet.")
}
func testIsToday(t *testing.T) {
	dateToday := time.Now()
	assert.Equal(t, true, IsToday(dateToday), "Die Funktion sollte true zurückgeben.")

	dateToday = time.Now().AddDate(-1, 0, 0)
	assert.Equal(t, false, IsToday(dateToday), "Die Funktion sollte false zurückgeben.")
}

func testNeedsBreak(t *testing.T) {

	testDate := createSpecificDate(2022, 21, 11)
	assert.False(t, NeedsBreak(testDate), "Das Ergebnis sollte false sein, da der Tag kein Sonntag ist.")

	testDate = createSpecificDate(2022, 20, 11)
	assert.True(t, NeedsBreak(testDate), "Das Ergebnis sollte true sein, da der Tag ein Sonntag ist.")
}

// testFilterCalendarEntriesCorrectAppending
// Testen ob Daten sortiert am richtigen Tag eingefügt werden
func testFilterCalendarEntriesCorrectAppending(t *testing.T) {
	var tv TableView
	//Hier wird August als Monat von Interesse gesetzt
	tv.ShownDate = createSpecificDate(2022, 1, 8)

	//Datum für Testtermine erstellen
	testTermin1 := createSpecificDate(2022, 15, 8)
	testTermin2 := createSpecificDate(2022, 10, 8)
	testTermin3 := createSpecificDate(2022, 20, 8)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 0, 3)
	testTermine = append(testTermine, ds.NewTerminObj("testTermin1", "test", ds.Never, testTermin1, testTermin1))
	testTermine = append(testTermine, ds.NewTerminObj("testTermin2", "test", ds.Never, testTermin2, testTermin2))
	testTermine = append(testTermine, ds.NewTerminObj("testTermin3", "test", ds.Never, testTermin3, testTermin3))

	//Testen ob Daten sortiert am richtigen Tag eingefügt werden
	filteresSlice := tv.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteresSlice[15-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteresSlice[10-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[2], filteresSlice[20-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
}

// testFilterCalendarEntriesCorrectRecurring
// Testen ob wöchentlich/Jährlich/moatlich auftetende Termine erkannt und richtig eingefügt werden
// Testen von Terminen, die am selben Tag sind
func testFilterCalendarEntriesCorrectRecurring(t *testing.T) {
	var tv TableView
	//Hier wird November 2022 als Monat von Interesse gesetzt
	tv.ShownDate = createSpecificDate(2022, 1, 11)

	//Datum für Testtermine erstellen
	testTermin1Starts := createSpecificDate(2022, 9, 11)
	testTermin1Ends := createSpecificDate(2022, 29, 11)

	testTermin2Starts := createSpecificDate(2022, 27, 10)
	testTermin2Ends := createSpecificDate(2023, 27, 10)

	testTermin3Starts := createSpecificDate(2021, 9, 11)
	testTermin3Ends := createSpecificDate(2023, 9, 11)

	testTermin4Starts := createSpecificDate(2022, 10, 10)
	testTermin4Ends := createSpecificDate(2023, 10, 10)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 4)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends)
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.MONTHLY, testTermin4Starts, testTermin4Ends)

	//Testen ob Daten sortiert am richtigen Tag eingefügt werden
	filteredSlice := tv.FilterCalendarEntries(testTermine)

	//Kontrolle ob Termine täglich eingefügt werden, Zeitraum 9-29.11
	for day := testTermin4Starts.Day(); day <= testTermin4Ends.Day(); day++ {
		assert.Equal(t, testTermine[0], filteredSlice[day-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
		assert.Equal(t, testTermine[0], filteredSlice[day-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
		assert.Equal(t, testTermine[0], filteredSlice[day-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
		assert.Equal(t, testTermine[0], filteredSlice[29-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
		assert.Equal(t, 0, len(filteredSlice[30-1].Dayentries), "An dieser Position sollte kein Termin eingetragen sein")
	}
	//Kontrolle ob Termine korrekt wöchentlich eingefügt werden, startet beim 27.10 -> 3.11 -> 10.11...
	//Kontrolle von korrektem Einfügen hinter schon eingefügten Terminen
	assert.Equal(t, testTermine[1], filteredSlice[3-1].Dayentries[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteredSlice[10-1].Dayentries[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteredSlice[17-1].Dayentries[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteredSlice[24-1].Dayentries[1], "Termine an diesen Positionen sollten identisch sein")

	//Kontrolle ob Termine korrekt jährlich eingefügt werden
	//Kontrolle von korrektem Einfügen hinter schon eingefügten Terminen
	assert.Equal(t, testTermine[2], filteredSlice[9-1].Dayentries[1], "Termine an diesen Positionen sollten identisch sein")

	//Kontrolle ob Termine korrekt monatlich eingefügt werden
	//Kontrolle von korrektem Einfügen hinter schon eingefügten Terminen
	assert.Equal(t, testTermine[3], filteredSlice[10-1].Dayentries[2], "Termine an diesen Positionen sollten identisch sein")

}
func TestTableView(t *testing.T) {
	//Tests zur Kontrolle der richtigen Wertanzeige in der Webansicht
	t.Run("testRuns ShowYear", testShowYear)
	t.Run("testRuns showMonth", testShowMonth)

	//Tests zum Navigieren in der Tabellenansicht
	t.Run("testRuns JumpMonthBack", testJumpMonthBack)
	t.Run("testRuns JumpMonthFor", testJumpMonthFor)
	t.Run("testRuns JumpYearForOrBack", testJumpToYear)
	t.Run("testRuns SelectMonth", testSelectMonth)
	t.Run("testRuns JumpToToday", testJumpToToday)

	//Tests zum Termine filtern
	t.Run("testRuns FilterCalendarEntriesCorrectAppending", testFilterCalendarEntriesCorrectAppending)
	t.Run("testRuns FilterCalendarEntriesCorrectRecurring", testFilterCalendarEntriesCorrectRecurring)

	//Hilfsfunktionen zum Termine filtern und Termine richtig in Webseite anzuzeigen
	t.Run("testRuns GetFirstDayOfMonth", testGetFirstDayOfMonth)
	t.Run("testRuns GetLastDayOfMonth", testGetLastDayOfMonth)
	t.Run("testRuns MonthStarts", testMonthStarts)
	t.Run("testRuns IsToday", testIsToday)
	t.Run("testRuns NeedsBreak", testNeedsBreak)
}
