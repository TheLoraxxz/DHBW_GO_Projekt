package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strings"
	"sync"
	"time"
)

type UserTermin struct {
	Url   string
	Name  string
	Votes map[string]bool
}
type TerminFindung struct {
	User             string
	Info             dateisystem.Termin
	FinalTermin      dateisystem.Termin
	VorschlagTermine []dateisystem.Termin
	Persons          map[string]UserTermin
}

// Shared
// mutex --> mutex to lock the thing
// shared --> a map of TerminFindung --> the string is user|idOfTheTermin --> is for terminfindung admin
// links --> all links for terminfindung user --> which are the one who don't have an id --> if loginrequest fails it checks
// the api key + user and then redirects to the actual clientwebsite --> links consists of personInvited|apikey and refers to sharerd
// string
type Shared struct {
	mutex  sync.RWMutex
	shared map[string]TerminFindung
	links  map[string]string
}

// implementation of Shared
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
		User:             *user,
		Info:             *termin,
		Persons:          make(map[string]UserTermin, 10),
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

// CreateNewProposedDate
// alreadyLocked per default should be false
func CreateNewProposedDate(startdate time.Time, endDate time.Time, user *string, terminID *string, alreadyLocked bool) (err error) {
	idTermin, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()+startdate.String()), 1)
	newProposedTermin := dateisystem.Termin{
		Date:    startdate,
		EndDate: endDate,
		ID:      string(idTermin),
	}
	if startdate.After(endDate) {
		return errors.New("can't insert startdate which has the wrong format")
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

func CreatePerson(name *string, terminID *string, user *string) (urlToShow string, err error) {
	//checks whether the input parameters are right
	if len(*name) == 0 || len(*terminID) == 0 || len(*user) == 0 {
		err = errors.New("name, TerminID and user need to be set")
	}
	bytesHash, err := bcrypt.GenerateFromPassword([]byte(*name+*user+url.QueryEscape(*terminID)), 1)
	newUser := UserTermin{
		Name:  *name,
		Url:   url.QueryEscape(string(bytesHash)),
		Votes: map[string]bool{},
	}
	if err != nil {
		return
	}
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	termin, err := GetTerminFromShared(user, terminID)
	for key := range termin.Persons {
		if strings.Compare(key, *name) == 0 {
			err = errors.New("user already existed in Termin")
			return
		}
	}
	if err != nil {
		return
	}
	termin.Persons[*name] = newUser
	urlToShow = "terminID=" + url.QueryEscape(*terminID) + "&name=" + *name + "&user=" + *user + "&apiKey=" + url.QueryEscape(string(bytesHash))
	allTermine.links[urlToShow] = *user + "|" + *terminID
	return urlToShow, nil
}

func GetAllLinks(user *string, terminId *string) (users []UserTermin, err error) {
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	termin, err := GetTerminFromShared(user, terminId)
	if err != nil {
		return
	}
	//add all the user that are there
	for _, element := range termin.Persons {
		users = append(users, element)
	}
	return
}

func SelectDate(idPropDate *string, terminID *string, user *string) (err error) {
	termin, err := GetTerminFromShared(user, terminID)
	if err != nil {
		return err
	}
	for _, elem := range termin.VorschlagTermine {
		if elem.ID == *idPropDate {
			termin.FinalTermin = elem
		}
	}
	return nil
}
