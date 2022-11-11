package kalenderansicht

import (
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
Test zur Datenbeschaffung für die Tabellenansicht
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testJahrAnzeige(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	jahr := ta.ShownDate.Year()
	assert.Equal(t, jahr, ta.ShowYear(), "Die Jahre sollten identisch sein.")
}

func testMonatsAnzeige(t *testing.T) {
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
func testSpringMonatVor(t *testing.T) {
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
func testSpringMonatZurueck(t *testing.T) {
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
func testSpringeJahr(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	aktuellesJahr := ta.ShownDate.Year()
	ta.JumpToYear(1)
	assert.Equal(t, aktuellesJahr+1, ta.ShownDate.Year(), "Die Jahre sollten identisch sein.")
	ta.JumpToYear(-1)
	assert.Equal(t, aktuellesJahr, ta.ShownDate.Year(), "Die Jahre sollten identisch sein.")
}

func testWaehleMonat(t *testing.T) {
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

func testSpringZuHeute(t *testing.T) {
	var ta TableView
	ta.ShownDate = ta.getFirstDayOfMonth(ganerateRandomDate())
	ta.JumpToToday()
	assert.Equal(t, ta.ShownDate, ta.getFirstDayOfMonth(time.Now()), "Das Datum sollte das heutige sein.")
}

func TestTabellenAnsicht(t *testing.T) {
	t.Run("testRuns", testJahrAnzeige)
	t.Run("testRuns", testMonatsAnzeige)
	t.Run("testRuns", testSpringMonatVor)
	t.Run("testRuns", testSpringMonatZurueck)
	t.Run("testRuns", testSpringeJahr)
	t.Run("testRuns", testWaehleMonat)
	t.Run("testRuns", testSpringZuHeute)

}
