package terminfindung

import (
	"errors"
	"strings"
)

func VoteForDay(terminID *string, userAdmin *string, name *string, votedDay *string, voted bool) (err error) {
	termin, err := GetTerminFromShared(userAdmin, terminID)
	if err != nil {
		return err
	}
	if _, ok := termin.Persons[*name]; ok {
		termin.Persons[*name].Votes[*votedDay] = voted
	} else {
		err = errors.New("cound't find username")
		return
	}
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	allTermine.shared[*userAdmin+"|"+*terminID] = termin
	return
}

func GetTerminViaApiKey(apikey *string) (termin TerminFindung, user string, err error) {
	if len(*apikey) == 0 {
		err = errors.New("API Key is not defined")
		return
	}
	if val, ok := allTermine.links[*apikey]; ok {
		termin = allTermine.shared[val]
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
