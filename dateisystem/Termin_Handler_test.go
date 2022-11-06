package dateisystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transitiver Termine

	termin := NewTerminObj("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC") //erzeugt dummy Termin

	assert.NotEqual(t, "", termin.Title)
	assert.NotEqual(t, "", termin.Description)
	assert.NotEqual(t, repeat(niemals), termin.Recurring)
	assert.NotEqual(t, "", termin.Date)
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", termin.EndDate.String())
}

func TestLoadTermin(t *testing.T) { //prüft das Laden von Objekten
	termin := CreateNewTermin("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	terminLoaded := LoadTermin("test", "mik")

	assert.Equal(t, termin, terminLoaded)
}

func TestGetTermine(t *testing.T) { //prüft ob das erzeugte Slice die korrekten Objekte geladen hat
	CreateNewTermin("testo", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	i := CreateNewTermin("testu", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	k := GetTermine("Mik")

	assert.Equal(t, i, k[2])
}

func TestAddToCache(t *testing.T) {
	k := GetTermine("Mik")
	ter := NewTerminObj("testa", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC")

	k = AddToCache(ter, k)

	assert.Equal(t, ter, k[3])
}

func TestDeleteTermin(t *testing.T) {
	DeleteTermin("test", "mik")
	DeleteTermin("testo", "mik")
	DeleteTermin("testu", "mik")
	DeleteTermin("testa", "mik")

	k := GetTermine("mik")
	assert.Equal(t, []Termin(nil), k)
}
