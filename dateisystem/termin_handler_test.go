package dateisystem

//Mat-Nr. 8689159
import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transitiver Termine
	termin := NewTerminObj("test", "test", Repeat(WEEKLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)) //erzeugt dummy Termin

	assert.NotEqual(t, "", termin.Title)
	assert.NotEqual(t, "", termin.Description)
	assert.NotEqual(t, Repeat(Never), termin.Recurring)
	assert.NotEqual(t, "", termin.Date)
	assert.Equal(t, "2007-03-02 15:02:05 +0000 UTC", termin.EndDate.String())
}

func TestLoadTermin(t *testing.T) { //prüft das Laden von Objekten
	termin := CreateNewTermin("test", "test", Repeat(WEEKLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik")
	terminLoaded := LoadTermin("test", "mik")

	assert.Equal(t, termin, terminLoaded)
}

func TestGetTermine(t *testing.T) { //prüft ob das erzeugte Slice die korrekten Objekte geladen hat
	CreateNewTermin("testo", "test", Repeat(WEEKLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik")
	i := CreateNewTermin("testu", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik")
	k := GetTermine("mik")

	assert.Equal(t, i, k[2])
}

func TestAddToCache(t *testing.T) { //prüft, ob Termin dem Caching hinzugefügt wurde
	k := GetTermine("mik")
	ter := NewTerminObj("testa", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC))

	k = AddToCache(ter, k)

	assert.Equal(t, ter, k[3])
}

func TestDeleteTermin(t *testing.T) { //prüft ob die JSONs korrekt gelöscht werden
	deleteTermin("test", "mik")
	deleteTermin("testo", "mik")
	deleteTermin("testu", "mik")

	k := GetTermine("mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestStoreAll(t *testing.T) { //prüft, ob sich der gesamte Cache speichern lässt
	k := GetTermine("mik")
	k = AddToCache(NewTerminObj("testa", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	k = AddToCache(NewTerminObj("testb", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	k = AddToCache(NewTerminObj("testc", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	StoreCache(k, "mik")

	assert.Equal(t, k, GetTermine("mik"))
}

func TestDeleteAll(t *testing.T) { //prüft, ob ein gesamter Kalender gelöscht werden kann(Username bleibt bestehen)
	k := GetTermine("mik")
	assert.NotEqual(t, []Termin(nil), k)
	k = DeleteAll(k, "mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestDeleteFromCache(t *testing.T) { //prüft, ob Termin aus dem Caching gelöscht werden kann, ohne Verzeichnis neu einlesen zu müssen
	k := GetTermine("mik")
	k = AddToCache(NewTerminObj("testa", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	k = AddToCache(NewTerminObj("testb", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	k = AddToCache(NewTerminObj("testc", "test", Repeat(YEARLY), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC)), k)
	k2 := DeleteFromCache(k, "testb", "mik")

	assert.Equal(t, k[2], k2[1])
}
