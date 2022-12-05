/*
@author: 2447899 8689159 3000685
*/
package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type VoteTermin struct {
	Votedfor int
	Term     dateisystem.Termin
}

type UserHTML struct {
	User    string
	APIKey  string
	Info    dateisystem.Termin
	ToVotes map[string]VoteTermin
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

// ConvertUserSiteToRightHTML --> used for converting user site to html object without weird objects in between
// user --> the name of the user (not admin-user)
func (termin TerminFindung) ConvertUserSiteToRightHTML(user *string, apikey *string) (newHTMLObj UserHTML) {
	if len(*user) == 0 {
		return
	}
	newHTMLObj = UserHTML{
		Info:    termin.Info,
		User:    *user,
		ToVotes: map[string]VoteTermin{},
		APIKey:  *apikey,
	}
	for _, elem := range termin.VorschlagTermine {
		voteTermin := VoteTermin{
			Votedfor: 0,
			Term:     elem,
		}
		newHTMLObj.ToVotes[elem.ID] = voteTermin
	}
	for vote, bol := range termin.Persons[*user].Votes {
		voted := newHTMLObj.ToVotes[vote]
		if bol {
			voted.Votedfor = 1
		} else {
			voted.Votedfor = -1
		}
		newHTMLObj.ToVotes[vote] = voted
	}
	return
}

// SaveSharedTermineToDisk --> saving it to the disk
func SaveSharedTermineToDisk(basepath *string) error {
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	// thorw out all empty terminfindung objects --> happens if one deletes a shared termin
	newTermine := map[string]TerminFindung{}
	for key, val := range allTermine.shared {
		if !(val.User == "") {
			newTermine[key] = val
		}
	}
	// set to ../data/shared-termin/shared-termin-data.json
	allTermine.shared = newTermine
	file, err := json.MarshalIndent(newTermine, "", " ")
	if err != nil {
		return fmt.Errorf("coudn't Convert to JSON data - Error: %w", err)
	}
	pathAbs := filepath.Join(*basepath, "data", "shared-termin", "shared-termin-data.json")
	err = os.WriteFile(pathAbs, file, 0644)
	if err != nil {
		return fmt.Errorf("Coudn't write JSON to path - Error: %w", err)
	}

	// save links in seperate JSON File
	file, err = json.MarshalIndent(allTermine.links, "", " ")
	if err != nil {
		return fmt.Errorf("coudn't Convert links to JSON: %w", err)
	}
	pathAbs = filepath.Join(*basepath, "data", "shared-termin", "links.json")
	err = os.WriteFile(pathAbs, file, 0644)
	if err != nil {
		return fmt.Errorf("coudn't write to file links to JSON: %w", err)
	}
	return nil

}

// LoadDataToSharedTermin
// loads the data from ../data/shared-termin/ and puts it into the shared muted
func LoadDataToSharedTermin(pathBase *string) (err error) {
	//lock file directly to reset the mutex
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	allTermine.shared = map[string]TerminFindung{}
	allTermine.links = map[string]string{}
	// read paths and read it
	pathShared := filepath.Join(*pathBase, "data", "shared-termin", "shared-termin-data.json")
	pathLinks := filepath.Join(*pathBase, "data", "shared-termin", "links.json")
	fileShared, err := os.ReadFile(pathShared)
	fileLinks, err := os.ReadFile(pathLinks)
	if err != nil {
		return fmt.Errorf("coudn't read files %w", err)
	}
	err = json.Unmarshal(fileShared, &allTermine.shared)
	err = json.Unmarshal(fileLinks, &allTermine.links)
	if err != nil {
		allTermine.shared = map[string]TerminFindung{}
		allTermine.links = map[string]string{}
		return fmt.Errorf("Coudn't Convert Jsons %w", err)
	}
	return nil
}
