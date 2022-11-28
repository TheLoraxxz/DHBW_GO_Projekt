package dateisystem

//ToDo Funktion Dateiname durch ID ersetzen
//Mat-Nr. 8689159
import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"strings"
	"time"
)

// NewTerminObj erzeugt einen transitiven Termin; NUR FÜR TESTS EMPFOHLEN
func NewTerminObj(title string, description string, rep Repeat, date time.Time, endDate time.Time) Termin {

	t := Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        date,
		EndDate:     endDate,
		ID:          createID(date, endDate)}

	return t
}

// AddToCache fügt Termin dem Caching hinzu
func AddToCache(termin Termin, kalender []Termin) []Termin {
	k := append(kalender, termin)
	return k
}

// StoreTerminObj exportiert Termine zu json, "username" mapped Termine und Nutzer
func StoreTerminObj(termin Termin, username string) {

	path := getFile(termin.ID, username)

	p, _ := json.MarshalIndent(termin, "", " ") //erzeugt die json
	_ = os.WriteFile(path, p, 0755)             //schreibt json in Datei
}

// CreateNewTermin erzeugt einen persistenten Termin
func CreateNewTermin(title string, description string, rep Repeat, date time.Time, endDate time.Time, username string) Termin {
	t := NewTerminObj(title, description, rep, date, endDate)
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

	files, err2 := f.Readdir(0) //liest alle Dateinamen ein
	if err2 != nil {
		fmt.Println(err2)
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
	var kNew []Termin
	file := getFile(kalender[0].ID, username)

	for i := 0; i < len(kalender); i++ {
		if kalender[i].ID != id {
			kNew = append(kNew, kalender[i]) //kopiert in neuen Kalender
		} else { //prüft, ob der Termin persistent, oder transitiv ist
			if _, err := os.Stat(file); err == nil {
				deleteTermin(id, username)
			}
		}
	}

	return kNew
}

// createID erzeugt neue ID
func createID(dat time.Time, endDat time.Time) string {

	u := time.Now().String()

	id := dat.String() + endDat.String() + u

	bytes, err := bcrypt.GenerateFromPassword([]byte(id), 14)
	id = string(bytes)

	f, err := os.OpenFile("id.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	id = strings.Replace(id, "/", "E", 99)
	id = strings.Replace(id, ".", "D", 99)

	if _, err = f.WriteString(id + "\n"); err != nil {
		panic(err)
	}

	return id
}
