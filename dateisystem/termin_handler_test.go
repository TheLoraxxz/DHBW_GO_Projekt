package dateisystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

package dateisystem

//Mat-Nr. 8689159
import (
"github.com/stretchr/testify/assert"
"testing"
"time"
)

func init() {
	DeleteAll(GetTermine("mik"), "mik")
}

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transitiver Termine
	termin := NewTerminObj("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1") //erzeugt dummy Termin

	assert.NotEqual(t, "", termin.Title)
	assert.NotEqual(t, "", termin.Description)
	assert.NotEqual(t, Never, termin.Recurring)
	assert.NotEqual(t, "", termin.Date)
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", termin.EndDate.String())
}

func TestLoadTermin(t *testing.T) { //prüft das Laden von Objekten
	termin := CreateNewTermin("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	terminLoaded := LoadTermin("1", "mik")

	assert.Equal(t, termin, terminLoaded)
	DeleteAll(GetTermine("mik"), "mik")
}

func TestGetTermine(t *testing.T) { //prüft ob das erzeugte Slice die korrekten Objekte geladen hat
	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	CreateNewTermin("testu", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	k := GetTermine("mik")

	i := NewTerminObj("testu", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "3")

	assert.Equal(t, i, k[2])

	DeleteAll(k, "mik")
}

func TestAddToCache(t *testing.T) { //prüft, ob Termin dem Caching hinzugefügt wurde
	k := GetTermine("mik")
	ter := NewTerminObj("testa", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1")

	k = AddToCache(ter, k)

	assert.Equal(t, ter, k[0])

	DeleteAll(k, "mik")
}

func TestDeleteTermin(t *testing.T) { //prüft ob die JSONs korrekt gelöscht werden
	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	deleteTermin("1", "mik")

	k := GetTermine("mik")
	assert.Equal(t, []Termin(nil), k)

	DeleteAll(k, "mik")
}

func TestStoreAll(t *testing.T) { //prüft, ob sich der gesamte Cache speichern lässt
	k := GetTermine("mik")
	k = AddToCache(NewTerminObj("testa", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1"), k)
	k = AddToCache(NewTerminObj("testb", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1"), k)
	k = AddToCache(NewTerminObj("testc", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1"), k)
	StoreCache(k, "mik")

	assert.Equal(t, k, GetTermine("mik"))

	DeleteAll(k, "mik")
}

func TestDeleteAll(t *testing.T) { //prüft, ob ein gesamter Kalender gelöscht werden kann(Username bleibt bestehen)
	CreateNewTermin("testo", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "1")
	k := GetTermine("mik")
	assert.NotEqual(t, []Termin(nil), k)
	k = DeleteAll(k, "mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestDeleteFromCache(t *testing.T) { //prüft, ob Termin aus dem Caching gelöscht werden kann, ohne Verzeichnis neu einlesen zu müssen
	k := GetTermine("mik")
	k = AddToCache(NewTerminObj("testa", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "1"), k)
	k = AddToCache(NewTerminObj("testb", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "2"), k)
	k = AddToCache(NewTerminObj("testc", "test", YEARLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "3"), k)
	k2 := DeleteFromCache(k, "2", "mik")

	assert.Equal(t, k[2], k2[1])

	DeleteAll(k, "mik")
}