package dateisystem

//Mat-Nr. 8689159
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

func TestAddToCache(t *testing.T) { //prüft, ob Termin dem Caching hinzugefügt wurde
	k := GetTermine("Mik")
	ter := NewTerminObj("testa", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC")

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
	k := GetTermine("Mik")
	k = AddToCache(NewTerminObj("testa", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	k = AddToCache(NewTerminObj("testb", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	k = AddToCache(NewTerminObj("testc", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	StoreCache(k, "Mik")

	assert.Equal(t, k, GetTermine("Mik"))
}

func TestDeleteAll(t *testing.T) { //prüft, ob ein gesamter Kalender gelöscht werden kann(Username bleibt bestehen)
	k := GetTermine("Mik")
	assert.NotEqual(t, []Termin(nil), k)
	k = DeleteAll(k, "Mik")
	assert.Equal(t, []Termin(nil), k)
}

func TestDeleteFromCache(t *testing.T) { //prüft, ob Termin aus dem Caching gelöscht werden kann, ohne Verzeichnis neu einlesen zu müssen
	k := GetTermine("Mik")
	k = AddToCache(NewTerminObj("testa", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	k = AddToCache(NewTerminObj("testb", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	k = AddToCache(NewTerminObj("testc", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC"), k)
	k2 := DeleteFromCache(k, "testb", "Mik")

	assert.Equal(t, k[2], k2[1])
}
