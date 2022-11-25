package dateisystem

//Mat-Nr. 8689159
import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transitiver Termine

	var k []Termin

	k = NewTerminObj("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), k) //erzeugt dummy Termin

	assert.NotEqual(t, "", k[0].Title)
	assert.NotEqual(t, "", k[0].Description)
	assert.NotEqual(t, Never, k[0].Recurring)
	assert.NotEqual(t, "", k[0].Date)
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", k[0].EndDate.String())
	DeleteFromCache(k, "0", "mik")
}

func TestLoadTermin(t *testing.T) { //prüft das Laden von Objekten

	var k []Termin

	k = CreateNewTermin("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	terminLoaded := LoadTermin(strconv.Itoa(getID()), "mik")

	assert.Equal(t, k[0], terminLoaded)
	DeleteAll(GetTermine("mik"), "mik")
}

func TestGetTermine(t *testing.T) { //prüft ob das erzeugte Slice die korrekten Objekte geladen hat

	var k []Termin
	var n []Termin

	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	CreateNewTermin("testu", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)

	kTest := GetTermine("mik")

	decrementID()
	n = NewTerminObj("testu", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), n)

	assert.Equal(t, n[0], kTest[2])

	DeleteAll(k, "mik")
}

func TestAddToCache(t *testing.T) { //prüft, ob Termin dem Caching hinzugefügt wurde

	k := GetTermine("mik")

	ter := NewTerminObj("testa", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), k)

	decrementID()
	k = AddToCache(ter[0], k)

	assert.Equal(t, ter[0], k[0])

	DeleteAll(k, "mik")
}

func TestDeleteTermin(t *testing.T) { //prüft ob die JSONs korrekt gelöscht werden

	DeleteAll(GetTermine("mik"), "mik")

	var k []Termin

	k = CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	deleteTermin("1", "mik")

	n := GetTermine("mik")

	assert.Equal(t, []Termin(nil), n)

	DeleteAll(k, "mik")
}

func TestStoreAll(t *testing.T) { //prüft, ob sich der gesamte Cache speichern lässt

	k := GetTermine("mik")
	k = CreateNewTermin("testa", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	k = CreateNewTermin("testb", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	k = CreateNewTermin("testc", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	StoreCache(k, "mik")

	assert.Equal(t, k, GetTermine("mik"))

	DeleteAll(k, "mik")
}

func TestDeleteAll(t *testing.T) { //prüft, ob ein gesamter Kalender gelöscht werden kann(Username bleibt bestehen)

	k := GetTermine("mik")
	k = CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)

	assert.NotEqual(t, []Termin(nil), k)
	k = DeleteAll(k, "mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestDeleteFromCache(t *testing.T) { //prüft, ob Termin aus dem Caching gelöscht werden kann, ohne Verzeichnis neu einlesen zu müssen
	k := GetTermine("mik")
	k = CreateNewTermin("testa", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	k = CreateNewTermin("testb", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	k = CreateNewTermin("testc", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", k)
	k2 := DeleteFromCache(k, "2", "mik")

	assert.Equal(t, k[2], k2[1])

	DeleteAll(k, "mik")
}
