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
Hier Folgen die Tests zum Termine erstellen/bearbeiten/löschen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testCreateTermin(t *testing.T) {
	vm := new(ViewManager)
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
	//Länge des TerminCaches vor dem Hinzufügen des neuen Termins
	oldLen := len(vm.TerminCache)
	vm.CreateTermin(r, vm.Username)
	assert.Equal(t, oldLen+1, len(vm.TerminCache), "Die Länge sollte um eins erhöht worden sein.")
	assert.Equal(t, termin, vm.TerminCache[0], "Die Termine sollten überein stimmen.")

}
func testEditTermin(t *testing.T) {
	//r *http.Request, username string, monthEntries []dayInfos
}

func TestViewManager(t *testing.T) {
	//createTestTermins()
	t.Run("testRuns CreateTermin", testCreateTermin)
	t.Run("testRuns EditTermin", testEditTermin)
}
