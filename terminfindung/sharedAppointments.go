package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"sync"
)

type User struct {
	Url       string
	Vorschlag map[string]bool
}
type TerminFindung struct {
	user             string
	info             dateisystem.Termin
	vorschlagTermine []dateisystem.Termin
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

func CreateSharedTermin(termin *dateisystem.Termin, user *string) string {
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	newTermin := TerminFindung{
		user:             *user,
		info:             *termin,
		vorschlagTermine: []dateisystem.Termin{*termin},
		persons:          make(map[string]User, 10),
	}
	allTermine.shared[*user+"|"+termin.ID] = newTermin
	return termin.ID
}
