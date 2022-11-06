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

func CreateNewTermin(title string, description string, rep repeat, date string, endDate string, username string) Termin { //erzeugt einen persistenten Termin
	t := NewTerminObj(title, description, rep, date, endDate)
	StoreTerminObj(t, username)
	return t
}

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

	for _, v := range files { // lädt der Reihe nach alle Dateien ein
		title := strings.Split(v.Name(), ".")
		k = append(k, LoadTermin(title[0], username))
	}
	defer func(f *os.File) { // schließt Pointer auf die Json
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	return k
}

func NewTerminObj(title string, description string, rep repeat, date string, endDate string) Termin { //erzeugt einen transitiven Termin

	dat, _ := time.Parse(dateLayoutISO, date)
	enddat, _ := time.Parse(dateLayoutISO, endDate)

	T := Termin{
		Title:       title,
		Description: description,
		Recurring:   rep,
		Date:        dat,
		EndDate:     enddat}
	return T
}

func AddToCache(termin Termin, kalender []Termin) []Termin { // fügt Termin dem Caching hinzu
	k := append(kalender, termin)
	return k
}

func StoreTerminObj(termin Termin, username string) { //exportiert Termine zu json, "username" mapped Termine und Nutzer
	file := username + "/" + termin.Title + ".json"
	ter := termin

	err := os.MkdirAll(username, 755)
	if err != nil {
		fmt.Println(err)
	}

	p, _ := json.MarshalIndent(ter, "", " ")
	_ = os.WriteFile(file, p, 0755)
}

func LoadTermin(tittle string, username string) Termin { //kreiert Termin aus json, "username" mapped Termine und Nutzer
	file := username + "/" + tittle + ".json"

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
	defer func(f *os.File) { // schließt Pointer auf die Json
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(open)
	return t
}

func DeleteTermin(tittle string, username string) { //löscht json mit den Informationen zum Termin, "username" mapped Termine und Nutzer
	file := username + "/" + tittle + ".json"
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}
