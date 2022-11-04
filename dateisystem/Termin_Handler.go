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

func createNewTermin(title string, description string, rep repeat, date string, endDate string, username string) Termin { //erzeugt ein persistenten Termin
	t := newTerminObj(title, description, rep, date, endDate)
	storeTerminObj(t, username)
	return t
}

func getTermine(username string) []Termin { //liefert slice mit allen terminen eines Users zurück
	var k []Termin

	path := "./" + username //öffnet das  Verzeichnis des Users
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	files, err := f.Readdir(0) //ließt alle Dateinamen ein
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, v := range files { // lädt der Reihe nach alle Dateien ein
		title := strings.Split(v.Name(), ".")
		k = append(k, loadTermin(title[0], username))
	}
	return k
}

func newTerminObj(title string, description string, rep repeat, date string, endDate string) Termin { //erzeug einen transistiven Termin

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

func storeTerminObj(termin Termin, username string) { //exportiert Termine zu json, "username" mapped Termine und Nutzer
	file := "./" + username + "/" + termin.Title + ".json"
	ter := termin

	p, _ := json.MarshalIndent(ter, "", " ")
	_ = os.WriteFile(file, p, 0755)
}

func deleteTermin(tittle string, username string) { //löscht json mit den Infos zum Termin, "username" mapped Termine und Nutzer
	file := "./" + username + "/" + tittle + ".json"
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}

func loadTermin(tittle string, username string) Termin { //kreirt Termin aus json, "username" mapped Termine und Nutzer
	file := "./" + username + "/" + tittle + ".json"

	jsonFile, err := os.Open(file) //öffnet json
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var t Termin //erzeugt Termin aus json
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		fmt.Println(err)
	}
	return t
}
