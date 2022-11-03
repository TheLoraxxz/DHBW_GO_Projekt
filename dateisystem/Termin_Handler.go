package dateisystem

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

func newTerminObj(title string, description string, rep repeat, date string, endDate string) Termin {

	dat, _ := time.Parse(dateLayoutISO, date)
	enddat, _ := time.Parse(dateLayoutISO, endDate)

	T := Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        dat,
		EndDate:     enddat}
	return T
}

func StoreNewTerminObj(termin Termin) {
	ter := termin
	file, _ := json.MarshalIndent(ter, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)
}

func createNewTermin(title string, description string, rep repeat, date string, endDate string) Termin {
	t := newTerminObj(title, description, rep, date, endDate)
	StoreNewTerminObj(t)
	return t
}
