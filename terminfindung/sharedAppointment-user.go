package terminfindung

import (
	"errors"
	"strings"
)

// VoteForDay
// give terminID, userAdmin, name of the user who votes, the identifier of the day voted for
// (Hashed value of the time it was created) and how he votes
func VoteForDay(terminID *string, userAdmin *string, name *string, votedDay *string, voted bool) (err error) {
	//get the termin
	termin, err := GetTerminFromShared(userAdmin, terminID)
	if err != nil {
		return err
	}
	// check whether the user even exists
	if _, ok := termin.Persons[*name]; ok {
		//if he exists we can automatically vote --> it will not be existing
		//before because we dont add votes on create proposal date
		// so it is "undefined" and we can add one
		termin.Persons[*name].Votes[*votedDay] = voted
	} else {
		err = errors.New("cound't find username")
		return
	}
	// lock the allTermin structure for writing back and "saving" the data
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	allTermine.shared[*userAdmin+"|"+*terminID] = termin
	return
}

func GetTerminViaApiKey(apikey *string) (termin TerminFindung, user string, err error) {
	//check if apikey is even there
	if len(*apikey) == 0 {
		err = errors.New("API Key is not defined")
		return
	}
	// make read lock to prevent more read --> dont need to handle if it is locked before
	//because it is supposed to be called directly by the handler in routerSharedAppointemtn
	allTermine.mutex.RLock()
	defer allTermine.mutex.RUnlock()
	// if it exists it returns the user (not the applicant of the calendar system)
	//   returns user so it can be found by the template and used for links
	if val, ok := allTermine.links[*apikey]; ok {
		termin = allTermine.shared[val]
		//if the object is empty ---> if the user isnt set it automatically discards it
		if len(termin.User) == 0 {
			return TerminFindung{}, "", errors.New("Termin Object is empty")
		}
		for key, pers := range termin.Persons {
			if strings.Compare(pers.Url, *apikey) == 0 {
				user = key
				break
			}
		}
		err = nil
		return
	} else {
		err = errors.New("coudn't found API key")
		return
	}
}
