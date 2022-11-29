package kalenderansicht

import (
	ds "DHBW_GO_Projekt/dateisystem"
	"fmt"
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
	assert.Equal(t, 1, lv.CurrentPage, "Die aktuelle Seite sollte Seite 1 sein.")
}

func testSelectEntriesPerPage(t *testing.T) {
	var lv = new(ListView)
	entriesAmountValue := 5
	lv.SelectEntriesPerPage(entriesAmountValue)
	assert.Equal(t, entriesAmountValue, lv.EntriesPerPage, "Die Nummern sollten identisch sein.")
	assert.Equal(t, 1, lv.CurrentPage, "Die aktuelle Seite sollte Seite 1 sein.")
}

/*
**************************************************************************************************************
Test zur Navigation innerhalb der Listenansicht und für die Anzeige der richtigen Anzahl an Einträgen
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
//testJumPageForLastPage
//Hier wird der Fall abgedeckt, dass der Benutzer sich bereits auf der letzten, zur Darstellung der Termine benötigten Seite
//befindet und dennoch eine Seite vorspringen möchte
func testJumPageForLastPage(t *testing.T) {
	//Liste erstellen, angezeigte Seite ist 1, keine Einträge vorhanden -> nur eine Seite benötigt
	lv := InitListView([]ds.Termin{})

	lv.JumpPageForward()
	assert.Equal(t, 1, lv.CurrentPage, "Die aktuelle Seite muss 1 sein, da nicht mehr Seiten erforderlich sind!")
}

// testJumpPageBackFirstPage
// Hier wird der Fall abgedeckt, falls nue eine Seite zud Darstellung der Termine benötigt wird
// und der Benutzer dennoch eine Seite vor springen möchte
func testJumpPageBackFirstPage(t *testing.T) {
	//Liste erstellen, angezeigte Seite ist 1, keine Einträge vorhanden -> nur eine Seite benötigt
	lv := InitListView([]ds.Termin{})

	lv.JumpPageBack()
	assert.Equal(t, 1, lv.CurrentPage, "Die aktuelle Seite muss 1 sein, da diese nicht negativ sein kann!")
}
func testJumPageFor(t *testing.T) {

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//Listen-Ansicht erstellen, angezeigte Seite ist 1
	lv := InitListView(testTermine30)

	//5 Einträge pro Tag -> das heißt es muss 6 Seiten geben
	lv.SelectEntriesPerPage(5)
	assert.Equal(t, 6, lv.RequiredPages(), "Es sollten 6 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 6; pageJump++ {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageForward()
	}

	// Angezeigte Seite wieder auf Seite 1 setzten
	lv.CurrentPage = 1

	//10 Einträge pro Tag -> das heißt es muss 3 Seiten geben
	lv.SelectEntriesPerPage(10)
	assert.Equal(t, 3, lv.RequiredPages(), "Es sollten 3 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 3; pageJump++ {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageForward()
	}
	// Angezeigte Seite wieder auf Seite 1 setzten
	lv.CurrentPage = 1

	//15 Einträge pro Tag -> das heißt es muss 2 Seiten geben
	lv.SelectEntriesPerPage(15)
	assert.Equal(t, 2, lv.RequiredPages(), "Es sollten 2 Seiten benötigt werden.")
	for pageJump := 1; pageJump <= 2; pageJump++ {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageForward()
	}
}
func testJumpPageBack(t *testing.T) {

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	//Listen-Ansicht erstellen
	lv := InitListView(testTermine30)

	//5 Einträge pro Tag -> das heißt es muss 6 Seiten geben
	lv.SelectEntriesPerPage(5)

	// Angezeigte Seite wird auf letzte Seite gesetzt
	lv.CurrentPage = lv.RequiredPages()
	assert.Equal(t, 6, lv.CurrentPage, "Es sollten 6 Seiten benötigt werden.")
	for pageJump := 6; pageJump > 0; pageJump-- {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageBack()
	}

	//10 Einträge pro Tag -> das heißt es muss 3 Seiten geben
	lv.SelectEntriesPerPage(10)
	// Angezeigte Seite wieder auf letzte Seite setzten
	lv.CurrentPage = lv.RequiredPages()
	assert.Equal(t, 3, lv.CurrentPage, "Es sollten 3 Seiten benötigt werden.")
	for pageJump := 3; pageJump > 0; pageJump-- {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageBack()
	}

	//15 Einträge pro Tag -> das heißt es muss 2 Seiten geben
	lv.SelectEntriesPerPage(15)
	// Angezeigte Seite wieder auf letzte Seite setzten
	lv.CurrentPage = lv.RequiredPages()
	assert.Equal(t, 2, lv.CurrentPage, "Es sollten 2 Seiten benötigt werden.")
	for pageJump := 2; pageJump > 0; pageJump-- {
		assert.Equal(t, pageJump, lv.CurrentPage, "Es sollten Seite "+fmt.Sprint(pageJump)+" angezeigt werden.")
		lv.JumpPageBack()
	}
}

// testGetEntriesCorrectNumber
// hier wird getestet, ob die Anzahl der Einträge mit der eingestellten überein stimmt
func testGetEntriesCorrectNumber(t *testing.T) {
	//Listen-Ansicht erstellen
	lv := new(ListView)
	//Hier wird das Startdatum gesetzt, abhängig vom heutigen Datum -> Testausführung unabhängig vom eig. Datum
	lv.SelectedDate = createSpecificDate(time.Now().Year()-1, time.Now().Day(), time.Now().Month())
	lv.EntriesPerPage = 5
	lv.CurrentPage = 1

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	lv.CreateTerminListEntries(testTermine30)

	assert.Equal(t, 30, len(lv.EntriesSinceSelDate), "Es sollten 30 Termine in dem Slice sein")

	//5 Einträge pro Tag
	lv.SelectEntriesPerPage(5)
	sliceWithEntries := lv.GetEntries()
	assert.Equal(t, 5, len(sliceWithEntries), "Es sollten 5 Termine in dem Slice sein")

	//10 Einträge pro Tag
	lv.SelectEntriesPerPage(10)
	sliceWithEntries = lv.GetEntries()
	assert.Equal(t, 10, len(sliceWithEntries), "Es sollten 10 Termine in dem Slice sein")

	//15 Einträge pro Tag
	lv.SelectEntriesPerPage(15)
	sliceWithEntries = lv.GetEntries()
	assert.Equal(t, 15, len(sliceWithEntries), "Es sollten 15 Termine in dem Slice sein")
}

// testGetCorrectEntries
// hier wird getestet, ob die Einträge stimmen, wenn die nächste Seite aufgerufen wird
func testGetEntriesCorrectEntries(t *testing.T) {

	//Listen-Ansicht erstellen
	lv := new(ListView)
	//Hier wird das Startdatum gesetzt, abhängig vom heutigen Datum -> Testausführung unabhängig vom eig. Datum
	lv.SelectedDate = createSpecificDate(time.Now().Year()-1, time.Now().Day(), time.Now().Month())
	lv.EntriesPerPage = 5
	lv.CurrentPage = 1

	//Falls der Slice mit Testterminen noch nicht erstellt worden ist, diesen erstellen
	//Ist der Fall, wenn Test einzeln ausgeführt wird
	if len(testTermine30) == 0 {
		testTermine30 = create30TestTermins()
	}

	lv.CreateTerminListEntries(testTermine30)

	//Termine filtern
	filteredSlice := lv.FilterCalendarEntries(testTermine30)
	assert.Equal(t, 30, len(filteredSlice), "Es sollten 30 Termine in dem Slice sein")

	//5 Einträge pro Tag, Test mit Seite vor und zurück springen
	lv.SelectEntriesPerPage(5)
	assert.Equal(t, filteredSlice[:5], lv.GetEntries(), "Es sollten die ersten 5 Termine in dem Slice sein")
	lv.JumpPageForward()
	assert.Equal(t, filteredSlice[5:10], lv.GetEntries(), "Es sollten die nächsten 5 Termine in dem Slice sein")
	lv.JumpPageBack()
	assert.Equal(t, filteredSlice[:5], lv.GetEntries(), "Es sollten die ersten 5 Termine in dem Slice sein")

	//10 Einträge pro Tag, Test mit Seite vor und zurück springen
	lv.SelectEntriesPerPage(10)
	assert.Equal(t, filteredSlice[:10], lv.GetEntries(), "Es sollten die ersten 10 Termine in dem Slice sein")
	lv.JumpPageForward()
	assert.Equal(t, filteredSlice[10:20], lv.GetEntries(), "Es sollten die nächsten 10 Termine in dem Slice sein")
	lv.JumpPageBack()
	assert.Equal(t, filteredSlice[:10], lv.GetEntries(), "Es sollten die ersten 10 Termine in dem Slice sein")

	//15 Einträge pro Tag, Test mit Seite vor und zurück springen
	lv.SelectEntriesPerPage(15)
	assert.Equal(t, filteredSlice[:15], lv.GetEntries(), "Es sollten die ersten 15 Termine in dem Slice sein")
	lv.JumpPageForward()
	assert.Equal(t, filteredSlice[15:30], lv.GetEntries(), "Es sollten die nächsten 15 Termine in dem Slice sein")
	lv.JumpPageBack()
	assert.Equal(t, filteredSlice[:15], lv.GetEntries(), "Es sollten die ersten 15 Termine in dem Slice sein")
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
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends, false)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends, false)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends, false)
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.MONTHLY, testTermin4Starts, testTermin4Ends, false)
	//Termin sollte nicht hinzugefügt werden
	testTermine[4] = ds.NewTerminObj("testTermin5", "test", ds.MONTHLY, testTermin5Starts, testTermin5Ends, false)

	//Testen ob Daten in Slice eingefügt wurden
	filteredSlice := lv.FilterCalendarEntries(testTermine)
	assert.Equal(t, testTermine[0], filteredSlice[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[1], filteredSlice[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[2], filteredSlice[2], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, testTermine[3], filteredSlice[3], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, 4, len(filteredSlice), "4 Termine sollten sich in der Slice befinden")

}

// testFilterCalendarEntriesEdgeCase
// hier wird getestet, ob die Funktion Termine nicht in das Slice mit aufnimmt,
// wenn das Enddatum des Termins zwar hinter dem aktuell angezeigtem Datum liegt ABER keine Wiederholung des
// Termins mehr stattfindet.
func testFilterCalendarEntriesEdgeCase(t *testing.T) {
	var lv ListView
	//Hier wird der 1.11.2022 als Startdatum gesetzt
	lv.SelectedDate = createSpecificDate(2022, 2, 11)

	//Datum für Testtermine erstellen
	testTermin1Starts := createSpecificDate(2022, 7, 10)
	testTermin1Ends := createSpecificDate(2022, 3, 11) //letzte Vorkommen: 28.10.2022

	testTermin2Starts := createSpecificDate(2020, 1, 11)
	testTermin2Ends := createSpecificDate(2022, 10, 11) //letzte Vorkommen: 1.11.2022

	testTermin3Starts := createSpecificDate(2021, 1, 10)
	testTermin3Ends := createSpecificDate(2022, 9, 11) //letzte Vorkommen: 1.11.2022

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 5)
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.WEEKLY, testTermin1Starts, testTermin1Ends, false)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.YEARLY, testTermin2Starts, testTermin2Ends, false)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.MONTHLY, testTermin3Starts, testTermin3Ends, false)

	//Testen, ob keine Daten in das Slice eingefügt wurden
	filteredSlice := lv.FilterCalendarEntries(testTermine)
	assert.Equal(t, 0, len(filteredSlice), "Es sollten keine Termine in der Slice sein.")

}

func testSortEntries(t *testing.T) {
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

	testTermin4Starts := createSpecificDate(2022, 2, 10)
	testTermin4Ends := createSpecificDate(2023, 2, 10)

	//Slice mit Testterminen erstellen
	testTermine := make([]ds.Termin, 4)

	// nächstes Vorkommen ab  1.11.2022: 9.11.2022
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends, false)

	// nächstes Vorkommen  1.11.2022: 3.11.2022
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends, false)

	// nächstes Vorkommen  1.11.2022: 9.11.2022
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends, false)

	// nächstes Vorkommen  1.11.2022: 2.11.2022
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.MONTHLY, testTermin4Starts, testTermin4Ends, false)

	//Daten für die Tests in ein Slice mit Terminen kopieren
	controlTermins := make([]ds.Termin, len(testTermine))
	copy(controlTermins, testTermine)

	//Testen ob Daten in Slice eingefügt wurden:
	lv.SortEntries(testTermine)
	assert.Equal(t, controlTermins[3], testTermine[0], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, controlTermins[1], testTermine[1], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, controlTermins[2], testTermine[2], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, controlTermins[0], testTermine[3], "Termine an diesen Positionen sollten identisch sein")
	assert.Equal(t, 4, len(testTermine), "4 Termine sollten sich in der Slice befinden")

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
	testTermine[0] = ds.NewTerminObj("testTermin1", "test", ds.DAILY, testTermin1Starts, testTermin1Ends, false)
	testTermine[1] = ds.NewTerminObj("testTermin2", "test", ds.WEEKLY, testTermin2Starts, testTermin2Ends, false)
	testTermine[2] = ds.NewTerminObj("testTermin3", "test", ds.YEARLY, testTermin3Starts, testTermin3Ends, false)
	testTermine[3] = ds.NewTerminObj("testTermin4", "test", ds.MONTHLY, testTermin4Starts, testTermin4Ends, false)
	testTermine[4] = ds.NewTerminObj("testTermin5", "test", ds.Never, testTermin5Starts, testTermin5Ends, false)

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
	//slice mit Testterminen erstellen, benötigt viel Zeit: daher ein globales Slice
	testTermine30 = create30TestTermins()

	//Tests für Custom-Settings innerhalb der Webseite
	t.Run("testRuns SelectDate", testSelectDate)
	t.Run("testRuns SelectEntriesPerPage", testSelectEntriesPerPage)

	//Test zur Navigation innerhalb der Listenansicht
	t.Run("testRuns JumPageForLastPage", testJumPageForLastPage)
	t.Run("testRuns JumpPageBackFirstPage", testJumpPageBackFirstPage)
	t.Run("testRuns JumPageFor", testJumPageFor)
	t.Run("testRuns JumpPageBack", testJumpPageBack)
	t.Run("testRuns GetEntriesCorrectNumber", testGetEntriesCorrectNumber)
	t.Run("test Run GetEntriesCorrectEntries", testGetEntriesCorrectEntries)

	//Tests zum Filtern und für die korrekte Darstellung der Termine in der Listenansicht
	t.Run("testRuns FilterCalendarEntries", testFilterCalendarEntries)
	t.Run("testRuns FilterCalendarEntriesEdgeCase", testFilterCalendarEntriesEdgeCase)
	t.Run("testRuns SortEntries", testSortEntries)
	t.Run("testRuns  NextOccurrences", testNextOccurrences)

}
