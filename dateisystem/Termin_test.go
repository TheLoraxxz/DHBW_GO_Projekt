package dateisystem

import (
	"testing"
)

func newTermin(title string) *Termin {
	t := Termin{title: title}
	t.description = "test"
	t.recurring = repeat(niemals)
	t.time = "15:04"
	t.date = "2022-10-03"
	return &t
}

func TestTermin(t *testing.T) {
	termin := newTermin("test")

	if termin.title == "" {
		t.Errorf("Titel nicht hinterlegt")
	}
	if termin.date == "" {
		t.Errorf("Datum nicht hinterlegt")
	}
	if termin.time == "" {
		t.Errorf("Zeit nicht hinterlegt")
	}
	if termin.description == "" {
		t.Errorf("Beschreibung nicht hinterlegt")
	}
	if termin.recurring != 4 {
		t.Errorf("Wiederholung nicht hinterlegt")
	}
}
