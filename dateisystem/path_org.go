package dateisystem

/*
Zweck: für das Dateisystem relative Pfade bereit zustellen um zu garantieren, das aus anderen Modulen korrekt auf Dateien zugegriffen werden kann
*/

//Mat-Nr. 8689159
import (
	"fmt"
	"os"
	"path/filepath"
)

func getDirectory(username string) (path string) { //liefert Pfad zu /dateisystem anhand des usernames
	wd, err := os.Getwd() // lädt aktuelles Verzeichnis
	if err != nil {
		fmt.Println(err)
	}

	if filepath.Base(wd) != "DHBW_GO_Projekt" { // prüft ob der Projektordner dem wd entspricht, falls nein wir das letzt Pfad Element gelöscht
		wd = filepath.Dir(wd)
	}
	wd = filepath.Join(wd, "dateisystem")

	path = filepath.Join(wd, username) //formt Dateipfad, indem username an das wd angehangen wird

	return path
}

func getFileNameByTitle(title string, username string) string { //liefert Pfad zur Datei anhand des Titels und des Usernamens, relevant zum Suchen gespeicherter Termine
	wd, err := os.Getwd() // lädt aktuelles Verzeichnis
	if err != nil {
		fmt.Println(err)
	}

	if filepath.Base(wd) != "DHBW_GO_Projekt" { // prüft ob der Projektordner dem wd entspricht, falls nein wir das letzt Pfad Element gelöscht
		wd = filepath.Dir(wd)
	}

	wd = filepath.Join(wd, "dateisystem")

	path := filepath.Join(wd, username) //formt Dateipfad, indem username an das wd angehangen wird
	fileName := title + ".json"

	path = filepath.Join(path, fileName)
	return path
}

func getFileNameByTerminObj(termin Termin, username string) (pathJSON string) { //liefert Pfad zur Datei anhand des Titel-Objektes und des Usernamens, relevant zum Speichern von Termine
	path := getDirectory(username)

	file := termin.Title + ".json"
	pathJSON = filepath.Join(path, file) //formt Dateipfad, inklusive Dateinamen

	return pathJSON
}
