package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/url"
	"strconv"
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

// GetTerminFromShared
// get termin from mutex and check if it exists
func GetTerminFromShared(user *string, terminID *string) (termin TerminFindung, err error) {
	// plausibilitäts check
	if len(*user) == 0 && len(*terminID) == 0 {
		err = errors.New("termin id and user is not valid")
		return
	}
	userAppID := *user + "|" + *terminID
	// checks that it finds it and the obj is not empty
	if _, ok := allTermine.shared[userAppID]; !ok || len(allTermine.shared[userAppID].User) == 0 {
		err = errors.New("can't find SharedTermin")
		return
	}
	return allTermine.shared[userAppID], nil
}

// CreateSharedTermin
// creates a shared termin and also create a proposed date
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
	// write back to termin
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
	// generate easy hash to create an ID that is unique
	idTermin, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()+startdate.String()), 1)
	newProposedTermin := dateisystem.Termin{
		Date:    startdate,
		EndDate: endDate,
		ID:      string(idTermin),
	}
	// check if the date connection is right
	if startdate.After(endDate) {
		return errors.New("can't insert startdate which has the wrong format")
	}
	// is because when it is called from CreateSharedTermin it is already locked --> would be deadlock
	if !alreadyLocked {
		allTermine.mutex.Lock()
		defer allTermine.mutex.Unlock()
	}
	// gertermin and update the object
	termin, err := GetTerminFromShared(user, terminID)
	if err != nil {
		return
	}
	termin.VorschlagTermine = append(termin.VorschlagTermine, newProposedTermin)
	allTermine.shared[*user+"|"+*terminID] = termin
	return nil
}

// CreatePerson
// creates a user taht is beeing surfed
// user --> the amdin user that created the shared appointment
// urltoshow is the full url with apikey --> is unieque to the termin and user
func CreatePerson(name *string, terminID *string, user *string) (urlToShow string, err error) {
	//checks whether the input parameters are right
	if len(*name) == 0 || len(*terminID) == 0 || len(*user) == 0 {
		err = errors.New("name, TerminID and user need to be set")
	}
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	termin, err := GetTerminFromShared(user, terminID)
	//check if name for shared user already exsists
	if len(termin.Persons[*name].Url) > 0 {
		return "", errors.New("user already exists")
	}
	// tupel beeing created from the username, the admin user (for the case that there are already two and the temrin id)
	bytesHash, err := bcrypt.GenerateFromPassword([]byte(*name+*user+url.QueryEscape(*terminID)), 1)
	newUser := UserTermin{
		Name:  *name,
		Url:   url.QueryEscape(string(bytesHash)),
		Votes: map[string]bool{},
	}
	if err != nil {
		return
	}
	//overwrite current object with new object to update the map
	termin.Persons[*name] = newUser
	allTermine.shared[*user+"|"+*terminID] = termin
	urlToShow = "apiKey=" + url.QueryEscape(string(bytesHash))
	allTermine.links[url.QueryEscape(string(bytesHash))] = *user + "|" + *terminID
	return urlToShow, nil
}

// GetAllLinks is for the Alle Links anzeigen page --> returns all users
func GetAllLinks(user *string, terminId *string) (users []UserTermin, err error) {
	allTermine.mutex.RLock()
	defer allTermine.mutex.RUnlock()
	//reads from a termin all persons
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

// SelectDate
// is the function that meets requirements 8.9 - 8.12
func SelectDate(idPropDate *string, terminID *string, user *string) (err error) {
	//get the termin
	termin, err := GetTerminFromShared(user, terminID)
	if err != nil {
		return err
	}
	// check where the final termin is and put it into the final termin
	for _, elem := range termin.VorschlagTermine {
		if strings.Compare(elem.ID, *idPropDate) == 0 {
			termin.FinalTermin = elem
		}
	}
	// count the votes for this termin and the one against it
	votedFor := 0
	for _, elem := range termin.Persons {
		if val, ok := elem.Votes[*idPropDate]; ok && val {
			votedFor++
		}
	}
	// tage over the recurring and the title
	termin.FinalTermin.Title = termin.Info.Title
	termin.FinalTermin.Recurring = termin.Info.Recurring
	//change description to the number who voted --> directly voted
	termin.FinalTermin.Description = termin.FinalTermin.Description + "| Dafür gestimmt: " +
		strconv.Itoa(votedFor) + " / Dagegen oder enthalten: " + strconv.Itoa(len(termin.Persons)-votedFor)
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	// shared Termine being saved back
	allTermine.shared[*user+"|"+*terminID] = termin
	// manage dateisystem and change the termin
	kalender := dateisystem.GetTermine(*user)
	kalender = dateisystem.DeleteFromCache(kalender, termin.Info.ID, *user)
	dateisystem.CreateNewTermin(termin.FinalTermin.Title, termin.FinalTermin.Description, dateisystem.Never,
		termin.FinalTermin.Date, termin.FinalTermin.EndDate, false, termin.User)
	return nil
}

func DeleteSharedTermin(terminID *string, user *string) (err error) {
	if len(*terminID) == 0 || len(*user) == 0 {
		return errors.New("terminid and user needs to be set")
	}
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	// shared Termine being saved back
	userAppID := *user + "|" + *terminID
	// checks that it finds it and the obj is not empty
	if _, ok := allTermine.shared[userAppID]; !ok || len(allTermine.shared[userAppID].User) == 0 {
		err = errors.New("can't find SharedTermin")
		return
	}
	allTermine.shared[*user+"|"+*terminID] = TerminFindung{}
	return nil
}
