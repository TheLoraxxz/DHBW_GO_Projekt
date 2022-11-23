package terminfindung

import "DHBW_GO_Projekt/dateisystem"

type User struct {
	Name      string
	Vorschlag map[int64]bool
}
type TerminFindung struct {
	Url              string
	vorschlagTermine []Vorschlag
	persons          []User
}
type Vorschlag struct {
	id     int64
	termin dateisystem.Termin
}

// make(map[string]Id)
