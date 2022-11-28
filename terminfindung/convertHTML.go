package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"strings"
)

type AdminTerminInfo struct {
	Termin       dateisystem.Termin
	VotesFor     int64
	PersonsVoted []string
	IsSelected   bool
}

type AdminHTML struct {
	TerminID         string
	Info             dateisystem.Termin
	VorschlagTermine []AdminTerminInfo
	Persons          []string
	IsLocked         bool
}

func (t TerminFindung) ConvertAdminToHTML() (rightHTML AdminHTML) {
	// create correct HTML
	rightHTML = AdminHTML{
		TerminID:         t.Info.ID,
		Info:             t.Info,
		Persons:          []string{},
		VorschlagTermine: []AdminTerminInfo{},
		IsLocked:         false,
	}
	//make for to get all the termine
	for _, elem := range t.VorschlagTermine {
		newTermin := AdminTerminInfo{
			Termin:     elem,
			IsSelected: false,
		}
		//if it is the final termin so it highlights it
		if strings.Compare(elem.ID, t.FinalTermin.ID) == 0 {
			//if the final termin is set it copys it to the info so
			//the right information stands then on the main page
			newTermin.IsSelected = true
			rightHTML.IsLocked = true
			rightHTML.Info.Date = t.FinalTermin.Date
			rightHTML.Info.EndDate = t.FinalTermin.EndDate
			rightHTML.Info.Description = t.FinalTermin.Description

		}
		rightHTML.VorschlagTermine = append(rightHTML.VorschlagTermine, newTermin)
	}
	for person, elem := range t.Persons {
		//append all persons including there names
		rightHTML.Persons = append(rightHTML.Persons, person)
		for id, isVoted := range elem.Votes {
			if isVoted {
				for i, search := range rightHTML.VorschlagTermine {
					if search.Termin.ID == id {
						rightHTML.VorschlagTermine[i].VotesFor += 1
						rightHTML.VorschlagTermine[i].PersonsVoted = append(rightHTML.VorschlagTermine[i].PersonsVoted, person)
					}
				}
			}
		}
	}
	return
}
