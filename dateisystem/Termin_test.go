package dateisystem

import (
	"fmt"
	"testing"
	"time"
)

func newTermin(title string) *Termin {
	t := Termin{title: title}
	t.description = "test"
	t.recurring = repeat(niemals)
	t.date = time.Date(2021, 8, 15, 14, 30, 45, 0, time.UTC)
	t.endDate = time.Date(2021, 8, 15, 15, 30, 45, 0, time.UTC)
	return &t
}

func updateTermin(termin *Termin) {
	setTitle(termin, "testo")
	setDescription(termin, "testo")
	setRecurring(termin, repeat(woechentlich))
	setDate(termin, "2007-03-02T13:02:05 UTC")
	setEndeDate(termin, "2007-03-02T15:02:05 UTC")
}

func TestTermin(t *testing.T) {
	termin := newTermin("test")

	if termin.title == "" {
		t.Errorf("Titel nicht hinterlegt")
	}
	if termin.description == "" {
		t.Errorf("Beschreibung nicht hinterlegt")
	}
	if termin.recurring != 4 {
		t.Errorf("Wiederholung nicht hinterlegt")
	}
	if termin.date.String() == "" {
		t.Errorf("Datum nicht hinterlegt")
	}

	if termin.endDate.String() != "2021-08-15 15:30:45 +0000 UTC" {
		t.Errorf("End-Datum nicht hinterlegt")
	}
}

func TestTerminUpdate(t *testing.T) {
	termin := newTermin("test")
	updateTermin(termin)

	if termin.title == "test" {
		fmt.Println(termin.title)
		t.Errorf("Titel nicht aktualisiert")
	}
	if termin.description == "test" {
		fmt.Println(termin.description)
		t.Errorf("Beschreibung nicht aktualisiert")
	}
	if termin.recurring == repeat(niemals) {
		fmt.Println(termin.recurring)
		t.Errorf("Wiederholung nicht aktualisiert")
	}
	if termin.date.String() == "2021-08-15 14:30:45 +0000 UTC" {
		t.Errorf("Datum nicht aktualisiert")
	}
	if termin.endDate.String() == "2021-08-15 15:30:45 +0000 UTC" {
		t.Errorf("Datum nicht aktualisiert")
	}
}
