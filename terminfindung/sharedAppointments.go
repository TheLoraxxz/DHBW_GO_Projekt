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

func GetTerminFromShared(user *string, terminID *string) (termin TerminFindung, err error) {
	if len(*user) == 0 && len(*terminID) == 0 {
		err = errors.New("termin id and user is not valid")
		return
	}
	userAppID := *user + "|" + *terminID
	if _, ok := allTermine.shared[userAppID]; !ok {
		err = errors.New("can't find SharedTermin")
		return
	}
	return allTermine.shared[userAppID], nil
}

func CreateSharedTermin(termin *dateisystem.Termin, user *string) (uuid string, err error) {
	newTermin := TerminFindung{
		user:             *user,
		info:             *termin,
		persons:          make(map[string]User, 10),
		VorschlagTermine: []dateisystem.Termin{},
	}
	if len(*user) == 0 || len(termin.ID) == 0 {
		err = errors.New("UserID or Termin id isn't zero")
		return
	}
	terminProp := *user + "|" + termin.ID
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	allTermine.shared[terminProp] = newTermin
	err = CreateNewProposedDate(termin.Date, termin.EndDate, user, &termin.ID, true)
	if err != nil {
		return
	}
	return termin.ID, nil
}
func CreateNewProposedDate(startdate time.Time, endDate time.Time, user *string, terminID *string, alreadyLocked bool) (err error) {
	newProposedTermin := dateisystem.Termin{
		Date:    startdate,
		EndDate: endDate,
		ID:      time.Now().String(),
	}
	if startdate.After(endDate) {
		return errors.New("Can't insert startdate which has the wrong format")
	}
	if !alreadyLocked {
		allTermine.mutex.Lock()
		defer allTermine.mutex.Unlock()
	}
	termin, err := GetTerminFromShared(user, terminID)
	if err != nil {
		return
	}
	termin.VorschlagTermine = append(termin.VorschlagTermine, newProposedTermin)
	allTermine.shared[*user+"|"+*terminID] = termin
	return nil
}
