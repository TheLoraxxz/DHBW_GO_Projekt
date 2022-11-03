package dateisystem

import (
	time "time"
)

type repeat int

const (
	taeglich repeat = iota
	woechentlich
	monatlich
	jaehrlich
	niemals
)

const (
	dateLayoutISO = "2006-01-02T15:04:05 UTC"
)

type Termin struct {
	Title       string
	Description string
	Recurring   repeat
	Date        time.Time
	EndDate     time.Time
}

func setTitle(t *Termin, newTitle string) {
	t.Title = newTitle
}

func setDescription(t *Termin, newDescription string) {
	t.Description = newDescription
}

func setRecurring(t *Termin, newRecurring repeat) {
	t.Recurring = newRecurring
}

func setDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.Date = d
}

func setEndeDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.EndDate = d
}
