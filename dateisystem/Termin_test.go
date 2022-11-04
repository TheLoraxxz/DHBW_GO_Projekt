package dateisystem

//Mat-Nr. 8689159
import (
	"fmt"
	"testing"
	"time"
)

func newTermin(title string) *Termin { //erzeugt Pointer auf dummy Termin
	t := Termin{Title: title}
	t.Description = "test"
	t.Recurring = repeat(niemals)
	t.Date = time.Date(2021, 8, 15, 14, 30, 45, 0, time.UTC)
	t.EndDate = time.Date(2021, 8, 15, 15, 30, 45, 0, time.UTC)
	return &t
}

func updateTermin(termin *Termin) { //führt die setter aus
	setTitle(termin, "testj")
	setDescription(termin, "testo yeet")
	setRecurring(termin, repeat(woechentlich))
	setDate(termin, "2007-03-02T14:02:05 UTC")
	setEndeDate(termin, "2007-03-02T15:02:05 UTC")
}

func TestTermin(t *testing.T) { //prüft ob der dummy Termin nicht Leer ist
	termin := newTermin("test")

	if termin.Title == "" {
		t.Errorf("Titel nicht hinterlegt")
	}
	if termin.Description == "" {
		t.Errorf("Beschreibung nicht hinterlegt")
	}
	if termin.Recurring != 4 {
		t.Errorf("Wiederholung nicht hinterlegt")
	}
	if termin.Date.String() == "" {
		t.Errorf("Datum nicht hinterlegt")
	}

	if termin.EndDate.String() != "2021-08-15 15:30:45 +0000 UTC" {
		t.Errorf("End-Datum nicht hinterlegt")
	}
}

func TestTerminUpdate(t *testing.T) { // prüft ob die updates durchgeführt werden
	termin := newTermin("test")
	updateTermin(termin) //ruft die setter auf

	if termin.Title == "test" {
		fmt.Println(termin.Title)
		t.Errorf("Titel nicht aktualisiert")
	}
	if termin.Description == "test" {
		fmt.Println(termin.Description)
		t.Errorf("Beschreibung nicht aktualisiert")
	}
	if termin.Recurring == repeat(niemals) {
		fmt.Println(termin.Recurring)
		t.Errorf("Wiederholung nicht aktualisiert")
	}
	if termin.Date.String() == "2007-03-02T14:02:05 UTC" {
		t.Errorf("Datum nicht aktualisiert")
	}
	if termin.EndDate.String() == "2007-03-02T15:02:05 UTC" {
		t.Errorf("Datum nicht aktualisiert")
	}
}
