package dateisystem

import (
	"fmt"
	"testing"
)

func TestNewTerminObj(t *testing.T) { //prüft das erstellen transistiver Termine

	termin := newTerminObj("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC") //erzeugt dummy Termin

	if termin.Title == "" {
		t.Errorf("Titel nicht hinterlegt")
	}
	if termin.Description == "" {
		t.Errorf("Beschreibung nicht hinterlegt")
	}
	if termin.Recurring != 1 {
		t.Errorf("Wiederholung nicht hinterlegt")
	}
	if termin.Date.String() == "" {
		t.Errorf("Datum nicht hinterlegt")
	}

	if termin.EndDate.String() != "2007-03-02 15:02:05 +0000 UTC" {
		t.Errorf("End-Datum nicht hinterlegt")
	}
}

func TestStoreTermin(t *testing.T) {
	termin := loadTermin("test", "mik")
	updateTermin(&termin)
	storeTerminObj(termin, "mik")
}

func TestLoadTermin(t *testing.T) {
	termin := loadTermin("test", "mik")
	if termin != createNewTermin("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik") {
		t.Errorf("Objekte passen nicht zusammen")
	}
}

func TestDeleteTermin(t *testing.T) {
	createNewTermin("testo", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	deleteTermin("testo", "Mik")
}

func TestGetTermine(t *testing.T) {
	createNewTermin("testo", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	i := createNewTermin("testu", "test", repeat(jaehrlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC", "Mik")
	k := getTermine("Mik")

	if k[3] != i {
		t.Errorf("Laden fehlgeschlagen")
	}

	for i := 0; i < len(k); i++ {
		fmt.Println(k[i])
	}

}
