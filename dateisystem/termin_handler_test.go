package dateisystem

//Mat-Nr. 8689159
import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transitiver Termine

	ter := NewTerminObj("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true) //erzeugt dummy Termin

	assert.NotEqual(t, "", ter.Title)
	assert.NotEqual(t, "", ter.Description)
	assert.NotEqual(t, Never, ter.Recurring)
	assert.NotEqual(t, "", ter.Date)
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", ter.EndDate.String())
	assert.Equal(t, true, ter.Shared)

}

func TestLoadTermin(t *testing.T) { //prüft das Laden von Objekten

	DeleteAll(GetTermine("mik"), "mik")

	ter := CreateNewTermin("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")

	terminLoaded := LoadTermin(ter.ID, "mik")

	assert.Equal(t, ter, terminLoaded)
	DeleteAll(GetTermine("mik"), "mik")
}

func TestGetTermine(t *testing.T) { //prüft ob das erzeugte Slice die korrekten Objekte geladen hat

	DeleteAll(GetTermine("mik"), "mik")

	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	ter := CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	CreateNewTermin("testu", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")

	kTest := GetTermine("mik")

	assert.Equal(t, ter, kTest[1])

	DeleteAll(kTest, "mik")
}

func TestAddToCache(t *testing.T) { //prüft, ob Termin dem Caching hinzugefügt wurde

	DeleteAll(GetTermine("mik"), "mik")
	var k []Termin

	ter := NewTerminObj("testa", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true)

	k = AddToCache(ter, k)

	assert.Equal(t, ter, k[0])

	DeleteAll(k, "mik")
}

func TestDeleteTermin(t *testing.T) { //prüft ob die JSONs korrekt gelöscht werden

	DeleteAll(GetTermine("mik"), "mik")

	n := CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	deleteTermin(n.ID, "mik")

	k := GetTermine("mik")

	assert.Equal(t, []Termin(nil), k)
}

func TestStoreCache(t *testing.T) { //prüft, ob sich der gesamte Cache speichern lässt

	DeleteAll(GetTermine("mik"), "mik")
	var k []Termin

	k = AddToCache(CreateNewTermin("testa", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)
	ter := CreateNewTermin("testb", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	k = AddToCache(ter, k)
	k = AddToCache(CreateNewTermin("testc", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)
	StoreCache(k, "mik")

	terExp := FindInCacheByID(k, ter.ID)
	terAct := FindInCacheByID(GetTermine("mik"), ter.ID)

	assert.Equal(t, terExp, terAct)

	DeleteAll(k, "mik")
}

func TestDeleteAll(t *testing.T) { //prüft, ob ein gesamter Kalender gelöscht werden kann(Username bleibt bestehen)

	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	k := GetTermine("mik")
	assert.NotEqual(t, []Termin(nil), k)
	k = DeleteAll(k, "mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestDeleteFromCache(t *testing.T) { //prüft, ob Termin aus dem Caching gelöscht werden kann, ohne Verzeichnis neu einlesen zu müssen

	DeleteAll(GetTermine("mik"), "mik")

	k := GetTermine("mik")
	k = AddToCache(CreateNewTermin("testa", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)
	ter := CreateNewTermin("testb", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	k = AddToCache(ter, k)
	k = AddToCache(CreateNewTermin("testc", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)
	k2 := DeleteFromCache(k, ter.ID, "mik")

	assert.Equal(t, k[2], k2[1])

	DeleteAll(k, "mik")
}

func TestFindInCacheByID(t *testing.T) {
	DeleteAll(GetTermine("mik"), "mik")

	k := GetTermine("mik")
	k = AddToCache(CreateNewTermin("testa", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)
	ter := CreateNewTermin("testb", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik")
	k = AddToCache(ter, k)
	k = AddToCache(CreateNewTermin("testc", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), true, "mik"), k)

	ter = FindInCacheByID(k, ter.ID)

	assert.Equal(t, k[1], ter)

	DeleteAll(k, "mik")
}
