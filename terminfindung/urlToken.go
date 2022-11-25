package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"sync"
)

type User struct {
	Url       string
	Name      string
	Vorschlag map[string]bool
}
type TerminFindung struct {
	user             string
	info             dateisystem.Termin
	vorschlagTermine []dateisystem.Termin
	persons          []User
}

type termine struct {
	mutex  sync.RWMutex
	shared []TerminFindung
}

var allTermine = termine{shared: []TerminFindung{}}

func CreateURLToken(termin *dateisystem.Termin, user *string) {

}

func CreateSharedTermin(termin *dateisystem.Termin, user *string) string {
	allTermine.mutex.Lock()
	defer allTermine.mutex.Unlock()
	newShared := TerminFindung{
		user:             *user,
		info:             *termin,
		vorschlagTermine: []dateisystem.Termin{},
		persons:          []User{},
	}
	allTermine.shared = append(allTermine.shared, newShared)
	return termin.ID
}
