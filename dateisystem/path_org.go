package dateisystem

/*
Zweck: für das Dateisystem relative Pfade bereit zustellen um zu garantieren, das aus anderen Modulen korrekt auf Dateien zugegriffen werden kann
ToDo Dry ausmerzen
*/

//Mat-Nr. 8689159
import (
	"fmt"
	"os"
	"path/filepath"
)

// getDirectory liefert Pfad zu /dateisystem anhand des usernames
func getDirectory(username string) (path string) {
	wd, err := os.Getwd() // lädt aktuelles Verzeichnis
	if err != nil {
		fmt.Println(err)
	}

	if filepath.Base(wd) != "DHBW_GO_Projekt" { // prüft ob der Projektordner dem wd entspricht, falls nein wird das letzt Pfad Element gelöscht
		wd = filepath.Dir(wd)
	}
	wd = filepath.Join(wd, "dateisystem")

	path = filepath.Join(wd, username) //formt Dateipfad, indem username an das wd angehangen wird

	return path
}

// getFileNameByTitle liefert Pfad zur Datei anhand des Titels und des Usernamens, relevant zum Suchen gespeicherter Termine
func getFileNameByTitle(title string, username string) string {
	wd, err := os.Getwd() // lädt aktuelles Verzeichnis
	if err != nil {
		fmt.Println(err)
	}

	if filepath.Base(wd) != "DHBW_GO_Projekt" { // prüft ob der Projektordner dem wd entspricht, falls nein wird das letzt Pfad Element gelöscht
		wd = filepath.Dir(wd)
	}

	wd = filepath.Join(wd, "dateisystem")

	path := filepath.Join(wd, username) //formt Dateipfad, indem username an das wd angehangen wird
	fileName := title + ".json"

	path = filepath.Join(path, fileName)
	return path
}

// getFileNameByTerminObj liefert Pfad zur Datei anhand des Titel-Objektes und des Usernamens, relevant zum Speichern von Termine
func getFileNameByTerminObj(termin Termin, username string) (pathJSON string) {
	path := getDirectory(username)

	file := termin.Title + ".json"
	pathJSON = filepath.Join(path, file) //formt Dateipfad, inklusive Dateinamen

	return pathJSON
}
