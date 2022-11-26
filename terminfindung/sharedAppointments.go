package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"errors"
	"sync"
	"time"
)

type User struct {
	Url       string
	Vorschlag map[string]bool
}
type TerminFindung struct {
	user             string
	info             dateisystem.Termin
	VorschlagTermine []dateisystem.Termin
	persons          map[string]User
}

// Shared
// mutex --> mutex to lock the thing
// shared --> a map of TerminFindung --> the string is user|idOfTheTermin --> is for terminfindung admin
// links --> all links for terminfindung user --> which are the one who dont have an id --> if loginrequest fails it checks
// the api key + user and then redirects to the actual clientwebsite --> links consists of personInvited|apikey and refers to sharerd
// string
type Shared struct {
	mutex  sync.RWMutex
	shared map[string]TerminFindung
	links  map[string]string
}

var allTermine = Shared{
	shared: make(map[string]TerminFindung, 5),
	links:  make(map[string]string, 10),
}

func GetTerminFromShared(user *string, terminID *string) TerminFindung {
	if len(*user) == 0 && len(*terminID) == 0 {
		return TerminFindung{}
	}
	return allTermine.shared[*user+"|"+*terminID]
}

func CreateSharedTermin(termin *dateisystem.Termin, user *string) (uuid string) {
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	newTermin := TerminFindung{
		user:             *user,
		info:             *termin,
		VorschlagTermine: []dateisystem.Termin{*termin},
		persons:          make(map[string]User, 10),
	}
	allTermine.shared[*user+"|"+termin.ID] = newTermin
	return *user + "|" + termin.ID
}
func CreateNewProposedDate(startdate time.Time, endDate time.Time, userAppID *string) error {
	newProposedTermin := dateisystem.Termin{
		Date:    startdate,
		EndDate: endDate,
	}
	if startdate.After(endDate) {
		return errors.New("Can't insert startdate which has the wrong format")
	}
	if _, ok := allTermine.shared[*userAppID]; !ok {
		return errors.New("User ID is not definied")
	}
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	termin := allTermine.shared[*userAppID]
	termin.VorschlagTermine = append(termin.VorschlagTermine, newProposedTermin)
	allTermine.shared[*userAppID] = termin
	return nil
}
