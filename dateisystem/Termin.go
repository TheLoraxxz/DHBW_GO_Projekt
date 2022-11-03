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
	title       string
	description string
	recurring   repeat
	date        time.Time
	endDate     time.Time
}

func setTitle(t *Termin, newTitle string) {
	t.title = newTitle
}

func setDescription(t *Termin, newDescription string) {
	t.description = newDescription
}

func setRecurring(t *Termin, newRecurring repeat) {
	t.recurring = newRecurring
}

func setDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.date = d
}

func setEndeDate(t *Termin, newDate string) {
	d, _ := time.Parse(dateLayoutISO, newDate)
	t.endDate = d
}
