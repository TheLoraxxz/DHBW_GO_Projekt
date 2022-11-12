package dateisystem

//Mat-Nr. 8689159
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func NewTerminObj(title string, description string, rep Repeat, date time.Time, endDate time.Time) Termin { //erzeugt einen transitiven Termin; NUR FÜR TESTS EMPFOHLEN

	//dat, _ := time.Parse(dateLayoutISO, date)
	//enddat, _ := time.Parse(dateLayoutISO, endDate)

	T := Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        date,
		EndDate:     endDate}
	return T
}

func CreateNewTermin(title string, description string, rep Repeat, date time.Time, endDate time.Time, username string) Termin { //erzeugt einen persistenten Termin
	t := NewTerminObj(title, description, rep, date, endDate)
	StoreTerminObj(t, username)
	return t
}

// liefert slice mit allen terminen eines Users zurück
func GetTermine(username string) []Termin { //liefert slice mit allen terminen eines Users zurück
	var k []Termin

	path := username //öffnet das Verzeichnis des Users
	f, err := os.Open(path)
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
		title := strings.Split(v.Name(), ".")
		k = append(k, LoadTermin(title[0], username))
	}
	defer func(f *os.File) { //schließt Pointer auf die json
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	return k
}

func StoreTerminObj(termin Termin, username string) { //exportiert Termine zu json, "username" mapped Termine und Nutzer
	file := username + "/" + termin.Title + ".json"
	ter := termin

	err := os.MkdirAll(username, 755) //erzeugt Verzeichnis passend zum User, falls noch nicht existent
	if err != nil {
		fmt.Println(err)
	}

	p, _ := json.MarshalIndent(ter, "", " ") //erzeugt die json
	_ = os.WriteFile(file, p, 0755)          //schreibt json in Datei
}

func AddToCache(termin Termin, kalender []Termin) []Termin { //fügt Termin dem Caching hinzu
	k := append(kalender, termin)
	return k
}

func StoreCache(kalender []Termin, username string) { //speichert alle Elemente Caches von User "username"
	k := kalender

	for i := 0; i < len(k); i++ {
		StoreTerminObj(k[i], username)
	}
}

func LoadTermin(title string, username string) Termin { //kreiert Termin aus json, "username" mapped Termine und Nutzer
	file := username + "/" + title + ".json"

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

func deleteTermin(title string, username string) { //löscht json mit den Informationen zum Termin, "username" mapped Termine und Nutzer
	file := username + "/" + title + ".json"
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteAll(kalender []Termin, username string) []Termin { //löscht alle Termine eines Users, liefert []Termin(nil) zurück
	k := kalender

	for i := 0; i < len(k); i++ {
		deleteTermin(k[i].Title, username)
	}

	k = GetTermine(username)
	return k
}

func DeleteFromCache(kalender []Termin, title string, username string) []Termin { //Löscht einzelnes Element aus dem Cache
	kOld := kalender
	var kNew []Termin
	file := username + "/" + title + ".json"

	for i := 0; i < len(kOld); i++ {
		if kOld[i].Title != title {
			kNew = append(kNew, kOld[i])
		} else { //prüft, ob der Termin persistent, oder transitiv ist
			if _, err := os.Stat(file); err == nil {
				deleteTermin(title, username)
			}
		}
	}

	return kNew
}
