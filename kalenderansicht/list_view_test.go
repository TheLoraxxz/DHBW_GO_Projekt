package kalenderansicht

import (
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
	entriesAmount := 10
	lv.SelectEntriesPerPage(entriesAmount)
	assert.Equal(t, entriesAmount, lv.EntriesPerPage, "Die Nummern sollten identisch sein.")
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
}
