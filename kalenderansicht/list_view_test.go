package kalenderansicht

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
**************************************************************************************************************
Tests für Custom-Settings innerhalb der Webseite
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testSelectDate(t *testing.T) {
	lv := InitListView()
	newDate := ganerateRandomDate()
	lv.SelectDate(newDate)
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
