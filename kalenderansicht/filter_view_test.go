/*
@author: 2447899 8689159 3000685
*/
package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
*********************************************************************************************************************
Tests für Funktionen, die den Benutzer Custom-Settings & Navigation innerhalb der Webseite ermöglichen.
(Bsp.: Seitenanzahl festlegen, Seite weiter navigieren...)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testFilterSelectEntriesPerPage(t *testing.T) {
	var fv = new(FilterView)
	entriesAmountValue := 5
	fv.SelectEntriesPerPage(entriesAmountValue)
	assert.Equal(t, entriesAmountValue, fv.EntriesPerPage, "Die Nummern sollten identisch sein.")
	assert.Equal(t, 1, fv.CurrentPage, "Die aktuelle Seite sollte Seite 1 sein.")
}

func testGetEntries(t *testing.T) {
	//Filter-Ansicht erstellen
	fv := new(FilterView)
	//Hier wird das Startdatum gesetzt, abhängig vom heutigen Datum -> Testausführung unabhängig vom eig. Datum
	fv.EntriesPerPage = 5
	fv.CurrentPage = 1

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	fv.CreateTerminFilterEntries(testTermine30)

	assert.Equal(t, 30, len(fv.FilteredTermins), "Es sollten 30 Termine in dem Slice sein")

	//5 Einträge pro Tag
	fv.SelectEntriesPerPage(5)
	sliceWithEntries := fv.GetEntries()
	assert.Equal(t, 5, len(sliceWithEntries), "Es sollten 5 Termine in dem Slice sein")

	//10 Einträge pro Tag
	fv.SelectEntriesPerPage(10)
	sliceWithEntries = fv.GetEntries()
	assert.Equal(t, 10, len(sliceWithEntries), "Es sollten 10 Termine in dem Slice sein")

	//15 Einträge pro Tag
	fv.SelectEntriesPerPage(15)
	sliceWithEntries = fv.GetEntries()
	assert.Equal(t, 15, len(sliceWithEntries), "Es sollten 15 Termine in dem Slice sein")
}

// testJumPageForLastPage
// Hier wird der Fall abgedeckt, dass der Benutzer sich bereits auf der letzten, zur Darstellung der Termine benötigten Seite
// befindet und dennoch eine Seite vorspringen möchte
func testFilterJumPageForwardLastPage(t *testing.T) {
	//Liste erstellen, angezeigte Seite ist 1, keine Einträge vorhanden -> nur eine Seite benötigt
	lv := InitListView([]ds.Termin{})

	lv.JumpPageForward()
	assert.Equal(t, 1, lv.CurrentPage, "Die aktuelle Seite muss 1 sein, da nicht mehr Seiten erforderlich sind!")
}

// testJumpPageBackFirstPage
// Hier wird der Fall abgedeckt, falls nue eine Seite zud Darstellung der Termine benötigt wird
// und der Benutzer dennoch eine Seite vor springen möchte
func testFilterJumpPageBackFirstPage(t *testing.T) {
	//Filteransicht erstellen, angezeigte Seite ist 1, keine Einträge vorhanden -> nur eine Seite benötigt
	fv := InitFilterView([]ds.Termin{})

	fv.JumpPageBack()
	assert.Equal(t, 1, fv.CurrentPage, "Die aktuelle Seite muss 1 sein, da diese nicht negativ sein kann!")
}
func testFilterJumpPageForward(t *testing.T) {

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//Filteransicht erstellen, angezeigte Seite ist 1, 30 Einträge vorhanden -> 6 Seite benötigt
	fv := InitFilterView(testTermine30)

	//5 Einträge pro Tag -> das heißt es muss 6 Seiten geben
	fv.SelectEntriesPerPage(5)
	assert.Equal(t, 6, fv.RequiredPages(), "Es sollten 6 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 6; pageJump++ {
		assert.Equal(t, pageJump, fv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		fv.JumpPageForward()
	}

	// Angezeigte Seite wieder auf Seite 1 setzten
	fv.CurrentPage = 1

	//10 Einträge pro Tag -> das heißt es muss 3 Seiten geben
	fv.SelectEntriesPerPage(10)
	assert.Equal(t, 3, fv.RequiredPages(), "Es sollten 3 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 3; pageJump++ {
		assert.Equal(t, pageJump, fv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		fv.JumpPageForward()
	}
	// Angezeigte Seite wieder auf Seite 1 setzten
	fv.CurrentPage = 1

	//15 Einträge pro Tag -> das heißt es muss 2 Seiten geben
	fv.SelectEntriesPerPage(15)
	assert.Equal(t, 2, fv.RequiredPages(), "Es sollten 2 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 2; pageJump++ {
		assert.Equal(t, pageJump, fv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		fv.JumpPageForward()
	}
}
func testFilterJumpPageBack(t *testing.T) {

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//Filteransicht erstellen, angezeigte Seite ist 1, 15 Einträge vorhanden -> 3 Seite benötigt bei 5 Einträgen pro Seite
	fv := InitFilterView(testTermine30[:15])

	//5 Einträge pro Tag -> das heißt es muss 3 Seiten geben
	fv.SelectEntriesPerPage(5)

	// Angezeigte Seite wird auf letzte Seite gesetzt
	fv.CurrentPage = fv.RequiredPages()
	assert.Equal(t, 3, fv.CurrentPage, "Es sollten 6 Seiten benötigt werden.")
	for pageJump := 3; pageJump > 0; pageJump-- {
		assert.Equal(t, pageJump, fv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		fv.JumpPageBack()
	}

	//10 Einträge pro Tag -> das heißt es muss 2 Seiten geben
	fv.SelectEntriesPerPage(10)
	// Angezeigte Seite wieder auf letzte Seite setzten
	fv.CurrentPage = fv.RequiredPages()
	assert.Equal(t, 2, fv.CurrentPage, "Es sollten 3 Seiten benötigt werden.")
	for pageJump := 2; pageJump > 0; pageJump-- {
		assert.Equal(t, pageJump, fv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		fv.JumpPageBack()
	}

	//15 Einträge pro Tag -> das heißt es muss 1 Seiten geben
	fv.SelectEntriesPerPage(15)
	// Angezeigte Seite wieder auf letzte Seite setzten
	fv.CurrentPage = fv.RequiredPages()
	assert.Equal(t, 1, fv.CurrentPage, "Es sollten 2 Seiten benötigt werden.")
	fv.JumpPageBack()
	//Da man sich auf der letzten Seite befindet, bleibt man auf Seite 1
	assert.Equal(t, 1, fv.CurrentPage, "Es sollte Seite 1 angezeigt werden.")
}

/*
*********************************************************************************************************************
Ab hier Folgen Funktionen, die dem Filtern und Sortieren der Termine in der Filteransicht dienen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testFilterSortEntries(t *testing.T) {
	var fv FilterView

	//Datum für Testtermine erstellen
	testTermin1Starts := createSpecificDate(2022, 9, 11)
	testTermin1Ends := createSpecificDate(2022, 29, 11)

	testTermin2Starts := createSpecificDate(2022, 27, 10)
	testTermin2Ends := createSpecificDate(2023, 27, 10)

	testTermin3Starts := createSpecificDate(2022, 9, 8)
	testTermin3Ends := createSpecificDate(2023, 9, 11)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 3)

	// Terminstart: 9.11.2022
	testTermine[0] = NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends, false)

	// Terminstart: 27.10.2022
	testTermine[1] = NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends, false)

	// Terminstart: 9.8.2022
	testTermine[2] = NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends, false)

	//Daten für die Tests in ein Slice mit Terminen kopieren
	controlTermins := make([]ds.Termin, len(testTermine))
	copy(controlTermins, testTermine)

	//Daten sortieren
	fv.SortEntries(testTermine)

	//Testen ob Daten in Slice richtig eingefügt wurden:
	assert.Equal(t, controlTermins[2], testTermine[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, controlTermins[1], testTermine[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, controlTermins[0], testTermine[2], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, 3, len(testTermine), "3 Termine sollten sich in der Slice befinden")
}

func testFilterNextOccurrences(t *testing.T) {
	var fv FilterView

	todayYear := time.Now().Year()
	todayMonth := time.Now().Month()
	todayDay := time.Now().Day()
	today := time.Date(todayYear, todayMonth, todayDay, 0, 0, 0, 0, time.UTC)
	//Slice mit den zu erwartenden Ergebnissen
	expectedResults := make([][]time.Time, 5)

	//Datum für Testtermine erstellen

	//täglicher Termin: liegt in vergangenheit, daher Ergebnis leer
	testTermin1Starts := createSpecificDate(2022, 9, 11)
	testTermin1Ends := createSpecificDate(2022, 29, 11)
	expectedResults[0] = make([]time.Time, 0, 3)

	//Wöchentlicher Termin
	testTermin2Starts := today
	testTermin2Ends := testTermin2Starts.AddDate(1, 0, 0)
	expectedResults[1] = []time.Time{createSpecificDate(todayYear, todayDay, todayMonth), createSpecificDate(todayYear, todayDay+7, todayMonth), createSpecificDate(todayYear, todayDay+14, todayMonth)}

	//Jährlicher Termin
	testTermin3Starts := today.AddDate(-1, 0, 0)
	testTermin3Ends := testTermin3Starts.AddDate(2, 0, 0)
	expectedResults[2] = []time.Time{createSpecificDate(todayYear, todayDay, todayMonth), createSpecificDate(todayYear+1, todayDay, todayMonth)}

	//Monatlicher Termin
	testTermin4Starts := today.AddDate(1, 0, 0)
	testTermin4Ends := testTermin4Starts.AddDate(2, 0, 0)
	expectedResults[3] = []time.Time{createSpecificDate(todayYear+1, todayDay, todayMonth), createSpecificDate(todayYear+1, todayDay, todayMonth+1), createSpecificDate(todayYear+1, todayDay, todayMonth+2)}

	//sich nie wiederholender Termin
	testTermin5Starts := today
	testTermin5Ends := testTermin5Starts
	expectedResults[4] = []time.Time{createSpecificDate(todayYear, todayDay, todayMonth)}

	//Slice mit Testterminen erstellen, jeder Wiederholungstyp dabei
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends, false)
	testTermine[1] = NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends, false)
	testTermine[2] = NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends, false)
	testTermine[3] = NewTerminObj("testTermin4", "test", ds.MONTHLY, testTermin4Starts, testTermin4Ends, false)
	testTermine[4] = NewTerminObj("testTermin5", "test", ds.Never, testTermin5Starts, testTermin5Ends, false)

	assert.Equal(t, expectedResults[0], fv.NextOccurrences(testTermine[0]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[1], fv.NextOccurrences(testTermine[1]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[2], fv.NextOccurrences(testTermine[2]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[3], fv.NextOccurrences(testTermine[3]), "Daten-in den Slices sollten identisch sein")
	assert.Equal(t, expectedResults[4], fv.NextOccurrences(testTermine[4]), "Daten-in den Slices sollten identisch sein")
}

func testFilterTermins(t *testing.T) {
	var fv FilterView

	todayYear := time.Now().Year()
	todayMonth := time.Now().Month()
	todayDay := time.Now().Day()
	today := time.Date(todayYear, todayMonth, todayDay, 0, 0, 0, 0, time.UTC)

	//Daten für Testtermine erstellen

	//Slice mit Testterminen erstellen, jeder Wiederholungstyp dabei
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = NewTerminObj("test go", "hui", ds.DAILY, today, today.AddDate(1, 0, 0), false)
	testTermine[1] = NewTerminObj("test ist", "lala", ds.WEEKLY, today, today.AddDate(1, 0, 0), false)
	testTermine[2] = NewTerminObj("test eine", "ich bin toll", ds.YEARLY, today, today.AddDate(1, 0, 0), false)
	testTermine[3] = NewTerminObj("test Sprache", "ich bin super", ds.MONTHLY, today, today.AddDate(1, 0, 0), false)
	testTermine[4] = NewTerminObj("test go", "tada", ds.Never, today, today, false)

	saerchTitle := "test"
	saerchDescription := ""

	fv.FilterTermins(saerchTitle, saerchDescription, testTermine)
	assert.Equal(t, 5, len(fv.FilteredTermins), "Alle Termine sollten herausgefiltert worden sein, die Länge sollte dementsprechend 5 sein.")

	saerchTitle = "test"
	saerchDescription = "ich bin"
	fv.FilterTermins(saerchTitle, saerchDescription, testTermine)
	assert.Equal(t, 2, len(fv.FilteredTermins), "2 Termine sollten herausgefiltert worden sein.")

	saerchTitle = "test go"
	saerchDescription = ""
	fv.FilterTermins(saerchTitle, saerchDescription, testTermine)
	assert.Equal(t, 2, len(fv.FilteredTermins), "2 Termine sollten herausgefiltert worden sein.")

	saerchTitle = ""
	saerchDescription = "bin"
	fv.FilterTermins(saerchTitle, saerchDescription, testTermine)
	assert.Equal(t, 2, len(fv.FilteredTermins), "2 Termine sollten herausgefiltert worden sein.")

	saerchTitle = ""
	saerchDescription = ""
	fv.FilterTermins(saerchTitle, saerchDescription, testTermine)
	assert.Equal(t, 5, len(fv.FilteredTermins), "Alle 5 Termine sollten herausgefiltert worden sein.")

}

/*
*********************************************************************************************************************
Alle tests ausführen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func TestFilterView(t *testing.T) {

	//Tests für Custom-Settings innerhalb der Webseite
	t.Run("testRuns SelectEntriesPerPage", testFilterSelectEntriesPerPage)
	t.Run("testRuns JumPageBackForFirstPage", testFilterJumpPageBackFirstPage)
	t.Run("testRuns JumPageForwardForLastPage", testFilterJumPageForwardLastPage)
	t.Run("testRuns JumpPageForward", testFilterJumpPageForward)
	t.Run("testRuns JumpPageBack", testFilterJumpPageBack)
	t.Run("testRuns GetEntries", testGetEntries)

	//Tests zum Filtern und für die korrekte Darstellung der Termine in der Filteransicht
	t.Run("testRuns SortEntries", testFilterSortEntries)
	t.Run("testRuns NextOccurrences", testFilterNextOccurrences)
	t.Run("testRuns FilterTermins", testFilterTermins)
}
