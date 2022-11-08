package kalenderansicht

import (
	"DHBW_GO_Projekt/dateisystem"
	"fmt"
	"net/http"
)

type kalenderEinträgeGenerator interface {
	ErstelleKalenderEintraege(userId string)
}

func CreateTermin(r *http.Request, username string) {
	title := r.FormValue("titel")
	description := r.FormValue("beschreibung")
	repStr := r.FormValue("wiederholung")
	date := r.FormValue("datum")
	endDate := r.FormValue("endDatum")
	var rep dateisystem.Repeat
	switch repStr {
	case "täglich":
		rep = 0
	case "wöchentlich":
		rep = 1
	case "jährlich":
		rep = 2
	case "niemals":
		rep = 3
	}
	fmt.Println(dateisystem.CreateNewTermin(title, description, rep, date, endDate, username))
}
