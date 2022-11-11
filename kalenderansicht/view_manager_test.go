package kalenderansicht

import (
	"DHBW_GO_Projekt/dateisystem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

import ds "DHBW_GO_Projekt/dateisystem"

/*
**************************************************************************************************************
Funktionen zum zufälligen generieren von Testdaten
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
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
func createTestTermins() {
	var k []ds.Termin
	i := 0
	for i < 50 {
		i += 1
		k = ds.AddToCache(ds.NewTerminObj("testa"+strconv.Itoa(i), "test", dateisystem.Repeat(ds.Never), generateRandomDateInSpecificMonth(2022, 11), generateRandomDateInSpecificMonth(2022, 11)), k)
		k = ds.AddToCache(ds.NewTerminObj("testb"+strconv.Itoa(i), "test", dateisystem.Repeat(ds.YEARLY), generateRandomDateInSpecificMonth(2022, 10), generateRandomDateInSpecificMonth(2022, 10)), k)
		k = ds.AddToCache(ds.NewTerminObj("testc"+strconv.Itoa(i), "test", dateisystem.Repeat(ds.WEEKLY), generateRandomDateInSpecificMonth(2022, 12), generateRandomDateInSpecificMonth(2022, 12)), k)
		ds.StoreCache(k, "mik")
	}
}
func createTestDate(year, day int, month time.Month) time.Time {
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
Hier Folgen die Tests
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testCreateTermin(t *testing.T) {
	//Erstellen der Termininfos, die über die Request gesendet werden
	data := url.Values{}
	data.Add("title", "Test Termin")
	data.Add("description", "Spaßiger Termin")
	data.Add("repeat", "täglich")
	data.Add("date", "2022-11-11")
	data.Add("endDate", "2030-11-11")

	//Erstellen der Request
	r, _ := http.NewRequest("POST", "/tabellenAnsicht?terminErstellen", strings.NewReader(data.Encode()))
	r.Header.Add("", "")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	termin := ds.Termin{
		Title:       "Test Termin",
		Description: "Spaßiger Termin",
		Recurring:   ds.DAILY,
		Date:        createTestDate(2022, 11, 11),
		EndDate:     createTestDate(2030, 11, 11),
	}
	assert.Equal(t, termin, CreateTermin(r, "Testuser"), "Die Termine sollten identisch sein.")

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
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", dateisystem.Repeat(ds.Never), testTermin1, testTermin1)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", dateisystem.Repeat(ds.Never), testTermin2, testTermin2)
	testTermine[2] = ds.NewTerminObj("testTermin2", "test", dateisystem.Repeat(ds.Never), testTermin3, testTermin3)

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
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", dateisystem.Repeat(ds.MONTHLY), testTermin1Starts, testTermin1Ends)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", dateisystem.Repeat(ds.WEEKLY), testTermin2Starts, testTermin2Ends)
	testTermine[2] = ds.NewTerminObj("testTermin2", "test", dateisystem.Repeat(ds.YEARLY), testTermin3Starts, testTermin3Ends)

	//Testen ob Daten sortiert am richtigen Tag eingefügt werden
	filteredSlice := ta.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteredSlice[10-1].Dayentries[0])
	assert.Equal(t, testTermine[1], filteredSlice[10-1].Dayentries[1])
	assert.Equal(t, testTermine[1], filteredSlice[17-1].Dayentries[0])
	assert.Equal(t, testTermine[1], filteredSlice[24-1].Dayentries[0])
	assert.Equal(t, testTermine[2], filteredSlice[9-1].Dayentries[0])
}

func TestCalendarView(t *testing.T) {
	//createTestTermins()
	t.Run("testRuns", testFilterCalendarEntriesCorrectAppending)
	t.Run("testRuns", testFilterCalendarEntriesCorrectRecurring)
	t.Run("testRuns", testCreateTermin)
}
