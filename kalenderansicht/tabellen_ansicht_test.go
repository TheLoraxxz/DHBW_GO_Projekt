package kalenderansicht

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func generiereZufallsdatum() time.Time {
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
		rand.Intn(25),
		rand.Intn(61),
		rand.Intn(61),
		rand.Intn(61),
		time.UTC,
	)
}

func testJahrAnzeige(t *testing.T) {
	var ta TabellenAnsicht
	ta.Datumsanzeige = generiereZufallsdatum()
	jahr := ta.Datumsanzeige.Year()
	assert.Equal(t, jahr, ta.JahrAnzeige(), "Die Jahre sollten identisch sein.")
}

func testMonatsAnzeige(t *testing.T) {
	var ta TabellenAnsicht
	ta.Datumsanzeige = generiereZufallsdatum()
	monat := ta.Datumsanzeige.Month()
	assert.Equal(t, monat, ta.MonatsAnzeige(), "Die Monate sollten identisch sein.")
}

func testErstelleKalenderEintraege(t *testing.T) {

}
func testSpringMonatVor(t *testing.T) {
	i := 100
	var ta TabellenAnsicht
	for i > 0 {
		i -= 1
		ta.Datumsanzeige = generiereZufallsdatum()
		aktuellerMonat := (int(ta.Datumsanzeige.Month()))
		ta.SpringMonatVor()
		neuerMonat := (int(ta.Datumsanzeige.Month()))
		if aktuellerMonat < 12 {
			assert.Equal(t, aktuellerMonat+1, neuerMonat)
		} else {
			assert.Equal(t, 1, neuerMonat)
		}
	}
}
func testSpringMonatZurueck(t *testing.T) {
	i := 100
	var ta TabellenAnsicht
	for i > 0 {
		i -= 1
		ta.Datumsanzeige = generiereZufallsdatum()
		aktuellerMonat := (int(ta.Datumsanzeige.Month()))
		ta.SpringMonatZurueck()
		neuerMonat := (int(ta.Datumsanzeige.Month()))
		if aktuellerMonat == 1 {
			assert.Equal(t, 12, neuerMonat, "Die Monate sollten identisch sein.")
		} else {
			assert.Equal(t, aktuellerMonat-1, neuerMonat, "Die Monate sollten identisch sein.")
		}
	}
}

func testWaehleMonat(t *testing.T) {
	i := 100
	var ta TabellenAnsicht
	for i > 0 {
		i -= 1
		ta.Datumsanzeige = generiereZufallsdatum()
		gewaehlterMonat := time.Month(rand.Intn(12-1) + 1)
		ta.WaehleMonat(gewaehlterMonat)
		assert.Equal(t, ta.Datumsanzeige.Month(), gewaehlterMonat, "Die Monate sollten identisch sein.")
	}
}

func testSpringZuHeute(t *testing.T) {
	var ta TabellenAnsicht
	ta.Datumsanzeige = generiereZufallsdatum()
	ta.SpringZuHeute()
	assert.Equal(t, ta.Datumsanzeige, time.Now(), "Das Datum sollte das heutige sein.")
}

func TestTabellenAnsicht(t *testing.T) {
	t.Run("testRuns", testJahrAnzeige)
	t.Run("testRuns", testMonatsAnzeige)
	t.Run("testRuns", testErstelleKalenderEinträge)
	t.Run("testRuns", testSpringMonatVor)
	t.Run("testRuns", testSpringMonatZurueck)
	t.Run("testRuns", testWaehleMonat)
	t.Run("testRuns", testSpringZuHeute)
}
