package dateisystem

//ToDo Funktion Dateiname durch ID ersetzen
//Mat-Nr. 8689159
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// NewTerminObj erzeugt einen transitiven Termin; NUR FÜR TESTS EMPFOHLEN
func NewTerminObj(title string, description string, rep Repeat, date time.Time, endDate time.Time, id string) Termin {
	T := Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        date,
		EndDate:     endDate,
		ID:          id}
	return T
}

// CreateNewTermin erzeugt einen persistenten Termin
func CreateNewTermin(title string, description string, rep Repeat, date time.Time, endDate time.Time, username string, id string) Termin {
	t := NewTerminObj(title, description, rep, date, endDate, id)
	StoreTerminObj(t, username)
	return t
}

// GetTermine liefert slice mit allen terminen eines Users zurück
func GetTermine(username string) []Termin {
	var k []Termin

	path := GetDirectory(username)

	f, err := os.Open(path) //öffnet das Verzeichnis des Users
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files, err := f.Readdir(0) //liest alle Dateinamen ein
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, v := range files { //lädt der Reihe nach alle Dateien ein
		id := strings.Split(v.Name(), ".")
		k = append(k, LoadTermin(id[0], username))
	}
	defer func(f *os.File) { //schließt Pointer auf die json
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	return k
}

// StoreTerminObj exportiert Termine zu json, "username" mapped Termine und Nutzer
func StoreTerminObj(termin Termin, username string) {
	ter := termin
	count, _ := os.ReadDir(GetDirectory(username))
	i := len(count) + 1
	id := strconv.Itoa(i)

	ter.setID(&ter, id)

	path := getFile(id, username)

	p, _ := json.MarshalIndent(ter, "", " ") //erzeugt die json
	_ = os.WriteFile(path, p, 0755)          //schreibt json in Datei
}

// AddToCache fügt Termin dem Caching hinzu
func AddToCache(termin Termin, kalender []Termin) []Termin {
	termin.setID(&termin, strconv.Itoa(len(kalender)+1))
	k := append(kalender, termin)
	return k
}

// StoreCache speichert alle Elemente Caches von User "username"
func StoreCache(kalender []Termin, username string) {
	k := kalender

	for i := 0; i < len(k); i++ {
		StoreTerminObj(k[i], username)
	}
}

// LoadTermin kreiert Termin aus json, "username" mapped Termine und Nutzer
func LoadTermin(id string, username string) Termin {
	file := getFile(id, username)

	open, err := os.Open(file) //öffnet json
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := io.ReadAll(open)

	var t Termin //erzeugt Termin aus json
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		fmt.Println(err)
	}
	defer func(f *os.File) { //schließt Pointer auf die json
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(open)
	return t
}

// deleteTermin löscht json mit den Informationen zum Termin, "username" mapped Termine und Nutzer
func deleteTermin(id string, username string) {
	file := getFile(id, username)
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}

// DeleteAll löscht alle Termine eines Users, liefert []Termin(nil) zurück
func DeleteAll(kalender []Termin, username string) []Termin {
	k := kalender

	for i := 0; i < len(k); i++ {
		deleteTermin(k[i].ID, username)
	}

	k = GetTermine(username)
	return k
}

// DeleteFromCache löscht einzelnes Element aus dem Cache und ggf. die dazugehörige json
func DeleteFromCache(kalender []Termin, id string, username string) []Termin {
	kOld := kalender
	var kNew []Termin
	file := getFile(kalender[0].ID, username)

	for i := 0; i < len(kOld); i++ {
		if kOld[i].ID != id {
			kNew = append(kNew, kOld[i])
		} else { //prüft, ob der Termin persistent, oder transitiv ist
			if _, err := os.Stat(file); err == nil {
				deleteTermin(id, username)
			}
		}
	}

	return kNew
}
