package dateisystem

import (
	"testing"
)

func TestNewTerminObj(t *testing.T) {

	termin := newTerminObj("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC")

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

func TestCreateNewTermin(t *testing.T) {
	createNewTermin("test", "test", repeat(woechentlich), "2007-03-02T15:02:05 UTC", "2007-03-02T15:02:05 UTC")
}
