package dateisystem

//Mat-Nr. 8689159
import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTermin() *Termin { //erzeugt Pointer auf dummy Termin
	t := Termin{Title: "testTermin"}
	t.Description = "test"
	t.Recurring = repeat(niemals)
	t.Date = time.Date(2021, 8, 15, 14, 30, 45, 0, time.UTC)
	t.EndDate = time.Date(2021, 8, 15, 15, 30, 45, 0, time.UTC)
	return &t
}

func updateTermin(termin *Termin) { //führt die setter aus
	termin.SetTitle(termin, "testTerminj")
	termin.SetDescription(termin, "testo yeet")
	termin.SetRecurring(termin, repeat(woechentlich))
	termin.SetDate(termin, "2007-03-02T14:02:05 UTC")
	termin.SetEndeDate(termin, "2007-03-02T15:02:05 UTC")
}

func TestTermin(t *testing.T) { //prüft ob der dummy Termin nicht Leer ist
	termin := newTermin()

	assert.NotEqual(t, "", termin.Title)
	assert.NotEqual(t, "", termin.Description)
	assert.Equal(t, repeat(niemals), termin.Recurring)
	assert.NotEqual(t, "", termin.Date)
	assert.Equal(t, "2021-08-15 15:30:45 +0000 UTC", termin.EndDate.String())
}

func TestTerminUpdate(t *testing.T) { // prüft, ob die updates durchgeführt wurden
	termin := newTermin()
	updateTermin(termin) //ruft die setter auf

	assert.Equal(t, "testTerminj", termin.Title)
	assert.Equal(t, "testo yeet", termin.Description)
	assert.Equal(t, repeat(woechentlich), termin.Recurring)
	assert.Equal(t, "2007-03-02 14:02:05 +0000 UTC", termin.Date.String())
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", termin.EndDate.String())

}
